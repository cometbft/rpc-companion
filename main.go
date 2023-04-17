package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Storage interface {
	Insert(table string, value string) (bool, error)
	Connect(conn string) error
}

type Fetcher interface {
	Fetch(endpoint string) ([]byte, error)
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

func (c *CometFetcher) Fetch(endpoint string) ([]byte, error) {
	fmt.Println("Fetching...")
	return nil, nil
}

func (c *PostgresStorage) Insert(table string, value string) (bool, error) {
	insertStmt := fmt.Sprintf("INSERT INTO comet.%s(blob) VALUES (%s)", table, value)
	_, err := c.Connection.Exec(insertStmt)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// Maybe use generics to return a database connection
func (c *PostgresStorage) Connect() error {
	db, err := sql.Open("postgres", c.ConnectionString)
	if err != nil {
		return err
	} else {
		c.Connection = db
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

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
		ConnectionString: "postgres://postgres:postgres@0.0.0.0:15432/postgres?sslmode=disable",
		Connection:       nil,
	}

	_, err := fetcher.Fetch("/block")

	err = storage.Connect()
	if err != nil {
		panic(err)
	}
	defer func(ps PostgresStorage) {
		err := ps.Connection.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(storage)

	inserted, err := storage.Insert("blocks", "'{\"jsonrpc\":\"2.0\",\"id\":-1,\"result\":{\"height\":\"3\",\"txs_results\":null,\"begin_block_events\":null,\"end_block_events\":null,\"validator_updates\":null,\"consensus_param_updates\":null}}'")
	if err != nil {
		fmt.Printf("Error inserting: %s\n", err)
	}
	if inserted {
		fmt.Println("Inserted successfully")
	}
}
