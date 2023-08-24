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
	logger = *logger.With("module", "Fetcher")
	return &Fetcher{
		logger: logger,
		config: cfg,
	}, nil
}

//----------------------------------------------------------------------------------------------------------------------
// Requests

// GetBlock returns block at a specific height
func (f *Fetcher) GetBlock(height int64) (*client.Block, error) {
	logger := *f.logger.With("method", "GetBlock")
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
	defer cancel()

	conn, err := client.New(ctx, f.config.GRPCClient.ListenAddress, client.WithBlockServiceEnabled(true), client.WithInsecure()) //TODO: In the future support secure connections
	if err != nil {
		logger.Error("New client", "error", err)
		return nil, fmt.Errorf("error creating new client")
	}
	defer conn.Close()

	//// Get Block By Height
	block, err := conn.GetBlockByHeight(ctx, height)
	if err != nil {
		logger.Error("Get block", "error", err, "height", height)
		return nil, fmt.Errorf("error getting block")
	}
	logger.Info("Get block", "height", height)
	return block, nil
}

// GetBlockResults returns block results at a specific height
func (f *Fetcher) GetBlockResults(height int64) (*client.BlockResults, error) {
	logger := *f.logger.With("method", "GetBlockResults")
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
	defer cancel()

	conn, err := client.New(ctx, f.config.GRPCClient.ListenAddress, client.WithBlockServiceEnabled(true), client.WithInsecure()) //TODO: In the future support secure connections
	if err != nil {
		logger.Error("New client", "error", err)
		return nil, fmt.Errorf("error creating new client")
	}
	defer conn.Close()

	//// Get Block Results By Height
	blockResults, err := conn.GetBlockResults(ctx, height)
	if err != nil {
		logger.Error("Get block results", "error", err)
		return nil, fmt.Errorf("error getting block results")
	}
	logger.Info("Get block results", "height", height)
	return blockResults, nil
}

// GetBlockRetainHeight Get Block Retain Height value
func (f *Fetcher) GetBlockRetainHeight() (privileged.RetainHeights, error) {
	logger := *f.logger.With("method", "GetBlockRetainHeight")
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
	defer cancel()

	// Privileged Services
	conn, err := privileged.New(ctx, f.config.GRPCClient.ListenAddressPrivileged, privileged.WithPruningServiceEnabled(true), privileged.WithInsecure())
	defer conn.Close()
	if err != nil {
		logger.Error("New privileged client", "error", err)
		return privileged.RetainHeights{
			App:            0,
			PruningService: 0,
		}, fmt.Errorf("error new priviledged client")
	}

	retainHeight, err := conn.GetBlockRetainHeight(ctx)
	if err != nil {
		logger.Error("Get block retain height", "error", err)
		return privileged.RetainHeights{
			App:            0,
			PruningService: 0,
		}, fmt.Errorf("error getting the block retain height")
	}
	logger.Info("Get block retain height", "retain_height", retainHeight.PruningService, "app_retain_height", retainHeight.App)
	return retainHeight, nil
}

// SetBlockRetainHeight Set Block Retain Height value
func (f *Fetcher) SetBlockRetainHeight(height uint64) error {
	logger := *f.logger.With("method", "SetBlockRetainHeight")
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
	defer cancel()

	// Privileged Services
	conn, err := privileged.New(ctx, f.config.GRPCClient.ListenAddressPrivileged, privileged.WithPruningServiceEnabled(true), privileged.WithInsecure())
	defer conn.Close()
	if err != nil {
		logger.Error("New privileged client", "error", err)
		return fmt.Errorf("error new priviledge client")
	}

	err = conn.SetBlockRetainHeight(ctx, height)
	if err != nil {
		logger.Error("Set block retain height", "error", err)
		return fmt.Errorf("error setting the block retain height")
	}
	logger.Info("Set block retain height", "height", height)
	return nil
}

// GetBlockResultsRetainHeight Get Block Retain Height value
func (f *Fetcher) GetBlockResultsRetainHeight() (uint64, error) {
	logger := *f.logger.With("method", "GetBlockResultsRetainHeight")
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
	defer cancel()

	// Privileged Services
	conn, err := privileged.New(ctx, f.config.GRPCClient.ListenAddressPrivileged, privileged.WithPruningServiceEnabled(true), privileged.WithInsecure())
	defer conn.Close()
	if err != nil {
		logger.Error("New privileged client", "error", err)
		return 0, fmt.Errorf("error new priviledge client")
	}

	retainHeight, err := conn.GetBlockResultsRetainHeight(ctx)
	if err != nil {
		logger.Error("Get block results retain height", "error", err)
		return 0, fmt.Errorf("error getting block results retain height")
	}
	logger.Info("Get block results retain height", "height", retainHeight)
	return retainHeight, nil
}

