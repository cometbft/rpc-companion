package ingest

import (
	"github.com/cometbft/rpc-companion/config"
	"github.com/cometbft/rpc-companion/service/base"
	"log/slog"
)

type Service struct {
	service.BaseService
	config *config.Config
	//client  client.IClient
	//storage storage.IStorage
}

func NewIngestService(
	logger slog.Logger,
) *Service {

	cfg := config.DefaultConfig()

	//// Database
	//db := storage.PostgresStorage{
	//	ConnectionString: "",
	//}
	//
	//// HTTP IClient
	//client := http.HTTP{
	//	Endpoint: os.Getenv("COMPANION_NODE_RPC"),
	//}

	// Ingest Service
	svc := &Service{
		config: cfg,
		//client:  &client,
		//storage: &db,
	}

	svc.BaseService = *service.NewBaseService(logger, "Ingest", svc)
	iLogger := svc.Logger.With("service", "ingest")
	svc.SetLogger(*iLogger)
	return svc
}

func (s *Service) OnStart() error {
	s.Logger.Info("Ingest Service starting")
	return nil
}

func (s *Service) OnStop() {
	s.BaseService.OnStop()
	s.Logger.Info("Stopping Node")
	//TODO: Add stopping logic
}
