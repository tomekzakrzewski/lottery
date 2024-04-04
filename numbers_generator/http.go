package main

import (
	"encoding/json"
	"net/http"
)

type HttpTransport struct {
	svc GeneratorServicer
}

func NewHttpTransport(svc GeneratorServicer) *HttpTransport {
	return &HttpTransport{
		svc: svc,
	}
}

func (h *HttpTransport) handleGetWinningNumbers(w http.ResponseWriter, r *http.Request) {
	winningNumbers := h.svc.GenerateWinningNumbers()
	writeJSON(w, http.StatusOK, winningNumbers, nil)
}

func writeJSON(w http.ResponseWriter, status int, v any, headers http.Header) error {
	js, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
