package client

import (
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type IClient interface {
	Header(height int64) (*ctypes.ResultHeader, error)
}
