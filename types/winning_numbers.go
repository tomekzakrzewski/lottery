package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WinningNumbers struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Numbers  []int              `json:"numbers" bson:"numbers"`
	DrawDate time.Time          `json:"drawDate" bson:"drawDate"`
}
