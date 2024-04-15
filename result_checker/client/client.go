package client

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
)

type Client interface {
	IsTicketWinning(ctx context.Context, ticket *types.Ticket) (*types.ResultResponse, error)
}
