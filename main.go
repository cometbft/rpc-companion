package main

import (
	"context"
	"database/sql"
	"fmt"
	client "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const connString = "postgres://postgres:postgres@0.0.0.0:15432/postgres?sslmode=disable"

type Storage interface {
	Insert(table string, value string) (bool, error)
	Get(table string, query string) ([]byte, error)
	Connect(conn string) error
}

type Fetcher interface {
	Fetch(endpoint string) ([]byte, error)
}

type Service interface {
	Serve()
}

type IngestService struct {
	Fetcher Fetcher
	Storage Storage
}

type CometFetcher struct {
	Endpoint string
}

type PostgresStorage struct {
	ConnectionString string
	Connection       *sql.DB
}

type RESTService struct {
	Version string
}

func (c *CometFetcher) FetchBlock(height int64) (*ctypes.ResultBlock, error) {

	httpClient, err := client.New("http://localhost:26657", "/websocket")
	if err != nil {
		return nil, err
	}

	resultBlock, err := httpClient.Block(context.Background(), &height)
	if err != nil {
		return nil, err
	} else {
		return resultBlock, nil
	}
}

func (c *PostgresStorage) InsertBlock(resultBlock ctypes.ResultBlock) (bool, error) {
	_, err := c.Connection.Exec("INSERT INTO comet.result_block (block_id_hash, block_id_parts_hash, block_id_parts_total, block_header_height, block_header_version_block, block_header_version_app, block_header_block_time, block_header_chain_id, block_last_block_id_hash, block_last_block_id_parts_hash, block_last_block_id_part_total) values ($1,$2,$3,$4, $5, $6, $7, $8, $9, $10, $11)",
		resultBlock.BlockID.Hash.String(),
		resultBlock.BlockID.PartSetHeader.Hash.String(),
		resultBlock.BlockID.PartSetHeader.Total,
		resultBlock.Block.Height,
		resultBlock.Block.Version.Block,
		resultBlock.Block.Version.App,
		resultBlock.Block.Time,
		resultBlock.Block.ChainID,
		resultBlock.Block.LastBlockID.Hash.String(),
		resultBlock.Block.LastBlockID.PartSetHeader.Hash.String(),
		resultBlock.Block.LastBlockID.PartSetHeader.Total)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (c *PostgresStorage) GetBlock(height int64) (int64, error) {
	var rowHeight int64
	row := c.Connection.QueryRow("SELECT block_header_height FROM comet.result_block WHERE block_header_height=$1", height)
	err := row.Scan(&rowHeight)
	if err != nil {
		return 0, err
	} else {
		return rowHeight, err
	}
}

func (c *PostgresStorage) Connect() error {
	db, err := sql.Open("postgres", c.ConnectionString)
	if err != nil {
		return err
	} else {
		c.Connection = db
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	fetcher := CometFetcher{
		Endpoint: "http://localhost:26657",
	}

	storage := PostgresStorage{
		ConnectionString: connString,
		Connection:       nil,
	}

	// Connect to the database
	err := storage.Connect()
	if err != nil {
		panic(err)
	}

	for height := 1; height <= 50; height++ {

		blockFetched, err := fetcher.FetchBlock(int64(height))
		if err != nil {
			log.Fatalf("Error fetching block at height %d: %s\n", height, err)
		}

		inserted, err := storage.InsertBlock(*blockFetched)
		if err != nil {
			fmt.Printf("Error inserting block at height %d: %s\n", height, err)
		}
		if inserted {
			fmt.Printf("Inserted height %d\n", height)
		}

		block, err := storage.GetBlock(int64(height))
		if err != nil {
			fmt.Printf("Error retrieving block at height %d: %s\n", height, err)
		} else {
			log.Printf("Block at height %d: %v", height, block)
		}
	}

	defer func(ps PostgresStorage) {
		err := ps.Connection.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(storage)
}
