package main

/*
import (
	"fmt"
	"github.com/cometbft/cometbft/libs/json"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

type Service interface {
	Serve()
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
			// TODO: If not records retrieved return a different status
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

func insert() {

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

	// Insert some blocks
	InsertBlocks(storage)

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
*/
