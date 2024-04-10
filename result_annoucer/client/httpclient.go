package client

import (
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

func (h *HTTPClient) CheckReuslt(hash string) *types.ResultRespone {
	endpoint := fmt.Sprintf("%s/win/%s", h.Endpoint, hash)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	var result types.ResultRespone
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil
	}

	return &result
}
