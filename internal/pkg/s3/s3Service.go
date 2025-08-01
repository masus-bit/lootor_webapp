package s3

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/rwcarlsen/goexif/exif"
)

type S3Service struct {
	endpoint        string
	accessKeyId     string
	secretAccessKey string
	bucketName      string
	region          string
	cacheTtl        int
	redisClient     *redis.Client
	logger          *log.Logger
	localCache      *cache.Cache
	httpClient      *http.Client
	inflight        sync.Map
	sem             chan struct{}
}

type RequestMetrics struct {
	CacheHit          bool
	CacheLatency      time.Duration
	S3Latency         time.Duration
	RedisLatency      time.Duration
	TotalLatency      time.Duration
	ResponseSizeBytes int
	Error             string
}

type ThumbnailPreset struct {
	Name    string
	Width   int
	Height  int
	Quality int
}

var DefaultThumbnailPresets = []ThumbnailPreset{
	{"small", 320, 240, 80},
	{"medium", 640, 480, 85},
	{"large", 1280, 720, 90},
}

type UploadOptions struct {
	Width   int
	Height  int
	Quality int
}

type DeleteFilesRequest struct {
	Keys []string `json:"keys"`
}

type DeleteFilesResponse struct {
	Success  bool     `json:"success"`
	Error    string   `json:"error,omitempty"`
	NotFound []string `json:"notFound,omitempty"`
}

func NewS3Service(redisClient *redis.Client) *S3Service {
	logger := log.New(os.Stdout, "[S3Service] ", log.LstdFlags|log.Lshortfile)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		logger.Printf("Failed to connect to Redis: %v", err)
	} else {
		logger.Println("Redis connected successfully")
	}

	return &S3Service{
		endpoint:        "https://lootor.storage.yandexcloud.net",
		accessKeyId:     os.Getenv("YC_ACCESS_KEY"),
		secretAccessKey: os.Getenv("YC_SECRET_KEY"),
		bucketName:      os.Getenv("YC_BUCKET_NAME"),
		region:          os.Getenv("YC_REGION"),
		cacheTtl:        36000,
		redisClient:     redisClient,
		logger:          logger,
		localCache:      cache.New(5*time.Minute, 10*time.Minute),
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 50,
				IdleConnTimeout:     90 * time.Second,
			},
			Timeout: 3 * time.Second,
		},
		sem: make(chan struct{}, 100),
	}
}

func (s *S3Service) logMetrics(ctx context.Context, key string, metrics *RequestMetrics) {
	s.logger.Printf("[METRICS] key=%s cache_hit=%v cache_latency=%dµs redis_latency=%dµs s3_latency=%dµs total=%dµs size=%d error=%q",
		key,
		metrics.CacheHit,
		metrics.CacheLatency.Microseconds(),
		metrics.RedisLatency.Microseconds(),
		metrics.S3Latency.Microseconds(),
		metrics.TotalLatency.Microseconds(),
		metrics.ResponseSizeBytes,
		metrics.Error,
	)

}

