package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

type ElasticService struct {
	client *elasticsearch.Client
	config ElasticConfig
	logger *log.Logger
}

func LoadConfig(path string) (ElasticConfig, error) {
	var config ElasticConfig

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return config, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}

func NewElasticService(configPath string) (*ElasticService, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Логируем загруженные настройки
	log.Printf("Elasticsearch config loaded: Address=%s, Indices=%v",
		config.Address, getIndexNames(config.Indices))

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.Address},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create elastic client: %w", err)
	}

	logger := log.New(os.Stdout, "[Elastic] ", log.LstdFlags|log.Lshortfile)

	// Проверка соединения
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

func getIndexNames(indices map[string]IndexConfig) []string {
	names := make([]string, 0, len(indices))
	for name := range indices {
		names = append(names, name)
	}
	return names
}

func (es *ElasticService) Search(ctx context.Context, index string, query map[string]interface{}) (*SearchResult, error) {
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}
	res, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex(index),
		es.client.Search.WithBody(strings.NewReader(string(queryJSON))))
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, parseErrorResponse(res)
	}

	var result SearchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}

func (es *ElasticService) SearchInIndices(ctx context.Context, indices []string, query string) (*SearchResult, error) {
	if strings.TrimSpace(query) == "" {
		return &SearchResult{}, nil
	}

	lowerQuery := strings.ToLower(query)

	// Полный поисковый запрос как в NestJS
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"_index": "users"}},
								{
									"bool": map[string]interface{}{
										"should": []map[string]interface{}{
											{
												"term": map[string]interface{}{
													"login.keyword": map[string]interface{}{
														"value": lowerQuery,
													},
												},
											},
											{
												"match": map[string]interface{}{
													"login.prefix": map[string]interface{}{
														"query": lowerQuery,
													},
												},
											},
											{
												"match": map[string]interface{}{
													"user_name": map[string]interface{}{
														"query": lowerQuery,
													},
												},
											},
											{
												"wildcard": map[string]interface{}{
													"login.keyword": map[string]interface{}{
														"value":            "*" + lowerQuery + "*",
														"case_insensitive": true,
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						"bool": map[string]interface{}{
							"must_not": map[string]interface{}{"term": map[string]interface{}{"_index": "users"}},
							"should": []map[string]interface{}{
								{
									"term": map[string]interface{}{
										"name.keyword": map[string]interface{}{
											"value": lowerQuery,
										},
									},
								},
								{
									"match": map[string]interface{}{
										"name.prefix": map[string]interface{}{
											"query": lowerQuery,
										},
									},
								},
								{
									"match": map[string]interface{}{
										"name.full": map[string]interface{}{
											"query": lowerQuery,
										},
									},
								},
								{
									"wildcard": map[string]interface{}{
										"name.keyword": map[string]interface{}{
											"value":            "*" + lowerQuery + "*",
											"case_insensitive": true,
										},
									},
								},
							},
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
		"size": 100,
	}

	return es.Search(ctx, strings.Join(indices, ","), searchQuery)
}

func (es *ElasticService) IndexDocument(ctx context.Context, index string, body map[string]interface{}) error {
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return fmt.Errorf("error encoding document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:   index,
		Body:    strings.NewReader(buf.String()),
		Refresh: "wait_for",
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

func (es *ElasticService) CreateIndexIfNotExists(ctx context.Context, indexName string) error {
	cfg, exists := es.config.Indices[indexName]
	if !exists {
		return fmt.Errorf("configuration for index %q not found", indexName)
	}

	// Проверяем существование индекса
	res, err := es.client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("error checking index existence: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil // Индекс уже существует
	}

	// Создаем индекс
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(cfg); err != nil {
		return fmt.Errorf("error encoding index config: %w", err)
	}

	createRes, err := es.client.Indices.Create(
		indexName,
		es.client.Indices.Create.WithBody(strings.NewReader(buf.String())),
		es.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		return parseErrorResponse(createRes)
	}

	es.logger.Printf("Index %q created successfully", indexName)
	return nil
}

func (es *ElasticService) UpsertDocument(ctx context.Context, index, id string, body map[string]interface{}) error {
	var buf strings.Builder
	updateBody := map[string]interface{}{
		"doc":           body,
		"doc_as_upsert": true,
	}

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

func (es *ElasticService) SafeIndexDocument(ctx context.Context, index, id string, body map[string]interface{}) error {
	exists, err := es.DocumentExists(ctx, index, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("document with id %q already exists", id)
	}

	return es.IndexDocument(ctx, index, body)
}

func (es *ElasticService) DocumentExists(ctx context.Context, index, id string) (bool, error) {
	req := esapi.ExistsRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return false, fmt.Errorf("exists request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return true, nil
	} else if res.StatusCode == 404 {
		return false, nil
	}

	return false, parseErrorResponse(res)
}

func (es *ElasticService) DeleteByQuery(ctx context.Context, index string, query map[string]interface{}) error {
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return fmt.Errorf("error encoding query: %w", err)
	}
	refresh := true
	req := esapi.DeleteByQueryRequest{
		Index:   []string{index},
		Body:    strings.NewReader(buf.String()),
		Refresh: &refresh, // Используем указатель на bool,
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("delete by query error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseErrorResponse(res)
	}

	return nil
}

func (es *ElasticService) UpdateDocument(ctx context.Context, index, id string, body map[string]interface{}) error {
	var buf strings.Builder
	updateBody := map[string]interface{}{
		"doc": body,
	}

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

func (es *ElasticService) DeleteDocument(ctx context.Context, index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "wait_for",
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("delete request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil // Документ не найден - считаем успехом
	} else if res.IsError() {
		return parseErrorResponse(res)
	}

	return nil
}

func (es *ElasticService) ReindexAll(ctx context.Context, dataProviders map[string]func() ([]map[string]interface{}, error)) error {
	// Удаляем все индексы
	for indexName := range es.config.Indices {
		res, err := es.client.Indices.Delete([]string{indexName})
		if err != nil {
			return fmt.Errorf("error deleting index %q: %w", indexName, err)
		}
		res.Body.Close()
		es.logger.Printf("Deleted index: %s", indexName)
	}

	// Создаем индексы заново
	for indexName := range es.config.Indices {
		if err := es.CreateIndexIfNotExists(ctx, indexName); err != nil {
			return fmt.Errorf("error recreating index %q: %w", indexName, err)
		}
	}

	// Индексируем данные
	for indexName, provider := range dataProviders {
		data, err := provider()
		if err != nil {
			return fmt.Errorf("error getting data for %q: %w", indexName, err)
		}

		var bulkBody strings.Builder
		for _, item := range data {
			id, ok := item["id"].(string)
			if !ok {
				return fmt.Errorf("missing or invalid id in document for index %q", indexName)
			}

			meta := map[string]interface{}{
				"index": map[string]interface{}{
					"_index": indexName,
					"_id":    id,
				},
			}

			metaJSON, _ := json.Marshal(meta)
			docJSON, _ := json.Marshal(item)

			bulkBody.WriteString(string(metaJSON) + "\n")
			bulkBody.WriteString(string(docJSON) + "\n")
		}

		res, err := es.client.Bulk(
			strings.NewReader(bulkBody.String()),
			es.client.Bulk.WithContext(ctx),
			es.client.Bulk.WithRefresh("true"),
		)
		if err != nil {
			return fmt.Errorf("bulk index error for %q: %w", indexName, err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return parseErrorResponse(res)
		}

		es.logger.Printf("Successfully indexed %d documents in %q", len(data), indexName)
	}

	return nil
}

func parseErrorResponse(res *esapi.Response) error {
	var e map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		return fmt.Errorf("error parsing error response: %w", err)
	}

	return fmt.Errorf("elasticsearch error [%d]: %v", res.StatusCode, e["error"])
}
