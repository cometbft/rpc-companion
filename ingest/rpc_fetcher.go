package ingest

import (
	"context"
	client "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type RPCFetcher struct {
	Endpoint string
}

func (c *RPCFetcher) FetchBlock(height int64) (*ctypes.ResultBlock, error) {

	httpClient, err := client.New(c.Endpoint, "/websocket")
	if err != nil {
		return nil, err
	}

	resultBlock, err := httpClient.Block(context.Background(), &height)
	if err != nil {
		return nil, err
	} else {
		return resultBlock, nil
	}
}

func (c *RPCFetcher) FetchABCIInfo() (*ctypes.ResultABCIInfo, error) {

	httpClient, err := client.New(c.Endpoint, "/websocket")
	if err != nil {
		return nil, err
	}

	abciInfo, err := httpClient.ABCIInfo(context.Background())
	if err != nil {
		return nil, err
	} else {
		return abciInfo, nil
	}
}
