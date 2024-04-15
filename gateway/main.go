package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	annoucer "github.com/tomekzakrzewski/lottery/result_annoucer/client"
	"github.com/tomekzakrzewski/lottery/types"
)

func main() {
	listenAddr := ":8000"
	//receiverAddr := "http://localhost:3000"
	receiverGRPC := "localhost:3006"
	//annoucerAddr := "http://localhost:6000"
	annoucerGRPC := "localhost:6006"

	//receiverHttpClient := receiver.NewHTTPClient(receiverAddr)
	receiverGRPCClient, err := receiver.NewGRPCClient(receiverGRPC)
	if err != nil {
		panic(err)
	}
	//	annoucerHttpClient := annoucer.NewHTTPClient(annoucerAddr)
	annoucerGRPCClient, err := annoucer.NewGRPCClient(annoucerGRPC)
	if err != nil {
		panic(err)
	}

	handler := NewHandler(*receiverGRPCClient, annoucerGRPCClient)
	r := chi.NewRouter()

	r.Post("/inputTicket", handler.handlePostTicket)
	r.Get("/result/{hash}", handler.handleCheckResult)
	http.ListenAndServe(listenAddr, r)
}

type Handler struct {
	receiverClient receiver.GRPCClient
	annoucerClient annoucer.Client
}

func NewHandler(receiverClient receiver.GRPCClient, annoucerClient annoucer.Client) *Handler {
	return &Handler{
		receiverClient: receiverClient,
		annoucerClient: annoucerClient,
	}
}

func (h *Handler) handlePostTicket(w http.ResponseWriter, r *http.Request) {
	var body *types.UserNumbers

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ticket, err := h.receiverClient.CreateTicket(context.Background(), body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, ticket, nil)
}

func (h *Handler) handleCheckResult(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	isWinning, err := h.annoucerClient.CheckResult(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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