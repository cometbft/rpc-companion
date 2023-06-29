package storage

import (
	"database/sql"
	"github.com/cometbft/cometbft/proto/tendermint/types"
)

type IStorage interface {
	Connect() (*sql.DB, error)
	Disconnect(db *sql.DB) error
	Ping() error
	InsertHeader(height int64, header types.Header) error
	//GetHeader(height int64) (types.Header, error)
}
