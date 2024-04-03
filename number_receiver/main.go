package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	store := NewTicketStore(client)
	svc := NewNumberReceiver(store)
	m := NewLogMiddleware(svc)

	p := &types.UserNumbers{
		Numbers: []int{1, 2, 3, 4, 5, 6},
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	srv := NewSrv(m)

	r := chi.NewRouter()
	r.Get("/drawDate", srv.handleGetDrawDate)
	http.ListenAndServe(":3000", r)

	m.CreateTicket(p)
	m.NextDrawDate()
}

type Srv struct {
	svc NumberReceiver
}

func NewSrv(svc NumberReceiver) *Srv {
	return &Srv{
		svc: svc,
	}
}

func (s *Srv) handleGetDrawDate(w http.ResponseWriter, r *http.Request) {
	date := s.svc.NextDrawDate()
	writeJSON(w, http.StatusOK, date, nil)
}

func writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
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
