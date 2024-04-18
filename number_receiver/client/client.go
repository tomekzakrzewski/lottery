package client

import (
	"context"
	"time"

	"github.com/tomekzakrzewski/lottery/types"
)

type Client interface {
	GetTicketByHash(ctx context.Context, hash string) (*types.Ticket, error)
	GetNextDrawDate(ctx context.Context, currentTime time.Time) (*types.DrawDate, error)
	CreateTicket(ctx context.Context, nums *types.UserNumbers) (*types.Ticket, error)
}
