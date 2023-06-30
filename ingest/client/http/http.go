package http

import (
	"context"
	"github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type HTTP struct {
	Endpoint string
}

func (c *HTTP) Header(height int64) (*ctypes.ResultHeader, error) {

	client, err := http.New(c.Endpoint, "/websocket")
	if err != nil {
		return nil, err
	}

	header, err := client.Header(context.Background(), &height)
	if err != nil {
		return nil, err
	} else {
		return header, nil
	}
}
