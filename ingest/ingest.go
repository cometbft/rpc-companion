package ingest

import (
	"fmt"
	"log/slog"

	"github.com/cometbft/rpc-companion/config"
)

// IngestService orchestrates the ingest services
type IngestService struct {
	BaseService
	config  *config.Config
	fetcher *Fetcher
	//storage storage.IStorage
}

func NewIngestService(
	logger slog.Logger,
	config config.Config,
) (*IngestService, error) {
	logger = *logger.With("service", "Ingest")

	// Instantiate new fetcher (gRPC client)
	fetcher, err := NewFetcher(logger, &config)
	if err != nil {
		logger.Error("Creating new fetcher", "error", err)
		return nil, fmt.Errorf("error creating new fetcher")
	}

	// Configure Fetcher service
	fetcher.BaseService = *NewBaseService(logger, "Fetcher", fetcher)

	// Ingest Service
	ingest := &IngestService{
		config:  &config,
		fetcher: fetcher,
		//storage: &db,
	}

	ingest.BaseService = *NewBaseService(logger, "Ingest", ingest)

	return ingest, nil
}

func (s *IngestService) OnStart() error {
	if s.IsRunning() {
		s.fetcher.Start()
	}
	return nil
}

func (s *IngestService) OnStop() {
	if s.fetcher.IsRunning() {
		s.fetcher.Stop()
	}
	s.BaseService.OnStop()
}
