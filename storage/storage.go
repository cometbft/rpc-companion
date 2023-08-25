package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/rpc/grpc/client"
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

func (c *Storage) InsertBlock(height uint64, block *client.Block) error {
	data, err := json.Marshal(block)
	if err != nil {
		return err
	} else {
		_, err = c.connection.Exec("INSERT INTO comet.block (height, data) values ($1,$2)", height, &data)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (c *Storage) GetHeader(height uint64) (*client.Block, error) {
	var block *client.Block
	var data []byte
	row := c.connection.QueryRow("SELECT block FROM comet.block WHERE height=$1", height)
	err := row.Scan(&data)
	if err != nil {
		return block, err
	}
	err = json.Unmarshal(data, &block)
	if err != nil {
		return block, err
	}
	return block, nil
}
