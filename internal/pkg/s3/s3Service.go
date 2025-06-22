package s3

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/redis/go-redis/v9"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	}
}

func (s *S3Service) GetFile(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Улучшенное логирование
	startTime := time.Now()
	s.logger.Printf("[GetFile] Starting request for key: %s", key)
	defer func() {
		s.logger.Printf("[GetFile] Completed request for key: %s (duration: %v)", key, time.Since(startTime))
	}()

	// 1. Проверка кеша с подробным логированием
	cacheStart := time.Now()
	cached, err := s.redisClient.Get(ctx, key).Bytes()
	if err == nil {
		s.logger.Printf("[GetFile] Cache HIT for key: %s (size: %d bytes, fetch time: %v)",
			key, len(cached), time.Since(cacheStart))
		return cached, nil
	}

	if err != redis.Nil {
		s.logger.Printf("[GetFile] Cache error for key: %s (error: %v, fetch time: %v)",
			key, err, time.Since(cacheStart))
	} else {
		s.logger.Printf("[GetFile] Cache MISS for key: %s (fetch time: %v)",
			key, time.Since(cacheStart))
	}

	// 2. Загрузка из S3
	s3Start := time.Now()
	file, err := s.downloadFromS3Optimized(ctx, key)
	if err != nil {
		s.logger.Printf("[GetFile] S3 download failed for key: %s (error: %v, duration: %v)",
			key, err, time.Since(s3Start))
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	s.logger.Printf("[GetFile] S3 download success for key: %s (size: %d bytes, duration: %v)",
		key, len(file), time.Since(s3Start))

	// 3. Асинхронное сохранение в кеш с обработкой ошибок
	go func() {
		cacheSaveStart := time.Now()
		if err := s.redisClient.Set(context.Background(), key, file, time.Duration(s.cacheTtl)*time.Second).Err(); err != nil {
			s.logger.Printf("[GetFile] Failed to save to cache for key: %s (error: %v, duration: %v)",
				key, err, time.Since(cacheSaveStart))
		} else {
			s.logger.Printf("[GetFile] Successfully saved to cache for key: %s (size: %d bytes, duration: %v)",
				key, len(file), time.Since(cacheSaveStart))
		}
	}()

	return file, nil
}

// Вспомогательная функция для определения источника данных
func resultSource(data []byte, key string) string {
	if len(data) > 0 && strings.HasPrefix(string(data), "REDIS_CACHE") {
		return "CACHE"
	}
	return "S3"
}

func (s *S3Service) downloadFromS3Optimized(ctx context.Context, key string) ([]byte, error) {
	s.logger.Printf("[S3] Starting download for key: %s", key)

	var httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     120 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: 5 * time.Second,
	}

	signed, err := s.signS3Request(ctx, "GET", "/images/"+key, nil, nil)
	if err != nil {
		s.logger.Printf("[S3] Signing failed for key %s: %v", key, err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, signed.method, signed.url, nil)
	if err != nil {
		s.logger.Printf("[S3] Request creation failed for key %s: %v", key, err)
		return nil, err
	}

	for k, v := range signed.headers {
		req.Header.Set(k, v)
	}

	start := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		s.logger.Printf("[S3] Request failed for key %s: %v (time: %v)",
			key, err, time.Since(start))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Printf("[S3] Non-200 status for key %s: %d (time: %v)",
			key, resp.StatusCode, time.Since(start))
		return nil, fmt.Errorf("S3 returned status: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		s.logger.Printf("[S3] Body read failed for key %s: %v (time: %v)",
			key, err, time.Since(start))
		return nil, err
	}

	s.logger.Printf("[S3] Download completed for key: %s (size: %d KB, total time: %v)",
		key, buf.Len()/1024, time.Since(start))

	return buf.Bytes(), nil
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
