package storage

import (
	"database/sql"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type IStorage interface {
	Connect() (*sql.DB, error)
	Disconnect(db *sql.DB) error
	Ping() error
	InsertHeader(height int64, header ctypes.ResultHeader) error
	GetHeader(height int64) (ctypes.ResultHeader, error)
}
