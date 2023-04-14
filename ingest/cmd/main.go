package main

import "github.com/cometbft/rpc-companion/ingest"

func main() {

	uri := "mongodb://postgres:postgres@0.0.0.0/ferretdb?authMechanism=PLAIN"

	ingest.InsertRecord(uri)
}
