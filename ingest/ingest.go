package ingest

import (
	"fmt"
	"log/slog"

	"github.com/cometbft/cometbft/rpc/grpc/client"
	"github.com/cometbft/cometbft/rpc/grpc/client/privileged"
	"github.com/cometbft/rpc-companion/config"
)

// IngestService orchestrates the ingest services
type IngestService struct {
	BaseService
	config       *config.Config
	fetcher      *Fetcher
	eventFetcher *EventFetcher
}

// ServiceClient GRPC clients
type ServiceClient struct {
	client           client.Client
	privilegedClient privileged.Client
}

func NewIngestService(
	logger slog.Logger,
	config config.Config,
) (*IngestService, error) {
	logger = *logger.With("service", "Ingest")

	// TODO Commented out to disable block fetching. Enable using both

	// Instantiate new fetcher (gRPC client)
	// fetcher, err := NewFetcher(logger, &config)
	// if err != nil {
	// 	logger.Error("Creating new fetcher", "error", err)
	// 	return nil, fmt.Errorf("error creating new fetcher")
	// }

	// // Configure Fetcher service
	// fetcher.BaseService = *NewBaseService(logger, "Fetcher", fetcher)

	eventFetcher, err := NewEventFetcher(logger, &config)

	if err != nil {
		logger.Error("creating new event fetcher", "error", err)
		return nil, fmt.Errorf("error creating new event fetcher")
	}

	eventFetcher.BaseService = *NewBaseService(logger, "EventFetcher", eventFetcher)
	// Ingest Service
	ingest := &IngestService{
		config: &config,
		// fetcher:      fetcher,
		eventFetcher: eventFetcher,
		//storage: &db,
	}

	ingest.BaseService = *NewBaseService(logger, "Ingest", ingest)

	return ingest, nil
}

func (s *IngestService) OnStart() error {
	if s.IsRunning() {
		// s.fetcher.Start()
		s.eventFetcher.Start()
	}
	return nil
}

func (s *IngestService) OnStop() {
	// if s.fetcher.IsRunning() {
	// 	s.fetcher.Stop()
	// } // Commented out to not start fetching blocks
	if s.eventFetcher.IsRunning() {
		s.eventFetcher.Stop()
	}
	s.BaseService.OnStop()
}
