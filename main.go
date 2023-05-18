package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cometbft/cometbft/libs/json"
	client "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
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

func (s *RESTService) Serve(storage *PostgresStorage) {
	// Handler for the block endpoint
	http.HandleFunc(fmt.Sprintf("/%s/block", s.Version), handleBlock)

	// Start the service
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("There's an error starting the REST service:", err)
	} else {
		log.Println("Started REST service...")
	}
}

// Handles the '/v1/block' endpoint
func handleBlock(writer http.ResponseWriter, request *http.Request) {

	// Database connection
	storage := PostgresStorage{
		ConnectionString: connString,
		Connection:       nil,
	}

	// Connect to the database
	err := storage.Connect()
	if err != nil {
		log.Println("Error connecting to storage in handleBlock: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Internal Server Error"))
	}

	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "GET" {
		h := request.URL.Query()["height"][0]
		height, err := strconv.ParseInt(h, 10, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Bad Request. Invalid height"))
		}
		fmt.Printf("Block Request. Height: %v\n", height)
		block, err := storage.GetBlock(height)
		if err != nil {
			log.Println("Error retrieving record from storage in handleBlock: ", err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Internal Server Error"))
		}
		resp, _ := json.Marshal(block)
		writer.Write(resp)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad Request"))
	}
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

func (c *PostgresStorage) GetBlock(height int64) (ctypes.ResultBlock, error) {
	resultBlock := ctypes.ResultBlock{}
	b := new(types.Block)
	row := c.Connection.QueryRow("SELECT block_header_height, block_header_chain_id, block_header_block_time FROM comet.result_block WHERE block_header_height=$1", height)
	err := row.Scan(&b.Header.Height, &b.Header.ChainID, &b.Header.Time)
	if err != nil {
		return resultBlock, err
	}
	resultBlock.Block = b
	return resultBlock, err
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

	// Ingest server
	fetcher := CometFetcher{
		Endpoint: "http://localhost:26657",
	}

	// Database storage
	storage := PostgresStorage{
		ConnectionString: connString,
		Connection:       nil,
	}

	// Connect to the database
	err := storage.Connect()
	if err != nil {
		panic(err)
	}

	for height := 1; height <= 100; height++ {

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
	}

	defer func(ps PostgresStorage) {
		err := ps.Connection.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(storage)

	// REST server
	service := RESTService{
		Version: "v1",
	}
	service.Serve(&storage)
}
