package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (h *HTTPClient) IsTicketWinning(ctx context.Context, hash string) *bool {
	endpoint := fmt.Sprintf("%s/win/%s", h.Endpoint, hash)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	var isWinning bool
	err = json.NewDecoder(resp.Body).Decode(&isWinning)
	if err != nil {
		return nil
	}

	return &isWinning
}
