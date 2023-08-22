package fetcher

import (
	"context"
	"fmt"
	"github.com/cometbft/cometbft/rpc/grpc/client"
	service "github.com/cometbft/rpc-companion/service/base"
	"log/slog"
	"time"
)

type Fetcher struct {
	service.BaseService
	endpoint string
	logger   slog.Logger
}

func NewFetcher(listenAddress string, logger slog.Logger) (*Fetcher, error) {
	return &Fetcher{
		endpoint: listenAddress,
		logger:   logger,
	}, nil
}

func (f *Fetcher) GetBlock(height int64) (*client.Block, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	conn, err := client.New(ctx, f.endpoint, client.WithInsecure()) //TODO: In the future support secure connections
	defer conn.Close()
	if err != nil {
		return nil, fmt.Errorf("error creating a gRPC client: %v", err)
	}

	//// Get Block
	block, err := conn.GetBlockByHeight(ctx, height)
	if err != nil {
		return nil, fmt.Errorf("error getting block: %v\n", err)
	}
	f.logger.Debug("got block at height: %d", height)
	return block, nil
}

func (f *Fetcher) OnStart() error {
	go func() {
		f.logger.Info("service running", "service", "fetcher", "msg", fmt.Sprintf("Fetcher service is running"))
	}()
	return nil
}
