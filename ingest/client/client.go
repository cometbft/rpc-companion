package client

import (
	"github.com/cometbft/cometbft/proto/tendermint/types"
)

type IClient interface {
	Header(height int64) (*types.Header, error)
}
