package main

import (
	"context"

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
	_ = NewResultCheckerService(*receiverHTTPClient, *generatorHTTPClient, *store)
}
