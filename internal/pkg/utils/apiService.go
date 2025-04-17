package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// RequestOptions - параметры запроса
type RequestOptions struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        map[string]string
	File        []byte
}

func SendRequest[T any](options RequestOptions) (T, error) {
	var result T

	queryParams := url.Values{}
	requestBody := url.Values{}

	var fileBody bytes.Buffer
	writer := multipart.NewWriter(&fileBody)
	part, err := writer.CreateFormFile("file", uuid.New().String())
	if err != nil {
		return result, err
	}

	modifiedUrl := options.URL

	if len(options.QueryParams) > 0 {
		for key, value := range options.QueryParams {
			queryParams.Add(key, value)
		}
		modifiedUrl = options.URL + "?" + queryParams.Encode()
	}

	if len(options.Body) > 0 {
		for key, value := range options.Body {
			requestBody.Add(key, value)
		}
	}

	var req *http.Request

	if len(options.File) > 0 {
		_, err = io.Copy(part, bytes.NewReader(options.File))
		writer.Close()
		req, err = http.NewRequest(options.Method, modifiedUrl, &fileBody)
		if err != nil {
			return result, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		req, err = http.NewRequest(options.Method, modifiedUrl, strings.NewReader(requestBody.Encode()))
		if err != nil {
			return result, fmt.Errorf("failed to create request: %w", err)
		}
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
