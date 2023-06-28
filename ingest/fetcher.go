package ingest

import (
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type Fetcher interface {
	FetchBlock(height int64) (*ctypes.ResultBlock, error)
	FetchABCIInfo() (*ctypes.ResultABCIInfo, error)
}
