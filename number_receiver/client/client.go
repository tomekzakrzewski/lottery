package client

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
)

type Client interface {
	GetTicketByHash(ctx context.Context, hash string) (*types.Ticket, error)
	GetNextDrawDate(ctx context.Context) *types.DrawDate
	CreateTicket(nctx context.Context, ums *types.UserNumbers) (*types.Ticket, error)
}
