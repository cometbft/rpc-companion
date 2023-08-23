package fetcher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cometbft/cometbft/rpc/grpc/client"
	"github.com/cometbft/cometbft/rpc/grpc/client/privileged"
	"github.com/cometbft/rpc-companion/config"
	service "github.com/cometbft/rpc-companion/service/base"
)

type Fetcher struct {
	service.BaseService
	config *config.Config
	logger slog.Logger
}

func NewFetcher(logger slog.Logger, cfg *config.Config) (*Fetcher, error) {
	return &Fetcher{
		logger: logger,
		config: cfg,
	}, nil
}

//----------------------------------------------------------------------------------------------------------------------
// Requests

func (f *Fetcher) GetBlock(height int64) (*client.Block, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	conn, err := client.New(ctx, f.config.GRPCClient.ListenAddress, client.WithBlockServiceEnabled(true), client.WithInsecure()) //TODO: In the future support secure connections
	if err != nil {
		return nil, fmt.Errorf("error creating a gRPC client: %v", err)
	}
	defer conn.Close()

	//// Get Block By Height
	block, err := conn.GetBlockByHeight(ctx, height)
	if err != nil {
		return nil, fmt.Errorf("error getting block: %v", err)
	}
	f.logger.Debug("got block at height: %d", height)
	return block, nil
}

// GetBlockRetainHeight Get Block Retain Height value
func (f *Fetcher) GetBlockRetainHeight() (privileged.RetainHeights, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	// Privileged Services
	conn, err := privileged.New(ctx, f.config.GRPCClient.ListenAddressPrivileged, privileged.WithPruningServiceEnabled(true), privileged.WithInsecure())
	defer conn.Close()
	if err != nil {
		fmt.Printf("error new priviledge client: %v", err)
	}

	retainHeight, err := conn.GetBlockRetainHeight(ctx)
	if err != nil {
		f.logger.Error("GetBlockRetainHeight request", "msg", "error getting the block retain height", "error", err)
		return privileged.RetainHeights{
			App:            0,
			PruningService: 0,
		}, fmt.Errorf("error getting the block retain height")
	}
	return retainHeight, nil
}

// SetBlockRetainHeight Set Block Retain Height value
func (f *Fetcher) SetBlockRetainHeight(height uint64) error {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	// Privileged Services
	conn, err := privileged.New(ctx, f.config.GRPCClient.ListenAddressPrivileged, privileged.WithPruningServiceEnabled(true), privileged.WithInsecure())
	defer conn.Close()
	if err != nil {
		fmt.Printf("error new priviledge client: %v\n", err)
	}

	err = conn.SetBlockRetainHeight(ctx, height)
	if err != nil {
		f.logger.Error("SetBlockRetainHeight request", "msg", "error setting the block retain height", "error", err)
		return fmt.Errorf("error setting the block retain height")
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
// Services methods

func (f *Fetcher) OnStart() error {

	f.logger.Info("service running", "msg", fmt.Sprintf("Fetcher service is running"))

	// Get latest block retain height
	rh, err := f.GetBlockRetainHeight()
	if err != nil {
		f.logger.Error("fetcher request", "error", err)
	} else {
		f.logger.Info("fetcher request", "msg", "Got pruning block retain height", "block_retain_height", rh.PruningService, "app_retain_height", rh.App)
	}

	// Fetch a block that is one block higher than the lowest block retain height
	height := min(rh.PruningService, rh.App) + 1

	b, err := f.GetBlock(int64(height))
	if err != nil {
		f.logger.Error("fetcher request", "error", err)
	} else {
		f.logger.Info("fetcher request", "msg", "Got block", "height", b.Block.Height)
	}

	// Set Block Retain Height if it's higher than zero
	if height > rh.PruningService {
		err = f.SetBlockRetainHeight(height)
		if err != nil {
			f.logger.Error("fetcher request", "error", err)
		} else {
			f.logger.Info("fetcher request", "msg", "Set pruning block retain height", "block_retain_height", height)
		}
	} else {
		f.logger.Info("fetcher request", "msg", "Skipping set pruning block retain height. Block retain height is higher than the app retain height", "block_retain_height", rh.PruningService, "app_retain_height", rh.App)
	}

	// Check Block Retain Height again
	rh, err = f.GetBlockRetainHeight()
	if err != nil {
		f.logger.Error("fetcher request", "error", err)
	} else {
		f.logger.Info("fetcher request", "msg", "Got pruning block retain height", "block_retain_height", rh.PruningService, "app_retain_height", rh.App)
	}

	return nil
}
