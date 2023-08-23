package ingest

import (
	"fmt"
	"github.com/cometbft/rpc-companion/config"
	"github.com/cometbft/rpc-companion/service/base"
	"github.com/cometbft/rpc-companion/service/ingest/fetcher"
	"log/slog"
)

type Service struct {
	service.BaseService
	config  *config.Config
	fetcher *fetcher.Fetcher
	//storage storage.IStorage
}

func NewIngestService(
	logger slog.Logger,
) (*Service, error) {

	cfg := config.DefaultConfig()

	//// Database
	//db := storage.PostgresStorage{
	//	ConnectionString: "",
	//}
	//

	// Instantiate new fetcher (gRPC client)
	fetcher, err := fetcher.NewFetcher(cfg.GRPCClient.ListenAddress, logger, cfg.GRPCClient)
	if err != nil {
		return nil, fmt.Errorf("error creating new fetcher: %v", err)
	}

	// Configure Fetcher service
	fetcher.BaseService = *service.NewBaseService(logger, "Fetcher", fetcher)
	fetcher.SetLogger(*logger.With("service", "fetcher"))

	// Ingest Service
	ingest := &Service{
		config:  cfg,
		fetcher: fetcher,
		//storage: &db,
	}

	ingest.BaseService = *service.NewBaseService(logger, "Ingest", ingest)
	ingest.SetLogger(*logger.With("service", "ingest"))

	return ingest, nil
}

func (s *Service) OnStart() error {
	if s.IsRunning() {
		if err := s.config.GRPCClient.ValidateBasic(); err != nil {
			return fmt.Errorf("error validating gRPC client configuration: %v", err)
		}
		s.fetcher.Start()
	}
	return nil
}

func (s *Service) OnStop() {
	if s.fetcher.IsRunning() {
		s.fetcher.Stop()
	}
	s.BaseService.OnStop()
	//TODO: Add stopping logic
}
