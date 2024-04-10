package main

import (
	"context"
	"time"

	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	winningNumbersColl = "winningNumbers"
)

type MongoNumbersStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewNumbersStore(client *mongo.Client) *MongoNumbersStore {
	return &MongoNumbersStore{
		client: client,
		coll:   client.Database("lottery").Collection(winningNumbersColl),
	}
}

func (s *MongoNumbersStore) InsertWinningNumbers(numbers *types.WinningNumbers) (*types.WinningNumbers, error) {
	res, err := s.coll.InsertOne(context.Background(), numbers)
	if err != nil {
		return nil, err
	}
	numbers.ID = res.InsertedID.(primitive.ObjectID)
	return numbers, nil
}

func (s *MongoNumbersStore) FindWinningNumbersByDate(drawDate time.Time) (*types.WinningNumbers, error) {
	var winningNumbers types.WinningNumbers
	filter := bson.M{"drawDate": drawDate}
	err := s.coll.FindOne(context.Background(), filter).Decode(&winningNumbers)
	if err != nil {
		return nil, err
	}

	return &winningNumbers, nil
}

/*
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

*/
