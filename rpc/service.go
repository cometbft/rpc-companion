package rpc

import (
	"github.com/cometbft/rpc-companion/storage"
	"log"
	"net/http"
)

type Service struct {
	Storage storage.IStorage
}

func NewService(connStr string) Service {

	// Database
	db := storage.PostgresStorage{
		ConnectionString: connStr,
	}

	// Return an Ingest Service
	return Service{
		Storage: &db,
	}
}

func (s *Service) Serve() {

	// Start the service
	log.Fatalln(http.ListenAndServe(":8080", nil)) // TODO: Make the port configurable
}
