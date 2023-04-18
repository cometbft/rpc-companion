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

const connstr = "postgres://postgres:postgres@0.0.0.0:15432/postgres?sslmode=disable"

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

func (c *CometFetcher) Fetch(endpoint string) ([]byte, error) {

	url := fmt.Sprintf("http://localhost:26657/%s", endpoint)
	method := "GET"
	fmt.Printf("Fetching %s\n", url)

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

func (c *PostgresStorage) Get(table string, parameter string) ([]byte, error) {
	var response []byte
	queryStmt := fmt.Sprintf("SELECT * FROM comet.%s WHERE blob @> '{\"result\":{\"block\":{\"header\":{\"height\": \"%s\"}}}}';", table, parameter)
	err := c.Connection.QueryRow(queryStmt).Scan(&response)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

//SELECT * FROM comet.blocks WHERE blob @> '{"result":{"block":{"header":{"height": "45"}}}}';

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
		ConnectionString: connstr,
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
		height := request.URL.Query()["height"][0]
		block, err := storage.Get("blocks", height)
		if err != nil {
			log.Println("Error retrieving record from storage in handleBlock: ", err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Internal Server Error"))
		}
		writer.Write(block)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad Request"))
	}
}

func main() {

	fetcher := CometFetcher{
		Endpoint: "http://localhost:26657",
	}

	storage := PostgresStorage{
		ConnectionString: connstr,
		Connection:       nil,
	}

	service := RESTService{
		Version: "v1",
	}

	// Connect to the database
	err := storage.Connect()
	if err != nil {
		panic(err)
	}

	for height := 1; height <= 100; height++ {

		resp, err := fetcher.Fetch(fmt.Sprintf("block?height=%d", height))
		if err != nil {
			log.Fatalf("Error fetching height %d: %s\n", height, err)
		}

		inserted, err := storage.Insert("blocks", resp)
		if err != nil {
			fmt.Printf("Error inserting height %d: %s\n", height, err)
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

	service.Serve(&storage)
}
