package rpc

import (
	"encoding/json"
	"github.com/cometbft/cometbft/rpc/jsonrpc/types"
	"github.com/cometbft/rpc-companion/storage"
	"log"
	"net/http"
	"strconv"
)

type Service struct {
	Storage storage.IStorage
}

func NewService(connStr string) Service {

	// Database
	db := storage.PostgresStorage{
		ConnectionString: connStr,
	}

	// Return an Ingest Service
	return Service{
		Storage: &db,
	}
}

func (s *Service) Serve() {

	// Handler for the block endpoint
	http.HandleFunc("/v1/header", s.HeaderHandler)

	// Start the service
	log.Fatalln(http.ListenAndServe(":8080", nil)) // TODO: Make the port configurable
}

func (s *Service) HeaderHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	conn, err := s.Storage.Connect()
	defer conn.Close()
	if err != nil {
		log.Println("error connecting to storage in HeaderHandler: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Internal Server Error"))
	} else {
		if request.Method == "GET" {
			h := request.URL.Query()["height"][0]
			height, err := strconv.ParseInt(h, 10, 64)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("Bad Request. Invalid height"))
			}
			log.Printf("header request at height: %v\n", height)
			header, err := s.Storage.GetHeader(height)
			if err != nil {
				// TODO: If not records retrieved return a different status
				log.Println("error retrieving header from storage in HeaderHandler: ", err)
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("Internal Server Error"))
			}
			// Return response
			id := types.JSONRPCStringID("id")
			resp := types.NewRPCSuccessResponse(id, header)
			respJson, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling header response: ", err)
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("Internal Server Error"))
			} else {
				writer.WriteHeader(http.StatusOK)
				writer.Write(respJson)
			}

		} else {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Bad Request"))
		}
	}
}
