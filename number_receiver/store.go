package main

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ticketColl = "tickets"
)

type MongoTicketStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewTicketStore(client *mongo.Client) *MongoTicketStore {
	return &MongoTicketStore{
		client: client,
		coll:   client.Database("lottery").Collection(ticketColl),
	}
}

func (s *MongoTicketStore) Insert(ticket *types.Ticket) (*types.Ticket, error) {
	res, err := s.coll.InsertOne(context.TODO(), ticket)
	if err != nil {
		return nil, err
	}
	ticket.ID = res.InsertedID.(primitive.ObjectID)
	return ticket, nil
}

func (s *MongoTicketStore) FindByHash(hash string) (*types.Ticket, error) {
	var ticket types.Ticket
	err := s.coll.FindOne(context.TODO(), bson.M{"hash": hash}).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// fix getting tickets by draw date and winning numbers
func (s *MongoTicketStore) FindByWinningNumersAndDrawDate(winningNumbers types.WinningNumbers) ([]*types.Ticket, error) {
	res, err := s.coll.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var tickets []*types.Ticket
	err = res.All(context.Background(), &tickets)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}
