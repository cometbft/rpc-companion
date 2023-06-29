package ingest

import (
	"github.com/cometbft/rpc-companion/ingest/client"
	"github.com/cometbft/rpc-companion/ingest/client/http"
	storage "github.com/cometbft/rpc-companion/storage"
	"os"
)

type Service struct {
	Client  client.IClient
	Storage storage.IStorage
}

func NewService(connStr string) Service {

	// Database
	db := storage.PostgresStorage{
		ConnectionString: connStr,
	}

	// HTTP IClient
	client := http.HTTP{
		Endpoint: os.Getenv("COMPANION_NODE_RPC"),
	}

	// Return an Ingest Service
	return Service{
		Client:  &client,
		Storage: &db,
	}
}
