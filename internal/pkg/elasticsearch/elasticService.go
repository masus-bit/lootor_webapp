package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

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
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.Address},
	})
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

func (es *ElasticService) ReindexAll(ctx context.Context, dataProviders map[string]func() ([]map[string]interface{}, error)) error {
	for indexName := range dataProviders {
		if err := es.deleteIndexIfExists(indexName); err != nil {
			return fmt.Errorf("failed to delete index %s: %w", indexName, err)
		}
	}

	for indexName := range dataProviders {
		if err := es.createIndex(indexName); err != nil {
			return fmt.Errorf("failed to create index %s: %w", indexName, err)
		}
	}

	for indexName, provider := range dataProviders {
		data, err := provider()
		if err != nil {
			return fmt.Errorf("failed to get data for %s: %w", indexName, err)
		}

		if err := es.bulkIndexDocuments(ctx, indexName, data); err != nil {
			return fmt.Errorf("failed to index data for %s: %w", indexName, err)
		}
	}

	return nil
}

func (es *ElasticService) SearchInIndices(ctx context.Context, indices []string, query string, limit string) (*SearchResultFormatted, error) {
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

func getIndexPosition(index string) int {
	switch index {
	case "users":
		return 0
	case "collections":
		return 1
	case "tags":
		return 2
	case "collection_items":
		return 3
	case "entities":
		return 4
	default:
		return -1
	}
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

	res, err := es.client.Indices.Create(
		indexName,
		es.client.Indices.Create.WithBody(strings.NewReader(buf.String())),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	es.logger.Printf("Index %s created", indexName)
	return nil
}

func (es *ElasticService) bulkIndexDocuments(ctx context.Context, indexName string, docs []map[string]interface{}) error {
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
								{"prefix": map[string]interface{}{
									"login.keyword": map[string]interface{}{
										"value":            lowerQuery,
										"case_insensitive": true,
									},
								}},
							},
							"must_not": prefixMustNot,
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
								{"term": map[string]interface{}{"_index": "entities"}},
								{"term": map[string]interface{}{"name.keyword": lowerQuery}},
							},
							"boost": 2.0,
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"bool": map[string]interface{}{
									"should": []map[string]interface{}{
										{"term": map[string]interface{}{"_index": "collections"}},
										{"term": map[string]interface{}{"_index": "tags"}},
										{"term": map[string]interface{}{"_index": "collection_items"}},
										{"term": map[string]interface{}{"_index": "entities"}},
									},
								}},
								{"match": map[string]interface{}{
									"name.full": map[string]interface{}{
										"query":    lowerQuery,
										"operator": "or",
									},
								}},
							},
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"bool": map[string]interface{}{
									"should": []map[string]interface{}{
										{"term": map[string]interface{}{"_index": "collections"}},
										{"term": map[string]interface{}{"_index": "tags"}},
										{"term": map[string]interface{}{"_index": "collection_items"}},
										{"term": map[string]interface{}{"_index": "entities"}},
									},
								}},
								{"prefix": map[string]interface{}{
									"name.keyword": map[string]interface{}{
										"value":            lowerQuery,
										"case_insensitive": true,
									},
								}},
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

func (es *ElasticService) search(ctx context.Context, indices []string, query map[string]interface{}) (*SearchResult, error) {
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
		Refresh:    "true", // Ждем обновления индекса
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
