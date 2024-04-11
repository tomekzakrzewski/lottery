package client

import (
	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	Endpoint string
	client   types.GeneratorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := types.NewGeneratorClient(conn)

	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}
