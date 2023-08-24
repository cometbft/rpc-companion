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

	// Instantiate new fetcher (gRPC client)
	fetcher, err := NewFetcher(logger, &config)
	if err != nil {
		return nil, fmt.Errorf("error creating new fetcher: %v", err)
	}

	// Configure Fetcher service
	fetcher.BaseService = *NewBaseService(logger, "Fetcher", fetcher)
	fetcher.SetLogger(*logger.With("service", "fetcher"))

	// Ingest Service
	ingest := &IngestService{
		config:  &config,
		fetcher: fetcher,
		//storage: &db,
	}

	ingest.BaseService = *NewBaseService(logger, "Ingest", ingest)
	ingest.SetLogger(*logger.With("service", "ingest"))

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
