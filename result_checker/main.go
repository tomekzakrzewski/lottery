package main

import (
	"context"
	"time"

	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	generator "github.com/tomekzakrzewski/lottery/numbers_generator/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	receiverHTTPClient := receiver.NewHTTPClient("http://localhost:3000")
	generatorHTTPClient := generator.NewHTTPClient("http://localhost:3001")

	store := NewWinningTicketStore(client)
	svc := NewResultCheckerService(*receiverHTTPClient, *generatorHTTPClient, *store)
	m := NewLogMiddleware(svc)
	time.Sleep(5 * time.Second)
	m.GetWinningTickets()

	time.Sleep(5 * time.Minute)
}