func (s *S3Service) GetFile(ctx context.Context, key string) ([]byte, error) {
	type result struct {
		data []byte
		err  error
	}

	resChan := make(chan result, 1)
	metrics := &RequestMetrics{}
	startTime := time.Now()

	go func() {
		defer func() {
			metrics.TotalLatency = time.Since(startTime)
			s.logMetrics(ctx, key, metrics)
		}()

		if data, found := s.localCache.Get(key); found {
			metrics.CacheHit = true
			metrics.ResponseSizeBytes = len(data.([]byte))
			resChan <- result{data.([]byte), nil}
			return
		}

		redisStart := time.Now()
		data, err := s.redisClient.Get(ctx, key).Bytes()
		metrics.RedisLatency = time.Since(redisStart)

		if err == nil {
			metrics.CacheHit = true
			metrics.ResponseSizeBytes = len(data)
			s.localCache.SetDefault(key, data)
			resChan <- result{data, nil}
			return
		}

		s.sem <- struct{}{}
		defer func() { <-s.sem }()

		s3Start := time.Now()
		s3Data, err := s.downloadFromS3Optimized(ctx, key)
		metrics.S3Latency = time.Since(s3Start)

		if err != nil {
			metrics.Error = err.Error()
			resChan <- result{nil, err}
			return
		}

		metrics.ResponseSizeBytes = len(s3Data)
		resChan <- result{s3Data, nil}

		go s.updateCaches(key, s3Data)
	}()
	source := "local_cache"
	if !metrics.CacheHit {
		source = "redis"
		if metrics.S3Latency > 0 {
			source = "s3"
		}
	}
	s.logger.Printf("[SOURCE] key=%s source=%s total=%dµs", key, source, metrics.TotalLatency.Microseconds())
	select {
	case res := <-resChan:
		return res.data, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *S3Service) getFileAsync(ctx context.Context, key string) ([]byte, error) {
	if data, found := s.localCache.Get(key); found {
		return data.([]byte), nil
	}

	val, _ := s.inflight.LoadOrStore(key, new(sync.WaitGroup))
	wg := val.(*sync.WaitGroup)
	wg.Add(1)
	defer wg.Done()
	defer s.inflight.Delete(key)

	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err == nil {
		s.localCache.SetDefault(key, data)
		return data, nil
	}

	s.sem <- struct{}{}
	defer func() { <-s.sem }()

	s3Data, err := s.downloadFromS3Optimized(ctx, key)
	if err != nil {
		return nil, err
	}

	go s.updateCaches(key, s3Data)

	return s3Data, nil
}

func (s *S3Service) updateCaches(key string, data []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.localCache.Set(key, data, cache.DefaultExpiration)
	}()

	go func() {
		defer wg.Done()
		s.redisClient.Set(ctx, key, data, time.Duration(s.cacheTtl)*time.Second)
	}()

	wg.Wait()
}

func (s *S3Service) downloadFromS3Optimized(ctx context.Context, key string) ([]byte, error) {
	start := time.Now()
	defer func() {
		s.logger.Printf("[S3_DOWNLOAD] key=%s duration=%v", key, time.Since(start))
	}()

	signStart := time.Now()
	signed, err := s.signS3Request(ctx, "GET", "/images/"+key, nil, nil)
	s.logger.Printf("[S3_SIGN] key=%s duration=%v", key, time.Since(signStart))

	if err != nil {
		return nil, err
	}

	reqStart := time.Now()
	req, err := http.NewRequestWithContext(ctx, signed.method, signed.url, nil)
	s.logger.Printf("[S3_REQ_CREATE] key=%s duration=%v", key, time.Since(reqStart))

	if err != nil {
		return nil, err
	}

	for k, v := range signed.headers {
		req.Header.Set(k, v)
	}

	s.logger.Printf("[S3_HEADERS] key=%s headers=%+v", key, signed.headers)

	clientStart := time.Now()
	resp, err := s.httpClient.Do(req)
	s.logger.Printf("[S3_CLIENT_DO] key=%s duration=%v", key, time.Since(clientStart))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	s.logger.Printf("[S3_RESPONSE] key=%s status=%d content_length=%s",
		key, resp.StatusCode, resp.Header.Get("Content-Length"))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("S3 returned status: %d", resp.StatusCode)
	}

	readStart := time.Now()
	data, err := io.ReadAll(resp.Body)
	s.logger.Printf("[S3_READ_BODY] key=%s duration=%v size=%d",
		key, time.Since(readStart), len(data))

	return data, err
}

func (s *S3Service) downloadFromS3(ctx context.Context, key string) ([]byte, error) {
	signed, err := s.signS3Request(ctx, "GET", "/images/"+key, nil, nil)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, signed.method, signed.url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range signed.headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (s *S3Service) GetStats() map[string]string {
	stats := make(map[string]string)

	stats["local_cache_items"] = strconv.Itoa(s.localCache.ItemCount())

	redisStats, err := s.redisClient.Info(context.Background(), "stats").Result()
	if err == nil {
		stats["redis_ops_per_sec"] = extractRedisStat(redisStats, "instantaneous_ops_per_sec")
		stats["redis_hit_rate"] = extractRedisStat(redisStats, "keyspace_hits") + "/" +
			extractRedisStat(redisStats, "keyspace_misses")
	}

	if t, ok := s.httpClient.Transport.(*http.Transport); ok {
		idleConns := reflect.ValueOf(t).Elem().FieldByName("idleConn")
		if idleConns.IsValid() {
			stats["http_idle_conns"] = strconv.Itoa(idleConns.Len())
		}
	}

	return stats
}

func extractRedisStat(info, key string) string {
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		metricName := strings.TrimSpace(parts[0])
		if metricName == key {
			return strings.TrimSpace(strings.TrimRight(parts[1], "\r"))
		}
	}
	return "0"
}

