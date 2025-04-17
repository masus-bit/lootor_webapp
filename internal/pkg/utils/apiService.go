package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// RequestOptions - параметры запроса
type RequestOptions struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        interface{}
}

func SendRequest[T any](options RequestOptions) (T, error) {
	var result T

	var bodyReader io.Reader
	if options.Body != nil {
		bodyBytes, err := json.Marshal(options.Body)
		if err != nil {
			return result, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	queryParams := url.Values{}

	modifiedUrl := options.URL

	if len(options.QueryParams) > 0 {
		for key, value := range options.QueryParams {
			queryParams.Add(key, value)
		}
		modifiedUrl = options.URL + "?" + queryParams.Encode()
	}

	req, err := http.NewRequest(options.Method, modifiedUrl, bodyReader)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("bad status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response: %w", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}
