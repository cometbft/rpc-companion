package ingest

import (
	storage "github.com/cometbft/rpc-companion/storage"
	"os"
)

type Service struct {
	Fetcher Fetcher
	Storage storage.IStorage
}

func NewService(connStr string) Service {

	// Database
	db := storage.PostgresStorage{
		ConnectionString: connStr,
	}

	// RPC Fetcher
	fetcher := RPCFetcher{
		Endpoint: os.Getenv("COMPANION_NODE_RPC"),
	}

	// Return an Ingest Service
	return Service{
		Fetcher: &fetcher,
		Storage: &db,
	}
}