func (s *S3Service) signS3Request(ctx context.Context, method, path string, headers map[string]string, body []byte) (struct {
	url     string
	headers map[string]string
	method  string
	data    []byte
}, error) {
	now := time.Now().UTC()
	date := now.Format("20060102T150405Z")

	host := strings.TrimPrefix(s.endpoint, "https://")
	canonicalURI := path
	canonicalQueryString := ""
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:UNSIGNED-PAYLOAD\nx-amz-date:%s\n", host, date)
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	payloadHash := "UNSIGNED-PAYLOAD"

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		payloadHash)

	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", now.Format("20060102"), s.region)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		date,
		credentialScope,
		hex.EncodeToString(hashSHA256([]byte(canonicalRequest))))

	signingKey := getSignatureKey(s.secretAccessKey, now.Format("20060102"), s.region, "s3")
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	authorizationHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.accessKeyId,
		credentialScope,
		signedHeaders,
		signature)

	resultHeaders := map[string]string{
		"Host":                 host,
		"x-amz-date":           date,
		"Authorization":        authorizationHeader,
		"x-amz-content-sha256": "UNSIGNED-PAYLOAD",
	}
	for k, v := range headers {
		resultHeaders[k] = v
	}

	return struct {
		url     string
		headers map[string]string
		method  string
		data    []byte
	}{
		url:     s.endpoint + path,
		headers: resultHeaders,
		method:  method,
		data:    body,
	}, nil
}

func hashSHA256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func hmacSHA256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func getSignatureKey(key, dateStamp, regionName, serviceName string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+key), []byte(dateStamp))
	kRegion := hmacSHA256(kDate, []byte(regionName))
	kService := hmacSHA256(kRegion, []byte(serviceName))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	return kSigning
}

func (s *S3Service) UploadFile(ctx context.Context, file []byte, key, contentType string) error {
	path := "/" + key
	signed, err := s.signS3Request(ctx, "PUT", path, map[string]string{
		"Content-Type": contentType,
	}, file)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, signed.method, signed.url, bytes.NewReader(file))
	if err != nil {
		return err
	}
	for k, v := range signed.headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	return nil
}

