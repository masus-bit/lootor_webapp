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
	Body        interface{}
	File        []byte
	BasicAuth   *struct {
		Username string
		Password string
	}
}

func SendRequest[T any](options RequestOptions) (T, error) {
	var result T

	queryParams := url.Values{}
	var reqBody io.Reader

	var fileBody bytes.Buffer
	writer := multipart.NewWriter(&fileBody)
	part, err := writer.CreateFormFile("file", uuid.New().String())
	if err != nil {
		return result, err
	}

	modifiedURL := options.URL

	if len(options.QueryParams) > 0 {
		for key, value := range options.QueryParams {
			queryParams.Add(key, value)
		}
		modifiedURL = options.URL + "?" + queryParams.Encode()
	}

	switch body := options.Body.(type) {
	case nil:
	case map[string]string:
		values := url.Values{}
		for key, value := range body {
			values.Add(key, value)
		}
		reqBody = strings.NewReader(values.Encode())
	case []byte:
		reqBody = bytes.NewReader(body)
	default:
		jsonData, err := json.Marshal(body)
		if err != nil {
			return result, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	var req *http.Request

	if len(options.File) > 0 {
		_, err = io.Copy(part, bytes.NewReader(options.File))
		if err != nil {
			return result, fmt.Errorf("failed to copy file: %w", err)
		}

		if err := writer.Close(); err != nil {
			return result, fmt.Errorf("failed to close multipart writer: %w", err)
		}

		req, err = http.NewRequest(options.Method, modifiedURL, &fileBody)
		if err != nil {
			return result, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		req, err = http.NewRequest(options.Method, modifiedURL, reqBody)
		if err != nil {
			return result, fmt.Errorf("failed to create request: %w", err)
		}
	}
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	if options.BasicAuth != nil {
		req.SetBasicAuth(options.BasicAuth.Username, options.BasicAuth.Password)
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
