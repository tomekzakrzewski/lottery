package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tomekzakrzewski/lottery/types"
)

type HttpTransport struct {
	svc NumberReceiver
}

func NewHttpTransport(svc NumberReceiver) *HttpTransport {
	return &HttpTransport{
		svc: svc,
	}
}

func (h *HttpTransport) handlePostTicket(w http.ResponseWriter, r *http.Request) {
	var body types.UserNumbers
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ticket, err := h.svc.CreateTicket(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, ticket, nil)
}

func (h *HttpTransport) handleGetNextDrawDate(w http.ResponseWriter, r *http.Request) {
	drawDate := h.svc.NextDrawDate(time.Now())
	writeJSON(w, http.StatusOK, drawDate, nil)
}

func (h *HttpTransport) handleFindByHash(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	ticket, err := h.svc.GetTicketByHash(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, ticket, nil)
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