func (s *S3Service) UploadOptimizedImages(ctx context.Context, files [][]byte, options UploadOptions) ([]string, error) {
	var keys []string

	for _, file := range files {
		key, err := s.UploadOptimizedImage(ctx, file, uuid.New().String(), options)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func (s *S3Service) UploadOptimizedImage(ctx context.Context, file []byte, filename string, options UploadOptions) (string, error) {
	width := options.Width
	if width == 0 {
		width = 1920
	}
	quality := options.Quality
	if quality == 0 {
		quality = 70
	}

	img, err := decodeImageWithOrientation(file)
	if err != nil {
		return "", err
	}

	if width > 0 || options.Height > 0 {
		img = imaging.Resize(img, width, options.Height, imaging.Lanczos)
	}

	var buf bytes.Buffer
	err = webp.Encode(&buf, img, &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	})
	if err != nil {
		return "", fmt.Errorf("failed to encode webp: %v", err)
	}

	ext := filepath.Ext(filename)
	if ext != "" {
		filename = strings.TrimSuffix(filename, ext)
	}
	key := fmt.Sprintf("images/%d-%s.webp", time.Now().Unix(), filename)

	err = s.UploadFile(ctx, buf.Bytes(), key, "image/webp")
	if err != nil {
		return "", err
	}

	return key, nil
}

func decodeImageWithOrientation(data []byte) (image.Image, error) {
	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	exifData, err := exif.Decode(bytes.NewReader(data))
	if err == nil {
		tag, err := exifData.Get(exif.Orientation)
		if err == nil {
			orientation, err := tag.Int(0)
			if err == nil {
				switch orientation {
				case 2:
					img = imaging.FlipH(img)
				case 3:
					img = imaging.Rotate180(img)
				case 4:
					img = imaging.FlipV(img)
				case 5:
					img = imaging.Transpose(img)
				case 6:
					img = imaging.Rotate270(img)
				case 7:
					img = imaging.Transverse(img)
				case 8:
					img = imaging.Rotate90(img)
				}
			}
		}
	}

	return img, nil
}

func (s *S3Service) DeleteFiles(ctx context.Context, req DeleteFilesRequest) (DeleteFilesResponse, error) {
	var deleted []string
	var notFound []string

	for _, key := range req.Keys {
		path := "/" + key
		signed, err := s.signS3Request(ctx, "DELETE", path, nil, nil)
		if err != nil {
			return DeleteFilesResponse{}, err
		}

		req, err := http.NewRequestWithContext(ctx, signed.method, signed.url, nil)
		if err != nil {
			return DeleteFilesResponse{}, err
		}
		for k, v := range signed.headers {
			req.Header.Set(k, v)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return DeleteFilesResponse{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			deleted = append(deleted, key)
			if err := s.redisClient.Del(ctx, key).Err(); err != nil {
				s.logger.Printf("Failed to delete key from cache: %s, error: %v", key, err)
			}
		} else if resp.StatusCode == http.StatusNotFound {
			notFound = append(notFound, key)
		}
	}

	go s.deleteRelatedThumbnails(context.Background(), req.Keys)

	if len(notFound) == 0 {
		return DeleteFilesResponse{
			Success: true,
		}, nil
	}

	return DeleteFilesResponse{
		Success:  len(deleted) > 0,
		Error:    "some files not found",
		NotFound: notFound,
	}, nil
}

func (s *S3Service) GenerateThumbnail(ctx context.Context, filename string, width, height, quality int) ([]byte, error) {
	originalKey := filename
	thumbKey := fmt.Sprintf("thumbs/%dx%d/%s", width, height, filename)

	if data, err := s.GetFile(ctx, thumbKey); err == nil {
		return data, nil
	}

	originalData, err := s.downloadFromS3Optimized(ctx, originalKey)
	if err != nil {
		return nil, fmt.Errorf("original image not found: %w", err)
	}

	img, err := decodeImageWithOrientation(originalData)
	if err != nil {
		return nil, fmt.Errorf("image decode failed: %w", err)
	}

	resized := imaging.Resize(img, width, height, imaging.Lanczos)

	var buf bytes.Buffer
	if err := webp.Encode(&buf, resized, &webp.Options{
		Quality: float32(quality),
	}); err != nil {
		return nil, fmt.Errorf("webp encode failed: %w", err)
	}

	thumbData := buf.Bytes()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.UploadFile(ctx, thumbData, thumbKey, "image/webp"); err != nil {
			s.logger.Printf("Failed to save thumbnail %s: %v", thumbKey, err)
		} else {
			s.updateCaches(thumbKey, thumbData)
		}
	}()

	return thumbData, nil
}

func (s *S3Service) FileExists(ctx context.Context, key string) (bool, error) {
	signed, err := s.signS3Request(ctx, "HEAD", key, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequestWithContext(ctx, signed.method, signed.url, nil)
	if err != nil {
		return false, err
	}

	for k, v := range signed.headers {
		req.Header.Set(k, v)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

func (s *S3Service) deleteRelatedThumbnails(ctx context.Context, keys []string) {
	for _, key := range keys {
		sizes := []struct{ w, h int }{
			{100, 0}, {120, 0}, {150, 0}, {180, 0},
			{200, 0}, {20, 0}, {300, 0}, {30, 0},
			{36, 0}, {38, 0}, {44, 0}, {470, 0},
			{50, 0}, {600, 0}, {700, 0}, {80, 0},
		}

		for _, size := range sizes {
			thumbKey := s.GetThumbnailKey(key, size.w, size.h)
			if err := s.redisClient.Del(ctx, thumbKey).Err(); err != nil {
				s.logger.Printf("Failed to delete thumbnail cache: %s", thumbKey)
			}
		}
	}
}

func (s *S3Service) GetOriginalKey(filename string) string {
	return "images/" + filename
}

func (s *S3Service) GetThumbnailKey(filename string, width, height int) string {
	return fmt.Sprintf("thumbs/%dx0/%s", width, filename)
}