// SetBlockResultsRetainHeight Set Block Results Retain Height value
func (f *Fetcher) SetBlockResultsRetainHeight(height uint64) error {
	logger := *f.logger.With("method", "SetBlockResultsRetainHeight")
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestDefaultTimeout)
	defer cancel()

	// Privileged Services
	conn, err := privileged.New(ctx, f.config.GRPCClient.ListenAddressPrivileged, privileged.WithPruningServiceEnabled(true), privileged.WithInsecure())
	defer conn.Close()
	if err != nil {
		logger.Error("New privileged client", "error", err)
		return fmt.Errorf("error new priviledge client: %v", err)
	}

	err = conn.SetBlockResultsRetainHeight(ctx, height)
	if err != nil {
		logger.Error("Set block results retain height", "error", err)
		return fmt.Errorf("error setting block results retain height")
	}
	logger.Info("Set block results retain height", "height", height)
	return nil
}

func (f *Fetcher) GetNewBlockStream(ctx *context.Context) (<-chan client.LatestHeightResult, error) {
	logger := *f.logger.With("method", "GetNewBlockStream")
	// Privileged Services
	conn, err := client.New(*ctx, f.config.GRPCClient.ListenAddress, client.WithBlockServiceEnabled(true), client.WithInsecure())

	newHeightCh, err := conn.GetLatestHeight(*ctx)
	if err != nil {
		logger.Error("Get new block stream", "error", err)
		return nil, fmt.Errorf("error get new block stream")
	}
	logger.Info("Get new block stream")
	return newHeightCh, nil
}

// WatchNewBlock watch for new block events streamed from the cometBFT server
func (f *Fetcher) WatchNewBlock() {
	logger := *f.logger.With("method", "WatchNewBlock")
	ctx := context.Background()
	newHeightCh, err := f.GetNewBlockStream(&ctx)
	if err != nil {
		logger.Error("New block stream", "error", err)
		ctx.Done()
	} else {
		logger.Info("Stream ready")
	}
	go func(c context.Context, ch <-chan client.LatestHeightResult, l slog.Logger) {
		for {
			select {
			case <-c.Done():
				l.Info("Context not available to stream new blocks")
			case latestHeightResult, ok := <-ch:
				if ok {
					if latestHeightResult.Error != nil {
						l.Error("Error in new block", "error", latestHeightResult.Error)
					} else {
						l.Info("New block", "height", latestHeightResult.Height)
					}
				} else {
					l.Info("New block streaming closed")
					//TODO: Instead of returning and exit, keep trying to connect or listening ?
					return
				}
			}
		}
	}(ctx, newHeightCh, logger)
}

//----------------------------------------------------------------------------------------------------------------------
// Services methods

func (f *Fetcher) OnStart() error {

	f.logger.Info("Service running")

	// Get latest block retain height
	rh, err := f.GetBlockRetainHeight()
	if err != nil {
		f.logger.Error("Get block retain height", "error", err)
	}

	// Fetch a block that is one block higher than the lowest block retain height
	height := min(rh.PruningService, rh.App) + 1

	_, err = f.GetBlock(int64(height))
	if err != nil {
		f.logger.Error("Get block", "error", err)
	}

	// Set Block Retain Height if it's higher than zero
	if height > rh.PruningService {
		err = f.SetBlockRetainHeight(height)
		if err != nil {
			f.logger.Error("Set Block Retain Height", "error", err)
		}
	} else {
		f.logger.Info("Skip set block retain height. Retain height higher than App retain height")
	}

	// Check Block Results Retain Height
	h, err := f.GetBlockResultsRetainHeight()
	if err != nil {
		f.logger.Error("Get Block Results Retain Height", "error", err)
	}

	// Get Block Results
	br, err := f.GetBlockResults(int64(h + 1))
	if err != nil {
		f.logger.Error("Get Block Results", "error", err)
	} else {
		// Set Block Results Retain Height
		err = f.SetBlockResultsRetainHeight(uint64(br.Height + 1))
		if err != nil {
			f.logger.Error("Set Block Results Retain Height", "error", err)
		}
	}

	// Stream new block events
	f.WatchNewBlock()

	return nil
}
