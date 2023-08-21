package storage

import (
	"context"
	"database/sql"
	"github.com/cometbft/cometbft/libs/json"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	_ "github.com/lib/pq"
	"time"
)

type PostgresStorage struct {
	ConnectionString string
}

func (c *PostgresStorage) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.ConnectionString)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func (c *PostgresStorage) Disconnect(conn *sql.DB) error {
	err := conn.Close()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *PostgresStorage) Ping() error {
	conn, err := c.Connect()
	if err != nil {
		return err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = conn.PingContext(ctx)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (c *PostgresStorage) InsertHeader(height int64, header ctypes.ResultHeader) error {
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return err
	} else {
		headerBytes, err := json.Marshal(header)
		if err != nil {
			return err
		} else {
			_, err = conn.Exec("INSERT INTO comet.header (height, header) values ($1,$2)", height, &headerBytes)
			if err != nil {
				return err
			} else {
				return nil
			}
		}
	}
}

func (c *PostgresStorage) GetHeader(height int64) (ctypes.ResultHeader, error) {
	var header ctypes.ResultHeader
	var headerBytes []byte
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return header, err
	} else {
		row := conn.QueryRow("SELECT header FROM comet.header WHERE height=$1", height)
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
}
