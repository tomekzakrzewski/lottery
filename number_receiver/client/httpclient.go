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

func (h *HTTPClient) CreateTicket(ctx context.Context, nums *types.UserNumbers) (*types.Ticket, error) {
	endpoint := fmt.Sprintf("%s/ticket", h.Endpoint)

	body, err := json.Marshal(nums)
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

	var ticket types.Ticket
	err = json.NewDecoder(resp.Body).Decode(&ticket)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (h *HTTPClient) GetNextDrawDate(ctx context.Context) (*types.DrawDate, error) {
	endpoint := fmt.Sprintf("%s/drawDate", h.Endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var date types.DrawDate
	err = json.NewDecoder(resp.Body).Decode(&date)
	if err != nil {
		return nil, err
	}

	return &date, nil
}
