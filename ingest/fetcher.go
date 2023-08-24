package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cometbft/cometbft/rpc/grpc/client"
	"github.com/cometbft/cometbft/rpc/grpc/client/privileged"
	"github.com/cometbft/rpc-companion/config"
)

var (
	requestDefaultTimeout = 10 * time.Second
)

type Fetcher struct {
	BaseService
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

// GetBlock returns block at a specific height
func (f *Fetcher) GetBlock(height int64) (*client.Block, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
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

// GetBlockResults returns block results at a specific height
func (f *Fetcher) GetBlockResults(height int64) (*client.BlockResults, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
	defer cancel()

	conn, err := client.New(ctx, f.config.GRPCClient.ListenAddress, client.WithBlockServiceEnabled(true), client.WithInsecure()) //TODO: In the future support secure connections
	if err != nil {
		return nil, fmt.Errorf("error creating a gRPC client: %v", err)
	}
	defer conn.Close()

	//// Get Block By Height
	blockResults, err := conn.GetBlockResults(ctx, height)
	if err != nil {
		return nil, fmt.Errorf("error getting block results: %v", err)
	}
	f.logger.Debug("got block results at height: %d", height)
	return blockResults, nil
}

// GetBlockRetainHeight Get Block Retain Height value
func (f *Fetcher) GetBlockRetainHeight() (privileged.RetainHeights, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
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

// GetBlockResultsRetainHeight Get Block Retain Height value
func (f *Fetcher) GetBlockResultsRetainHeight() (uint64, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
	defer cancel()

	// Privileged Services
	conn, err := privileged.New(ctx, f.config.GRPCClient.ListenAddressPrivileged, privileged.WithPruningServiceEnabled(true), privileged.WithInsecure())
	defer conn.Close()
	if err != nil {
		fmt.Printf("error new priviledge client: %v", err)
	}

	retainHeight, err := conn.GetBlockResultsRetainHeight(ctx)
	if err != nil {
		f.logger.Error("GetBlockResultsRetainHeight request", "msg", "error getting the block results retain height", "error", err)
		return 0, fmt.Errorf("error getting the block results retain height")
	}
	return retainHeight, nil
}

// SetBlockResultsRetainHeight Set Block Results Retain Height value
func (f *Fetcher) SetBlockResultsRetainHeight(height uint64) error {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
	defer cancel()

	// Privileged Services
	conn, err := privileged.New(ctx, f.config.GRPCClient.ListenAddressPrivileged, privileged.WithPruningServiceEnabled(true), privileged.WithInsecure())
	defer conn.Close()
	if err != nil {
		fmt.Printf("error new priviledge client: %v\n", err)
	}

	err = conn.SetBlockResultsRetainHeight(ctx, height)
	if err != nil {
		f.logger.Error("SetBlockResultsRetainHeight request", "msg", "error setting the block results retain height", "error", err)
		return fmt.Errorf("error setting the block results retain height")
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

	// Check Block Results Retain Height
	h, err := f.GetBlockResultsRetainHeight()
	if err != nil {
		f.logger.Error("fetcher request", "error", err)
	} else {
		f.logger.Info("fetcher request", "msg", "Got pruning block results retain height", "block_results_retain_height", h)
	}

	// Get Block Results
	br, err := f.GetBlockResults(int64(h + 1))
	if err != nil {
		f.logger.Error("fetcher request", "error", err)
	} else {
		f.logger.Info("fetcher request", "msg", "Got block results", "height", br.Height)
	}

	// Set Block Results Retain Height
	err = f.SetBlockResultsRetainHeight(uint64(br.Height + 1))
	if err != nil {
		f.logger.Error("fetcher request", "error", err)
	} else {
		f.logger.Info("fetcher request", "msg", "Set pruning block results retain height", "block_retain_height", br.Height+1)
	}

	return nil
}
