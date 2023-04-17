package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
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
	url := fmt.Sprintf("http://localhost:26657/%s", endpoint)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}

func (c *PostgresStorage) Insert(table string, value []byte) (bool, error) {

	data := json.RawMessage(value)
	insertStmt := fmt.Sprintf("INSERT INTO comet.%s(blob) VALUES ($1)", table)
	_, err := c.Connection.Exec(insertStmt, data)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

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

	for height := 1; height <= 50; height++ {

		resp, err := fetcher.Fetch(fmt.Sprintf("block?height=%d", height))
		if err != nil {
			log.Fatalf("Error fetching height %d: %s\n", height, err)
		}

		// Connect to the database
		err = storage.Connect()
		if err != nil {
			panic(err)
		}

		inserted, err := storage.Insert("blocks", resp)
		if err != nil {
			fmt.Printf("Error inserting height %d: %s\n", height, err)
		}
		if inserted {
			fmt.Printf("Inserted height %d successfully\n", height)
		}
	}

	defer func(ps PostgresStorage) {
		err := ps.Connection.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(storage)
}
