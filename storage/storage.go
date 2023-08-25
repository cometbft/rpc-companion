package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/cometbft/cometbft/libs/json"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	_ "github.com/lib/pq"
)

var (
	driverName = "postgres"
)

type Storage struct {
	connection *sql.DB
}

func NewStorage(connectionString string) (Storage, error) {
	db := Storage{}
	conn, err := db.Connect(connectionString)
	if err != nil {
		return db, err
	}
	db.connection = conn
	return db, nil
}

func (c *Storage) Connect(conn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, conn)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func (c *Storage) Disconnect() error {
	err := c.connection.Close()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.connection.PingContext(ctx)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Storage) InsertHeader(height int64, header ctypes.ResultHeader) error {
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return err
	} else {
		_, err = c.connection.Exec("INSERT INTO comet.header (height, header) values ($1,$2)", height, &headerBytes)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (c *Storage) GetHeader(height int64) (ctypes.ResultHeader, error) {
	var header ctypes.ResultHeader
	var headerBytes []byte
	row := c.connection.QueryRow("SELECT header FROM comet.header WHERE height=$1", height)
	err := row.Scan(
		&headerBytes)
	if err != nil {
		return header, err
	}
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return header, err
	}
	return header, nil
}
