package types

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Numbers  []int              `json:"numbers" bson:"numbers"`
	DrawDate time.Time          `json:"drawDate" bson:"drawDate"`
	Hash     string             `json:"hash" bson:"hash"`
}

type CreateTicketParams struct {
	Numbers  []int     `json:"numbers"`
	DrawDate time.Time `json:"drawDate"`
}

func NewTicketFromParams(params *CreateTicketParams) *Ticket {
	return &Ticket{
		Numbers:  params.Numbers,
		Hash:     uuid.New().String(),
		DrawDate: params.DrawDate,
	}
}
