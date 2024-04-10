package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HttpTransport struct {
	svc ResultAnnoucer
}

func NewHttpTransport(svc ResultAnnoucer) *HttpTransport {
	return &HttpTransport{
		svc: svc,
	}
}

func (h *HttpTransport) handleCheckResult(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	isWinning, _ := h.svc.CheckResult(hash)
	if isWinning == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, isWinning, nil)
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
