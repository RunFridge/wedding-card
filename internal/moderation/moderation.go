package moderation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const openAIModerationURL = "https://api.openai.com/v1/moderations"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type ModerationResult struct {
	Flagged        bool
	CategoryScores map[string]float64
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) ModerateText(ctx context.Context, text string) (*ModerationResult, error) {
	body := map[string]any{
		"model": "omni-moderation-latest",
		"input": text,
	}
	return c.doRequest(ctx, body)
}

func (c *Client) ModerateImage(ctx context.Context, imageURL string) (*ModerationResult, error) {
	body := map[string]any{
		"model": "omni-moderation-latest",
		"input": []any{
			map[string]any{
				"type": "image_url",
				"image_url": map[string]string{
					"url": imageURL,
				},
			},
		},
	}
	return c.doRequest(ctx, body)
}

func (c *Client) doRequest(ctx context.Context, body any) (*ModerationResult, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIModerationURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI moderation API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Results []struct {
			Flagged        bool               `json:"flagged"`
			CategoryScores map[string]float64 `json:"category_scores"`
		} `json:"results"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("OpenAI moderation API returned no results")
	}

	return &ModerationResult{
		Flagged:        result.Results[0].Flagged,
		CategoryScores: result.Results[0].CategoryScores,
	}, nil
}
