package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type TagDocument struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Slug                 string `json:"slug"`
	PrimaryID            string `json:"primaryId"`
	SeriesID             string `json:"seriesId"`
	TotalCollectionItems int64  `json:"totalCollectionItems"`
	TotalCollections     int64  `json:"totalCollections"`
	TotalPhotos          int64  `json:"totalPhotos"`
	TotalPosts           int64  `json:"totalPosts"`
}

type ElasticConfig struct {
	Address string                 `json:"address"`
	Indices map[string]IndexConfig `json:"indices"`
}

type IndexConfig struct {
	Settings map[string]interface{} `json:"settings"`
	Mappings map[string]interface{} `json:"mappings"`
}

type SearchResult struct {
	Hits struct {
		Hits []struct {
			Index  string                 `json:"_index"`
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type SearchResultFormatted struct {
	Result []struct {
		Index  string                 `json:"_index"`
		Source map[string]interface{} `json:"_source"`
	} `json:"result"`
}

type ElasticService struct {
	client *elasticsearch.Client
	config ElasticConfig
	logger *log.Logger
}

func LoadConfig(path string) (ElasticConfig, error) {
	// #nosec G304
	file, err := os.ReadFile(path)
	if err != nil {
		return ElasticConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ElasticConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return ElasticConfig{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}

func NewElasticService(configPath string) (*ElasticService, error) {
	// #nosec G304
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	es, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{os.Getenv("ELASTIC_URL")},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create elastic client: %w", err)
	}

	logger := log.New(os.Stdout, "[Elastic] ", log.LstdFlags|log.Lshortfile)

	res, err := es.Ping()
	if err != nil {
		return nil, fmt.Errorf("elasticsearch ping failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("elasticsearch ping error")
	}

	return &ElasticService{
		client: es,
		config: config,
		logger: logger,
	}, nil
}

func (es *ElasticService) ReindexAll(
	ctx context.Context,
	dataProviders map[string]func() ([]map[string]interface{}, error),
) error {
	for indexName, provider := range dataProviders {
		data, err := provider()
		if err != nil {
			return fmt.Errorf("failed to get data for %s: %w", indexName, err)
		}

		if len(data) == 0 {
			es.logger.Printf("Skipping %s: no data", indexName)
			continue
		}

		// НЕ УДАЛЯЕМ ИНДЕКС!
		// if err := es.deleteIndexIfExists(indexName); err != nil {
		// 	return fmt.Errorf("failed to delete index %s: %w", indexName, err)
		// }
	}

	for indexName, provider := range dataProviders {
		data, err := provider()
		if err != nil {
			return fmt.Errorf("failed to get data for %s: %w", indexName, err)
		}

		if len(data) == 0 {
			es.logger.Printf("Skipping %s: no data", indexName)
			continue
		}

		// Создаем индекс только если он не существует
		exists, err := es.indexExists(context.Background(), indexName)
		if err != nil {
			return fmt.Errorf("failed to check index existence: %w", err)
		}

		if !exists {
			if err := es.createIndex(indexName); err != nil {
				return fmt.Errorf("failed to create index %s: %w", indexName, err)
			}
		}

		// Используем bulk update вместо bulk create
		if err := es.bulkUpdateDocuments(ctx, indexName, data); err != nil {
			return fmt.Errorf("failed to update data for %s: %w", indexName, err)
		}
	}

	return nil
}

func (es *ElasticService) bulkUpdateDocuments(
	ctx context.Context,
	indexName string,
	docs []map[string]interface{},
) error {
	var buf strings.Builder

	for _, doc := range docs {
		id, ok := doc["id"].(string)
		if !ok {
			return fmt.Errorf("document missing id field")
		}

		// Используем update вместо index
		meta := map[string]interface{}{
			"update": map[string]interface{}{
				"_index": indexName,
				"_id":    id,
			},
		}

		// Используем doc_as_upsert для обновления или создания
		update := map[string]interface{}{
			"doc":           doc,
			"doc_as_upsert": true,
		}

		metaJSON, _ := json.Marshal(meta)
		updateJSON, _ := json.Marshal(update)

		buf.Write(metaJSON)
		buf.WriteString("\n")
		buf.Write(updateJSON)
		buf.WriteString("\n")
	}

	res, err := es.client.Bulk(
		strings.NewReader(buf.String()),
		es.client.Bulk.WithContext(ctx),
		es.client.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("error parsing bulk response: %w", err)
	}

	if errors, ok := result["errors"].(bool); ok && errors {
		return fmt.Errorf("bulk operation contains errors")
	}

	es.logger.Printf("Updated %d documents in %s", len(docs), indexName)
	return nil
}

func (es *ElasticService) SearchInIndices(
	ctx context.Context,
	indices []string,
	query string,
	limit string,
) (*SearchResultFormatted, error) {
	if strings.TrimSpace(query) == "" {
		return &SearchResultFormatted{}, nil
	}
	limitInt, _ := strconv.Atoi(limit)

	searchQuery := es.buildSearchQuery(query, limit)
	res, err := es.search(ctx, indices, searchQuery)
	if err != nil {
		return nil, err
	}

	var filteredHits []struct {
		Index  string                 `json:"_index"`
		Source map[string]interface{} `json:"_source"`
	}

	indexCounts := make(map[string]int)
	for _, hit := range res.Hits.Hits {
		if indexCounts[hit.Index] >= limitInt {
			continue
		}
		filteredHits = append(filteredHits, hit)
		indexCounts[hit.Index]++
	}

	return &SearchResultFormatted{Result: filteredHits}, nil
}

func (es *ElasticService) deleteIndexIfExists(indexName string) error {
	res, err := es.client.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil
	}

	res, err = es.client.Indices.Delete([]string{indexName})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	es.logger.Printf("Index %s deleted", indexName)
	return nil
}

func (es *ElasticService) createIndex(indexName string) error {
	cfg, exists := es.config.Indices[indexName]
	if !exists {
		return fmt.Errorf("no config found for index %s", indexName)
	}

	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(cfg); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, err := es.client.Indices.Create(
		indexName,
		es.client.Indices.Create.WithBody(strings.NewReader(buf.String())),
		es.client.Indices.Create.WithContext(ctx), // Важно: добавляем контекст
	)
	if err != nil {
		return fmt.Errorf("create index request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	es.logger.Printf("Index %s created successfully", indexName)
	return nil
}

func (es *ElasticService) bulkIndexDocuments(
	ctx context.Context,
	indexName string,
	docs []map[string]interface{},
) error {
	var buf strings.Builder

	for _, doc := range docs {
		id, ok := doc["id"].(string)
		if !ok {
			return fmt.Errorf("document missing id field")
		}

		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_id":    id,
			},
		}

		metaJSON, _ := json.Marshal(meta)
		docJSON, _ := json.Marshal(doc)

		buf.Write(metaJSON)
		buf.WriteString("\n")
		buf.Write(docJSON)
		buf.WriteString("\n")
	}

	res, err := es.client.Bulk(
		strings.NewReader(buf.String()),
		es.client.Bulk.WithContext(ctx),
		es.client.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	// Проверяем ошибки в bulk операции
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("error parsing bulk response: %w", err)
	}

	if errors, ok := result["errors"].(bool); ok && errors {
		return fmt.Errorf("bulk operation contains errors")
	}

	es.logger.Printf("Indexed %d documents to %s", len(docs), indexName)
	return nil
}

func (es *ElasticService) buildSearchQuery(query string, limit string) map[string]interface{} {
	limitInt, _ := strconv.Atoi(limit)
	lowerQuery := strings.ToLower(strings.TrimSpace(query))
	queryLen := len(lowerQuery)

	var prefixMustNot []map[string]interface{}
	if queryLen > 3 {
		prefixMustNot = []map[string]interface{}{
			{"term": map[string]interface{}{"login.keyword": lowerQuery}},
		}
	}

	var namePrefixMustNot []map[string]interface{}
	if queryLen > 3 {
		namePrefixMustNot = []map[string]interface{}{
			{"term": map[string]interface{}{"name.keyword": lowerQuery}},
		}
	}

	var profileNamePrefixMustNot []map[string]interface{}
	if queryLen > 3 {
		profileNamePrefixMustNot = []map[string]interface{}{
			{"term": map[string]interface{}{"profileName.keyword": lowerQuery}},
		}
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"_index": "users"}},
								{"term": map[string]interface{}{"login.keyword": lowerQuery}},
							},
							"boost": 2.0,
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"_index": "users"}},
								{
									"prefix": map[string]interface{}{
										"login.keyword": map[string]interface{}{
											"value":            lowerQuery,
											"case_insensitive": true,
										},
									},
								},
							},
							"must_not": prefixMustNot,
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"_index": "users"}},
								{"term": map[string]interface{}{"profileName.keyword": lowerQuery}},
							},
							"boost": 2.0,
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"_index": "users"}},
								{
									"prefix": map[string]interface{}{
										"profileName.keyword": map[string]interface{}{
											"value":            lowerQuery,
											"case_insensitive": true,
										},
									},
								},
							},
							"must_not": profileNamePrefixMustNot,
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"_index": "collections"}},
								{"term": map[string]interface{}{"name.keyword": lowerQuery}},
							},
							"boost": 2.0,
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"_index": "tags"}},
								{"term": map[string]interface{}{"name.keyword": lowerQuery}},
							},
							"boost": 2.0,
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"_index": "collection_items"}},
								{"term": map[string]interface{}{"name.keyword": lowerQuery}},
							},
							"boost": 2.0,
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{
									"bool": map[string]interface{}{
										"should": []map[string]interface{}{
											{"term": map[string]interface{}{"_index": "collections"}},
											{"term": map[string]interface{}{"_index": "tags"}},
											{"term": map[string]interface{}{"_index": "collection_items"}},
										},
									},
								},
								{
									"match": map[string]interface{}{
										"name.full": map[string]interface{}{
											"query":    lowerQuery,
											"operator": "or",
										},
									},
								},
							},
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{
									"bool": map[string]interface{}{
										"should": []map[string]interface{}{
											{"term": map[string]interface{}{"_index": "collections"}},
											{"term": map[string]interface{}{"_index": "tags"}},
											{"term": map[string]interface{}{"_index": "collection_items"}},
										},
									},
								},
								{
									"prefix": map[string]interface{}{
										"name.keyword": map[string]interface{}{
											"value":            lowerQuery,
											"case_insensitive": true,
										},
									},
								},
							},
							"must_not": namePrefixMustNot,
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
		"size": limitInt * 5,
	}
}

