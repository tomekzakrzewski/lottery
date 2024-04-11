package client

import "github.com/tomekzakrzewski/lottery/types"

type Client interface {
	GenerateWinningNumbers() *types.WinningNumbers
}
