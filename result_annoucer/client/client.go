package client

import "github.com/tomekzakrzewski/lottery/types"

type Client interface {
	CheckResult(hash string) (*types.ResultResponse, error)
}
