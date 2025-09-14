package flux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	aspect_ratio_9_16 = "9:16"
)

type FluxClient struct {
	apiKey string
	client *http.Client
}

// ProcessingResponse represents the response when an image is being processed
type ProcessingResponse struct {
	Status      string                 `json:"status"`
	Tip         string                 `json:"tip"`
	ETA         int                    `json:"eta"`
	Message     string                 `json:"message"`
	FetchResult string                 `json:"fetch_result"`
	ID          int64                  `json:"id"`
	Output      []string               `json:"output"`
	Meta        map[string]interface{} `json:"meta"`
	FutureLinks []string               `json:"future_links"`
}

// GenerationResult holds either the processing response or error
type GenerationResult struct {
	IsProcessing bool
	Processing   *ProcessingResponse
	RawResponse  map[string]interface{}
}

func NewFluxClient(apiKey string) *FluxClient {
	return &FluxClient{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *FluxClient) GenerateImage(
	initImageURL,
	style,
	roomType string,
) (*GenerationResult, error) {
	payload := map[string]interface{}{
		"init_image":   initImageURL,
		"prompt":       c.getPrompt(style, roomType),
		"model_id":     "flux-kontext-pro",
		"aspect_ratio": aspect_ratio_9_16,
		"key":          c.apiKey,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://modelslab.com/api/v7/images/image-to-image", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var rawResult map[string]interface{}
	if err := json.Unmarshal(body, &rawResult); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := &GenerationResult{
		RawResponse: rawResult,
	}

	// Check if this is a processing response
	if status, ok := rawResult["status"].(string); ok && status == "processing" {
		var processingResp ProcessingResponse
		if err := json.Unmarshal(body, &processingResp); err != nil {
			return nil, fmt.Errorf("failed to parse processing response: %w", err)
		}
		result.IsProcessing = true
		result.Processing = &processingResp
	}

	return result, nil
}

// FetchResult fetches the final result using the fetch_result URL from processing response
func (c *FluxClient) FetchResult(fetchURL string) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", fetchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create fetch request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send fetch request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read fetch response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse fetch response: %w", err)
	}

	return result, nil
}

// Helper method to wait for and fetch the final result
func (c *FluxClient) GenerateImageAndWait(
	initImageURL,
	style,
	roomType string,
) (string, error) {
	result, err := c.GenerateImage(initImageURL, style, roomType)
	if err != nil {
		return "", err
	}

	if !result.IsProcessing {
		if len(result.Processing.FutureLinks) == 0 {
			return "", fmt.Errorf("no future links found in response")
		}

		return result.Processing.FutureLinks[0], nil
	}

	waitTime := time.Duration(result.Processing.ETA) * time.Second

	time.Sleep(waitTime)

	if len(result.Processing.FutureLinks) == 0 {
		return "", fmt.Errorf("no future links found in response")
	}

	return result.Processing.FutureLinks[0], nil
}
