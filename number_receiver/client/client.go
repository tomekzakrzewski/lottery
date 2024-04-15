package client

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
)

type Client interface {
	GetTicketByHash(ctx context.Context, hash string) (*types.Ticket, error)
	GetNextDrawDate(ctx context.Context) *types.DrawDate
	CreateTicket(ctx context.Context, nums *types.UserNumbers) (*types.Ticket, error)
}
