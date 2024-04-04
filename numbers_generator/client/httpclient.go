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

func (h *HTTPClient) GenerateWinningNumbers() *types.WinningNumbers {
	endpoint := fmt.Sprintf("%s/winningNumbers", h.Endpoint)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	var winningNumbers types.WinningNumbers
	err = json.NewDecoder(resp.Body).Decode(&winningNumbers)
	if err != nil {
		return nil
	}

	return &winningNumbers
}
