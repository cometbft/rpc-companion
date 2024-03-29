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
						block, err := f.GetBlock(latestHeightResult.Height)
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
				// Update block retain height if lower
				// Get latest block retain height
				rh, err := f.GetBlockRetainHeight()
				if err != nil {
					logger.Error("Get block retain height", "error", err)
				} else {
					if rh.PruningService < uint64(job.cometType.Block.Height) {
						// This is a naive way of setting the retain height,
						// ideally there should be a process that checks the storage
						// to query inserted blocks and if there's a gap in the last
						// inserted block and the block in the job. Setting to the job
						// height will prune previous blocks that were not inserted yet.
						err := f.SetBlockRetainHeight(uint64(job.cometType.Block.Height))
						if err != nil {
							logger.Error("Set block retain height", "error", err)
						}
					}
				}
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
	// Stream new block events
	f.WatchNewBlock()

	return nil
}

func (f *Fetcher) OnStop() {
	f.logger.Info("Service stopping")
}
