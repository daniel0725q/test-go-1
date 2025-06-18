package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/truora/microservice/internal/dto"
)

type ExternalAPIRepository interface {
	GetHello(ctx context.Context) ([]*dto.StockRatingResponse, error)
}

type externalAPIRepository struct {
	baseURL string
	timeout time.Duration
	token   string
}

func NewExternalAPIRepository(baseURL string, timeout time.Duration, token string) ExternalAPIRepository {
	return &externalAPIRepository{
		baseURL: baseURL,
		timeout: timeout,
		token:   token,
	}
}

func (r *externalAPIRepository) GetHello(ctx context.Context) ([]*dto.StockRatingResponse, error) {
	var allItems []*dto.StockRatingResponse
	nextPage := ""

	for {
		url := fmt.Sprintf("%s/production/swechallenge/list", r.baseURL)
		if nextPage != "" {
			url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
		}

		client := &http.Client{
			Timeout: r.timeout,
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.token))

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response dto.Response
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allItems = append(allItems, response.Items...)

		// Check if there are more pages
		if response.NextPage == "" {
			break
		}
		nextPage = response.NextPage
	}

	return allItems, nil
}
