package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	winningTicketsColl = "winningTickets"
)

type MongoWinningTicketStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewWinningTicketStore(client *mongo.Client) *MongoWinningTicketStore {
	return &MongoWinningTicketStore{
		client: client,
		coll:   client.Database("lottery").Collection(winningTicketsColl),
	}
}

func (s *MongoWinningTicketStore) InsertWinningTickets(winningTickets []*types.Ticket) error {
	if len(winningTickets) == 0 {
		logrus.Info("no winning tickets to insert")
		return fmt.Errorf("no winning tickets to insert")
	}
	var tickets []interface{}
	for _, t := range winningTickets {
		tickets = append(tickets, t)
	}

	_, err := s.coll.InsertMany(context.Background(), tickets)
	if err != nil {
		return err
	}
	logrus.Info("inserted " + fmt.Sprint(len(winningTickets)) + " tickets")
	return nil
}

func (s *MongoWinningTicketStore) CheckIfTicketIsWinning(hash string) bool {
	filter := bson.M{"hash": hash}

	res, err := s.coll.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}

	return res > 0
}
