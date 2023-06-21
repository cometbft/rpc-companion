package rpc

import (
	"fmt"
	cmtjson "github.com/cometbft/cometbft/libs/json"
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

func (s *Service) Serve() error {

	// Handler for the block endpoint
	http.HandleFunc("/v1/block", s.handleBlock)

	// Start the service
	err := http.ListenAndServe(":8080", nil) // TODO: Make the port configurable
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Handles the '/v1/block' endpoint
func (s *Service) handleBlock(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	conn, err := s.Storage.Connect()
	defer conn.Close()
	if err != nil {
		log.Println("Error connecting to storage in handleBlock: ", err)
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
			fmt.Printf("Block Request. Height: %v\n", height)
			block, err := s.Storage.GetBlock(height)
			if err != nil {
				// TODO: If not records retrieved return a different status
				log.Println("Error retrieving record from storage in handleBlock: ", err)
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("Internal Server Error"))
			}
			// Return response
			//TODO: Empty objects return 'null' should return '[]'
			blockJSON, err := cmtjson.Marshal(block)
			if err != nil {
				log.Println("Error marshalling block: ", err)
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("Internal Server Error"))
			} else {
				var RPCResponse = types.RPCResponse{
					JSONRPC: "2.0",
					ID:      nil, //TODO: Figure out a way to properly return this avoiding error 'cannot encode unregistered type types.JSONRPCIntID'
					Result:  blockJSON,
					Error:   nil,
				}
				resp, err := cmtjson.Marshal(RPCResponse)
				if err != nil {
					log.Println("Error marshalling RPCResponse: ", err)
					writer.WriteHeader(http.StatusInternalServerError)
					writer.Write([]byte("Internal Server Error"))
				} else {
					writer.WriteHeader(http.StatusOK)
					writer.Write(resp)
				}
			}
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Bad Request"))
		}
	}
}
