package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HttpTransport struct {
	svc ResultChecker
}

func NewHttpTransport(svc ResultChecker) *HttpTransport {
	return &HttpTransport{
		svc: svc,
	}
}

// can handle better, but for now it's ok. write better response. return numbers
func (h *HttpTransport) handleCheckIsTicketWinning(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	isWinning := h.svc.IsTicketWinning(hash)
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
