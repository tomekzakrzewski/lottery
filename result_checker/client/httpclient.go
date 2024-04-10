package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tomekzakrzewski/lottery/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (h *HTTPClient) IsTicketWinning(ctx context.Context, ticket *types.Ticket) (*types.ResultResponse, error) {
	endpoint := fmt.Sprintf("%s/win", h.Endpoint)

	body, err := json.Marshal(ticket)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var resultResponse types.ResultResponse
	err = json.NewDecoder(resp.Body).Decode(&resultResponse)
	if err != nil {
		return nil, err
	}

	return &resultResponse, nil
}