func (es *ElasticService) search(ctx context.Context, indices []string, query map[string]interface{}) (
	*SearchResult,
	error,
) {
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex(indices...),
		es.client.Search.WithBody(strings.NewReader(buf.String())),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, parseErrorResponse(res)
	}

	var result SearchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func parseErrorResponse(res *esapi.Response) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading error response: %w", err)
	}

	var e map[string]interface{}
	if err := json.Unmarshal(body, &e); err != nil {
		return fmt.Errorf("error parsing error response: %w", err)
	}

	return fmt.Errorf("elasticsearch error [%d]: %v", res.StatusCode, e["error"])
}

func (es *ElasticService) IndexDocument(ctx context.Context, index string, doc map[string]interface{}) error {
	ctxCheck, cancelCheck := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCheck()

	exists, err := es.indexExists(ctxCheck, index)
	if err != nil {
		return fmt.Errorf("index check failed: %w", err)
	}

	if !exists {
		if err := es.createIndex(index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	id, ok := doc["id"].(string)
	if !ok {
		return fmt.Errorf("document missing id field")
	}

	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return fmt.Errorf("error encoding document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(buf.String()),
		Refresh:    "wait_for", // Ждем обновления индекса
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("index request error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	return nil
}

func (es *ElasticService) BulkIndexDocuments(ctx context.Context, index string, docs []map[string]interface{}) error {
	if len(docs) == 0 {
		return nil
	}

	ctxCheck, cancelCheck := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCheck()

	exists, err := es.indexExists(ctxCheck, index)
	if err != nil {
		return fmt.Errorf("index check failed: %w", err)
	}

	if !exists {
		if err := es.createIndex(index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	var buf strings.Builder
	for _, doc := range docs {
		id, ok := doc["id"].(string)
		if !ok {
			continue
		}

		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": index,
				"_id":    id,
			},
		}

		metaJSON, err := json.Marshal(meta)
		if err != nil {
			return fmt.Errorf("error encoding metadata: %w", err)
		}

		docJSON, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("error encoding document: %w", err)
		}

		buf.WriteString(string(metaJSON))
		buf.WriteString("\n")
		buf.WriteString(string(docJSON))
		buf.WriteString("\n")
	}

	if buf.Len() == 0 {
		return fmt.Errorf("no valid documents to index")
	}

	req := esapi.BulkRequest{
		Body:    strings.NewReader(buf.String()),
		Refresh: "wait_for",
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("bulk request error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return fmt.Errorf("error parsing bulk response: %w", err)
	}

	if errors, exists := response["errors"].(bool); exists && errors {
		log.Printf("Bulk operation completed with errors: %v", response)
	}

	return nil
}

func (es *ElasticService) DeleteDocument(ctx context.Context, index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("delete request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil
	} else if res.IsError() {
		return parseErrorResponse(res)
	}

	return nil
}

func (es *ElasticService) indexExists(ctx context.Context, indexName string) (bool, error) {
	res, err := es.client.Indices.Exists(
		[]string{indexName},
		es.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, nil
}

func (es *ElasticService) GetDocument(ctx context.Context, index, id string) (map[string]interface{}, error) {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return nil, fmt.Errorf("get document request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("document not found: %s", id)
	} else if res.IsError() {
		return nil, parseErrorResponse(res)
	}

	var result struct {
		Source map[string]interface{} `json:"_source"`
		Found  bool                   `json:"found"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing get response: %w", err)
	}

	if !result.Found {
		return nil, fmt.Errorf("document not found: %s", id)
	}

	return result.Source, nil
}

func (es *ElasticService) IncrementField(ctx context.Context, index, id, field string, increment int64) error {
	searchBody := fmt.Sprintf(
		`{
		"query": {
			"term": {
				"id": "%s"
			}
		},
		"_source": ["*"],
		"size": 1
	}`, id,
	)

	searchReq := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(searchBody),
	}

	searchRes, err := searchReq.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("search document error: %w", err)
	}
	defer searchRes.Body.Close()

	searchBodyBytes, _ := io.ReadAll(searchRes.Body)

	var tagDoc TagDocument

	if searchRes.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(searchBodyBytes, &result); err != nil {
			return fmt.Errorf("decode search response error: %w", err)
		}

		hits, ok := result["hits"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid search response format")
		}

		hitsList, ok := hits["hits"].([]interface{})
		if !ok {
			return fmt.Errorf("invalid hits format")
		}

		if len(hitsList) == 0 {
			tagDoc = TagDocument{
				ID:                   id,
				Name:                 "",
				Slug:                 "",
				PrimaryID:            "",
				SeriesID:             "",
				TotalCollectionItems: 0,
				TotalCollections:     0,
				TotalPhotos:          0,
				TotalPosts:           0,
			}
		} else {
			firstHit := hitsList[0].(map[string]interface{})

			if source, ok := firstHit["_source"].(map[string]interface{}); ok {

				tagDoc = convertToTagDocument(source, id)
			} else {
				log.Printf("[ERROR] No _source in hit")
				return fmt.Errorf("no source in document")
			}
		}
	} else {
		return parseErrorResponse(searchRes)
	}

	if tagDoc.TotalCollectionItems < 0 {
		tagDoc.TotalCollectionItems = 0
	}
	if tagDoc.TotalCollections < 0 {
		tagDoc.TotalCollections = 0
	}
	if tagDoc.TotalPhotos < 0 {
		tagDoc.TotalPhotos = 0
	}
	if tagDoc.TotalPosts < 0 {
		tagDoc.TotalPosts = 0
	}

	switch field {
	case "totalCollectionItems":
		tagDoc.TotalCollectionItems += increment
		if tagDoc.TotalCollectionItems < 0 {
			tagDoc.TotalCollectionItems = 0
		}
	case "totalCollections":
		tagDoc.TotalCollections += increment
		if tagDoc.TotalCollections < 0 {
			tagDoc.TotalCollections = 0
		}
	case "totalPhotos":
		tagDoc.TotalPhotos += increment
		if tagDoc.TotalPhotos < 0 {
			tagDoc.TotalPhotos = 0
		}
	case "totalPosts":
		tagDoc.TotalPosts += increment
		if tagDoc.TotalPosts < 0 {
			tagDoc.TotalPosts = 0
		}
	default:
		return fmt.Errorf("unknown field: %s", field)
	}

	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(tagDoc); err != nil {
		return fmt.Errorf("error encoding document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(buf.String()),
		Refresh:    "wait_for",
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("index document error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	return nil
}

func (es *ElasticService) DecrementField(ctx context.Context, index, id, field string, decrement int64) error {
	searchBody := fmt.Sprintf(
		`{
		"query": {
			"term": {
				"id": "%s"
			}
		},
		"_source": ["*"],
		"size": 1
	}`, id,
	)

	searchReq := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(searchBody),
	}

	searchRes, err := searchReq.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("search document error: %w", err)
	}
	defer searchRes.Body.Close()

	searchBodyBytes, _ := io.ReadAll(searchRes.Body)

	var tagDoc TagDocument

	if searchRes.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(searchBodyBytes, &result); err != nil {
			return fmt.Errorf("decode search response error: %w", err)
		}

		hits, ok := result["hits"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid search response format")
		}

		hitsList, ok := hits["hits"].([]interface{})
		if !ok || len(hitsList) == 0 {
			tagDoc = TagDocument{
				ID:                   id,
				Name:                 "",
				Slug:                 "",
				PrimaryID:            "",
				SeriesID:             "",
				TotalCollectionItems: 0,
				TotalCollections:     0,
				TotalPhotos:          0,
				TotalPosts:           0,
			}
		} else {
			firstHit := hitsList[0].(map[string]interface{})
			if source, ok := firstHit["_source"].(map[string]interface{}); ok {
				tagDoc = convertToTagDocument(source, id)
			}
		}
	} else {
		return parseErrorResponse(searchRes)
	}

	var newValue int64

	switch field {
	case "totalCollectionItems":
		newValue = tagDoc.TotalCollectionItems - decrement
		if newValue < 0 {
			newValue = 0
		}
		tagDoc.TotalCollectionItems = newValue

	case "totalCollections":
		newValue = tagDoc.TotalCollections - decrement
		if newValue < 0 {
			newValue = 0
		}
		tagDoc.TotalCollections = newValue

	case "totalPhotos":
		newValue = tagDoc.TotalPhotos - decrement
		if newValue < 0 {
			newValue = 0
		}
		tagDoc.TotalPhotos = newValue

	case "totalPosts":
		newValue = tagDoc.TotalPosts - decrement
		if newValue < 0 {
			newValue = 0
		}
		tagDoc.TotalPosts = newValue

	default:
		return fmt.Errorf("unknown field: %s", field)
	}

	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(tagDoc); err != nil {
		return fmt.Errorf("error encoding document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(buf.String()),
		Refresh:    "wait_for",
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("index document error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}
	return nil
}

func convertToTagDocument(source map[string]interface{}, id string) TagDocument {
	doc := TagDocument{
		ID:                   id,
		Name:                 "",
		Slug:                 "",
		PrimaryID:            "",
		SeriesID:             "",
		TotalCollectionItems: 0,
		TotalCollections:     0,
		TotalPhotos:          0,
		TotalPosts:           0,
	}

	if v, ok := source["name"]; ok {
		if s, isString := v.(string); isString {
			doc.Name = s
		}
	}

	if v, ok := source["slug"]; ok {
		if s, isString := v.(string); isString {
			doc.Slug = s
		}
	}

	if v, ok := source["primaryId"]; ok {
		if s, isString := v.(string); isString {
			doc.PrimaryID = s
		}
	}

	if v, ok := source["seriesId"]; ok {
		if s, isString := v.(string); isString {
			doc.SeriesID = s
		}
	}

	if v, ok := source["totalCollectionItems"]; ok {
		doc.TotalCollectionItems = toInt64(v)
	}

	if v, ok := source["totalCollections"]; ok {
		doc.TotalCollections = toInt64(v)
	}

	if v, ok := source["totalPhotos"]; ok {
		doc.TotalPhotos = toInt64(v)
	}

	if v, ok := source["totalPosts"]; ok {
		doc.TotalPosts = toInt64(v)
	}

	return doc
}

func toInt64(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int64:
		return val
	case int:
		return int64(val)
	case json.Number:
		if intVal, err := val.Int64(); err == nil {
			return intVal
		}
	}
	return 0
}

func (es *ElasticService) UpdateDocument(ctx context.Context, index, id string, fields map[string]interface{}) error {
	updateBody := map[string]interface{}{
		"doc": fields,
	}

	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(updateBody); err != nil {
		return fmt.Errorf("error encoding update body: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(buf.String()),
		Refresh:    "wait_for",
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("update request error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	return nil
}
