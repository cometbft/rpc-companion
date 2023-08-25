package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cometbft/cometbft/rpc/grpc/client"
	"github.com/cometbft/cometbft/rpc/grpc/client/privileged"
	"github.com/cometbft/rpc-companion/config"
	"github.com/cometbft/rpc-companion/storage"
)

var (
	requestDefaultTimeout = 10 * time.Second
	blockQueue            = make(chan Job[client.Block]) // Queue to process blocks
)

type CometType interface {
	client.Block | client.BlockResults
}

type Fetcher struct {
	BaseService
	config   *config.Config
	services *ServiceClient
	context  context.Context
	logger   slog.Logger
	storage  *storage.Storage
}

type Job[T CometType] struct {
	done      bool
	cometType T
}

func NewJob[T CometType](cType T) Job[T] {
	return Job[T]{
		done:      false,
		cometType: cType,
	}
}

func NewFetcher(logger slog.Logger, cfg *config.Config) (*Fetcher, error) {
	logger = *logger.With("module", "Fetcher")

	ctx := context.Background()

	// Service
	conn, err := client.New(ctx, cfg.GRPCClient.ListenAddress, client.WithBlockServiceEnabled(true), client.WithInsecure()) //TODO: In the future support secure connections
	if err != nil {
		logger.Error("New client", "error", err)
		return nil, fmt.Errorf("error creating new client")
	}

	// Privileged ServiceClient
	privConn, err := privileged.New(ctx, cfg.GRPCClient.ListenAddressPrivileged, privileged.WithPruningServiceEnabled(true), privileged.WithInsecure())
	if err != nil {
		logger.Error("New privileged client", "error", err)
		return nil, fmt.Errorf("error creating new privileged client")
	}

	// Client services connections
	services := ServiceClient{
		client:           conn,
		privilegedClient: privConn,
	}

	// Storage
	db, err := storage.NewStorage(cfg.Storage.Connection)
	if err != nil {
		logger.Error("New storage", "error", err)
		return nil, fmt.Errorf("error creating new storage")
	}

	return &Fetcher{
		logger:   logger,
		config:   cfg,
		context:  ctx,
		services: &services,
		storage:  &db,
	}, nil
}

//----------------------------------------------------------------------------------------------------------------------
// Requests

// GetBlock returns block at a specific height
func (f *Fetcher) GetBlock(height int64) (*client.Block, error) {
	logger := *f.logger.With("method", "GetBlock")

	block, err := f.services.client.GetBlockByHeight(f.context, height)
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

	blockResults, err := f.services.client.GetBlockResults(f.context, height)
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

	retainHeight, err := f.services.privilegedClient.GetBlockRetainHeight(f.context)
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

	err := f.services.privilegedClient.SetBlockRetainHeight(f.context, height)
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

	retainHeight, err := f.services.privilegedClient.GetBlockResultsRetainHeight(f.context)
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

	err := f.services.privilegedClient.SetBlockResultsRetainHeight(f.context, height)
	if err != nil {
		logger.Error("Set block results retain height", "error", err)
		return fmt.Errorf("error setting block results retain height")
	}
	logger.Info("Set block results retain height", "height", height)
	return nil
}

func (f *Fetcher) GetNewBlockStream() (<-chan client.LatestHeightResult, error) {
	logger := *f.logger.With("method", "GetNewBlockStream")

	newHeightCh, err := f.services.client.GetLatestHeight(f.context)
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
	newHeightCh, err := f.GetNewBlockStream()
	if err != nil {
		logger.Error("New block stream", "error", err)
		ctx.Done()
	} else {
		logger.Info("Stream ready")
	}

	// Start the queue processor
	f.ProcessBlockJob()

	go func(f *Fetcher, c context.Context, ch <-chan client.LatestHeightResult, l slog.Logger) {
		for {
			select {
			case <-c.Done():
				l.Info("Connection not available to stream new blocks")
			case latestHeightResult, ok := <-ch:
				if ok {
					if latestHeightResult.Error != nil {
						l.Error("Error in new block", "error", latestHeightResult.Error)
					} else {
						l.Info("New block", "height", latestHeightResult.Height)
						block, err := f.services.client.GetBlockByHeight(c, latestHeightResult.Height)
						if err != nil {
							l.Error("Get block from storage", "error", err)
						} else {
							job := NewJob(*block)
							blockQueue <- job
						}
					}
				} else {
					l.Info("New block streaming closed")
					//TODO: Instead of returning and exit, keep trying to connect or listening ?
					return
				}
			}
		}
	}(f, ctx, newHeightCh, logger)
}

func (f *Fetcher) ProcessBlockJob() {
	logger := *f.logger.With("method", "ProcessBlockJob")
	logger.Info("Starting Worker")
	go func(fetcher *Fetcher) {
		for {
			job := <-blockQueue
			fetcher.logger.Info("Processing job", "height", job.cometType.Block.Height)
			err := fetcher.storage.InsertBlock(uint64(job.cometType.Block.Height), &job.cometType)
			if err != nil {
				logger.Error("Process block job", "error", err)
			} else {
				job.done = true
				logger.Info("Processed block job", "height", job.cometType.Block.Height)
			}
		}
	}(f)
}

//----------------------------------------------------------------------------------------------------------------------
// ServiceClient methods

func (f *Fetcher) OnStart() error {

	f.logger.Info("Service running")

	//// Get latest block retain height
	//rh, err := f.GetBlockRetainHeight()
	//if err != nil {
	//	f.logger.Error("Get block retain height", "error", err)
	//}
	//
	//// Fetch a block that is one block higher than the lowest block retain height
	//height := min(rh.PruningService, rh.App) + 1
	//
	//_, err = f.GetBlock(int64(height))
	//if err != nil {
	//	f.logger.Error("Get block", "error", err)
	//}
	//
	//// Set Block Retain Height if it's higher than zero
	//if height > rh.PruningService {
	//	err = f.SetBlockRetainHeight(height)
	//	if err != nil {
	//		f.logger.Error("Set Block Retain Height", "error", err)
	//	}
	//} else {
	//	f.logger.Info("Skip set block retain height. Retain height higher than App retain height")
	//}
	//
	//// Check Block Results Retain Height
	//h, err := f.GetBlockResultsRetainHeight()
	//if err != nil {
	//	f.logger.Error("Get Block Results Retain Height", "error", err)
	//}
	//
	//// Get Block Results
	//br, err := f.GetBlockResults(int64(h + 1))
	//if err != nil {
	//	f.logger.Error("Get Block Results", "error", err)
	//} else {
	//	// Set Block Results Retain Height
	//	err = f.SetBlockResultsRetainHeight(uint64(br.Height + 1))
	//	if err != nil {
	//		f.logger.Error("Set Block Results Retain Height", "error", err)
	//	}
	//}

	// Stream new block events
	f.WatchNewBlock()

	return nil
}

func (f *Fetcher) OnStop() {
	f.logger.Info("Service stopping")
}
