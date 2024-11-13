package ingest

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cometbft/cometbft/rpc/grpc/client"
	"github.com/cometbft/cometbft/rpc/grpc/client/privileged"
	"github.com/cometbft/rpc-companion/config"
	"github.com/cometbft/rpc-companion/storage"
)

type EventFetcher struct {
	BaseService
	config   *config.Config
	services *ServiceClient
	context  context.Context
	logger   slog.Logger
	storage  *storage.Storage
}

var blockQueueResults = make(chan Job[client.BlockResults])

func NewEventFetcher(logger slog.Logger, cfg *config.Config) (*EventFetcher, error) {
	logger = *logger.With("module", "EventFetcher")

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

	return &EventFetcher{
		logger:   logger,
		config:   cfg,
		context:  ctx,
		services: &services,
		storage:  &db,
	}, nil
}

//----------------------------------------------------------------------------------------------------------------------
// Requests

// GetBlockResults returns block results at a specific height
func (f *EventFetcher) GetBlockResults(height int64) (*client.BlockResults, error) {
	logger := *f.logger.With("method", "GetBlockResults")

	blockResults, err := f.services.client.GetBlockResults(f.context, height)
	if err != nil {
		logger.Error("Get block results", "error", err)
		return nil, fmt.Errorf("error getting block results")
	}
	logger.Info("Get block results", "height", height)
	return blockResults, nil
}

// GetBlockResultsRetainHeight Get Block Retain Height value
func (f *EventFetcher) GetBlockResultsRetainHeight() (uint64, error) {
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
func (f *EventFetcher) SetBlockResultsRetainHeight(height uint64) error {
	logger := *f.logger.With("method", "SetBlockResultsRetainHeight")

	err := f.services.privilegedClient.SetBlockResultsRetainHeight(f.context, height)
	if err != nil {
		logger.Error("Set block results retain height", "error", err)
		return fmt.Errorf("error setting block results retain height")
	}
	logger.Info("Set block results retain height", "height", height)
	return nil
}

func (f *EventFetcher) GetNewBlockStream() (<-chan client.LatestHeightResult, error) {
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
func (f *EventFetcher) WatchNewBlock() {
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

	go func(f *EventFetcher, c context.Context, ch <-chan client.LatestHeightResult, l slog.Logger) {
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
						blockResult, err := f.GetBlockResults(latestHeightResult.Height)
						if err != nil {
							l.Error("Get block from storage", "error", err)
						} else {
							job := NewJob(*blockResult)
							blockQueueResults <- job
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

func (f *EventFetcher) ProcessBlockJob() {
	logger := *f.logger.With("method", "ProcessBlockJob")
	logger.Info("Starting Worker")
	go func(fetcher *EventFetcher) {
		for {
			job := <-blockQueueResults
			fetcher.logger.Info("Processing job", "height", job.cometType.Height)

			// fmt.Println(&job.cometType.FinalizeBlockEvents)
			// fmt.Println(&job.cometType.TxResults)

			ForwardData(job)
			// TODO Understand if data is stored and instruct CometBFT to delete it
		}
	}(f)
}

//----------------------------------------------------------------------------------------------------------------------
// ServiceClient methods

func (f *EventFetcher) OnStart() error {
	f.logger.Info("Service running")
	// Stream new block events
	f.WatchNewBlock()

	return nil
}

func (f *EventFetcher) OnStop() {
	f.logger.Info("Service stopping")
}
