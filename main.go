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
	"os"
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
	FetchBlock(height int64) (*ctypes.ResultBlock, error)
	FetchABCIInfo() (*ctypes.ResultABCIInfo, error)
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

func (c *CometFetcher) FetchBlock(height int64) (*ctypes.ResultBlock, error) {

	httpClient, err := client.New(c.Endpoint, "/websocket")
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

func (c *CometFetcher) FetchABCIInfo() (*ctypes.ResultABCIInfo, error) {

	httpClient, err := client.New(c.Endpoint, "/websocket")
	if err != nil {
		return nil, err
	}

	abciInfo, err := httpClient.ABCIInfo(context.Background())
	if err != nil {
		return nil, err
	} else {
		return abciInfo, nil
	}
}

func (c *PostgresStorage) InsertBlock(resultBlock ctypes.ResultBlock) (bool, error) {
	_, err := c.Connection.Exec(
		"INSERT INTO comet.result_block"+
			"(block_id_hash, "+
			"block_id_parts_hash, "+
			"block_id_parts_total, "+
			"block_header_height, "+
			"block_header_version_block, "+
			"block_header_version_app, "+
			"block_header_block_time, "+
			"block_header_chain_id, "+
			"block_header_last_commit_hash, "+
			"block_header_data_hash, "+
			"block_header_validators_hash, "+
			"block_header_next_validators_hash, "+
			"block_header_consensus_hash, "+
			"block_header_app_hash, "+
			"block_header_last_results_hash, "+
			"block_header_evidence_hash, "+
			"block_header_proposer_address, "+
			"block_header_last_block_id_hash, "+
			"block_header_last_block_id_parts_hash, "+
			"block_header_last_block_id_part_total, "+
			"block_last_commit_height, "+
			"block_last_commit_round, "+
			"block_last_commit_block_id_hash, "+
			"block_last_commit_block_id_parts_total, "+
			"block_last_commit_block_id_parts_hash) "+
			"VALUES ($1,$2,$3,$4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)",
		resultBlock.BlockID.Hash,
		resultBlock.BlockID.PartSetHeader.Hash,
		resultBlock.BlockID.PartSetHeader.Total,
		resultBlock.Block.Header.Height,
		resultBlock.Block.Header.Version.Block,
		resultBlock.Block.Header.Version.App,
		resultBlock.Block.Header.Time,
		resultBlock.Block.Header.ChainID,
		resultBlock.Block.Header.LastCommitHash,
		resultBlock.Block.Header.DataHash,
		resultBlock.Block.Header.ValidatorsHash,
		resultBlock.Block.Header.NextValidatorsHash,
		resultBlock.Block.Header.ConsensusHash,
		resultBlock.Block.Header.AppHash,
		resultBlock.Block.Header.LastResultsHash,
		resultBlock.Block.Header.EvidenceHash,
		resultBlock.Block.Header.ProposerAddress,
		resultBlock.Block.LastBlockID.Hash,
		resultBlock.Block.LastBlockID.PartSetHeader.Hash,
		resultBlock.Block.LastBlockID.PartSetHeader.Total,
		resultBlock.Block.LastCommit.Height,
		resultBlock.Block.LastCommit.Round,
		resultBlock.Block.LastCommit.BlockID.Hash,
		resultBlock.Block.LastCommit.BlockID.PartSetHeader.Total,
		resultBlock.Block.LastCommit.BlockID.PartSetHeader.Hash)

	// Insert transactions if they exist
	for _, tx := range resultBlock.Block.Data.Txs {
		_, err := c.InsertTransaction(resultBlock.Block.Header.Height, tx)
		if err != nil {
			return false, err
		}
	}

	// Insert last commit signatures
	for _, signature := range resultBlock.Block.LastCommit.Signatures {
		_, err := c.InsertSignature(resultBlock.Block.LastCommit.Height, signature)
		if err != nil {
			return false, err
		}
	}

	// Insert evidences
	for _, evidence := range resultBlock.Block.Evidence.Evidence {
		switch ev := evidence.(type) {
		case *types.DuplicateVoteEvidence:
			var dve *types.DuplicateVoteEvidence
			dve = ev
			_, err = c.InsertDuplicateVoteEvidence(resultBlock.Block.Header.Height, dve)
			if err != nil {
				return false, err
			}
		case *types.LightClientAttackEvidence:
			fmt.Printf("Light Client Attack")
		default:
			fmt.Printf("Evidence not supported")
		}
	}

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (c *PostgresStorage) InsertTransaction(height int64, tx types.Tx) (bool, error) {
	_, err := c.Connection.Exec("INSERT INTO comet.block_data (height, transaction) values ($1,$2)",
		height,
		tx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (c *PostgresStorage) InsertDuplicateVoteEvidence(height int64, dve *types.DuplicateVoteEvidence) (bool, error) {

	//TODO: Find how to get the evidence type property e.g. 'tendermint/DuplicateVoteEvidence'

	_, err := c.Connection.Exec("INSERT INTO comet.duplicate_vote_evidence ("+
		"height, "+
		"evidence_type, "+
		"vote_a_type, "+
		"vote_a_height, "+
		"vote_a_round, "+
		"vote_a_block_id_hash, "+
		"vote_a_block_id_parts_hash, "+
		"vote_a_block_id_parts_total, "+
		"vote_a_timestamp, "+
		"vote_a_validator_address, "+
		"vote_a_validator_index, "+
		"vote_a_signature, "+
		"vote_b_type, "+
		"vote_b_height, "+
		"vote_b_round, "+
		"vote_b_block_id_hash, "+
		"vote_b_block_id_parts_hash, "+
		"vote_b_block_id_parts_total, "+
		"vote_b_timestamp, "+
		"vote_b_validator_address, "+
		"vote_b_validator_index, "+
		"vote_b_signature, "+
		"total_voting_power, "+
		"validator_voting_power, "+
		"evidence_timestamp) "+
		"VALUES ($1,$2,$3,$4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)",
		height,
		"tendermint/DuplicateVoteEvidence",
		dve.VoteA.Type,
		dve.VoteA.Height,
		dve.VoteA.Round,
		dve.VoteA.BlockID.Hash,
		dve.VoteA.BlockID.PartSetHeader.Hash,
		dve.VoteA.BlockID.PartSetHeader.Total,
		dve.VoteA.Timestamp,
		dve.VoteA.ValidatorAddress,
		dve.VoteA.ValidatorIndex,
		dve.VoteA.Signature,
		dve.VoteB.Type,
		dve.VoteB.Height,
		dve.VoteB.Round,
		dve.VoteB.BlockID.Hash,
		dve.VoteB.BlockID.PartSetHeader.Hash,
		dve.VoteB.BlockID.PartSetHeader.Total,
		dve.VoteB.Timestamp,
		dve.VoteB.ValidatorAddress,
		dve.VoteB.ValidatorIndex,
		dve.VoteB.Signature,
		dve.TotalVotingPower,
		dve.ValidatorPower,
		dve.Timestamp)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (c *PostgresStorage) InsertSignature(height int64, signature types.CommitSig) (bool, error) {
	_, err := c.Connection.Exec("INSERT INTO comet.last_commit_signature (height, block_id_flag, validator_address, signature_timestamp, signature) values ($1,$2,$3,$4,$5)",
		height,
		signature.BlockIDFlag,
		signature.ValidatorAddress,
		signature.Timestamp,
		signature.Signature)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (c *PostgresStorage) GetBlock(height int64) (ctypes.ResultBlock, error) {
	resultBlock := ctypes.ResultBlock{}
	lastCommit := types.Commit{}
	bId := types.BlockID{}
	b := new(types.Block)

	row := c.Connection.QueryRow("SELECT "+
		"block_id_hash, "+
		"block_id_parts_hash, "+
		"block_id_parts_total, "+
		"block_header_height, "+
		"block_header_chain_id, "+
		"block_header_block_time, "+
		"block_header_version_block, "+
		"block_header_version_app, "+
		"block_header_last_commit_hash, "+
		"block_header_data_hash, "+
		"block_header_validators_hash, "+
		"block_header_next_validators_hash, "+
		"block_header_consensus_hash, "+
		"block_header_app_hash, "+
		"block_header_last_results_hash, "+
		"block_header_evidence_hash, "+
		"block_header_proposer_address, "+
		"block_header_last_block_id_hash, "+
		"block_header_last_block_id_part_total, "+
		"block_header_last_block_id_parts_hash, "+
		"block_last_commit_height, "+
		"block_last_commit_round, "+
		"block_last_commit_block_id_hash, "+
		"block_last_commit_block_id_parts_total, "+
		"block_last_commit_block_id_parts_hash "+
		"FROM comet.result_block WHERE block_header_height=$1", height)
	err := row.Scan(
		&bId.Hash,
		&bId.PartSetHeader.Hash,
		&bId.PartSetHeader.Total,
		&b.Header.Height,
		&b.Header.ChainID,
		&b.Header.Time,
		&b.Header.Version.Block,
		&b.Header.Version.App,
		&b.Header.LastCommitHash,
		&b.Header.DataHash,
		&b.Header.ValidatorsHash,
		&b.Header.NextValidatorsHash,
		&b.Header.ConsensusHash,
		&b.Header.AppHash,
		&b.Header.LastResultsHash,
		&b.Header.EvidenceHash,
		&b.Header.ProposerAddress,
		&b.LastBlockID.Hash,
		&b.LastBlockID.PartSetHeader.Total,
		&b.LastBlockID.PartSetHeader.Hash,
		&lastCommit.Height, // *Commit
		&lastCommit.Round,
		&lastCommit.BlockID.Hash,
		&lastCommit.BlockID.PartSetHeader.Total,
		&lastCommit.BlockID.PartSetHeader.Hash)
	if err != nil {
		return resultBlock, err
	}
	b.LastCommit = &lastCommit
	resultBlock.BlockID = bId
	resultBlock.Block = b

	// Retrieve transactions if any
	var txBytes []byte
	txs, err := c.Connection.Query("SELECT transaction FROM comet.block_data WHERE height=$1", height)
	if err != nil {
		return resultBlock, err
	}
	defer txs.Close()
	for txs.Next() {
		err := txs.Scan(&txBytes)
		if err != nil {
			return resultBlock, err
		} else {
			resultBlock.Block.Data.Txs = append(resultBlock.Block.Data.Txs, txBytes)
		}
	}

	// Retrieve commit signatures
	var signature types.CommitSig
	signatures, err := c.Connection.Query("SELECT block_id_flag, validator_address, signature_timestamp, signature FROM comet.last_commit_signature WHERE height=$1", height-1)
	if err != nil {
		return resultBlock, err
	}
	defer signatures.Close()
	for signatures.Next() {
		err := signatures.Scan(&signature.BlockIDFlag, &signature.ValidatorAddress, &signature.Timestamp, &signature.Signature)
		if err != nil {
			return resultBlock, err
		} else {
			resultBlock.Block.LastCommit.Signatures = append(resultBlock.Block.LastCommit.Signatures, signature)
		}
	}

	// Retrieve duplicate vote evidences
	dve := types.DuplicateVoteEvidence{}
	dves, err := c.Connection.Query("SELECT "+
		"vote_a_type "+
		"FROM comet.duplicate_vote_evidence "+
		"WHERE height=$1", height)
	if err != nil {
		return resultBlock, err
	}
	defer dves.Close()
	for dves.Next() {
		voteA := types.Vote{}
		err := dves.Scan(&voteA.Type)
		if err != nil {
			return resultBlock, err
		} else {
			dve.VoteA = &voteA
			resultBlock.Block.Evidence.Evidence = append(resultBlock.Block.Evidence.Evidence, &dve)
		}
	}

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

func InsertBlocks(storage PostgresStorage) {

	// Ingest server
	fetcher := CometFetcher{
		Endpoint: os.Getenv("COMPANION_NODE_RPC"),
	}

	numberHeights := int64(3)
	initialHeightParameter := os.Getenv("COMPANION_INITIAL_HEIGHT")
	initialHeight, err := strconv.ParseInt(initialHeightParameter, 10, 64)
	if err != nil {
		fmt.Printf("Invalid initial height %s: %s\n", initialHeightParameter, err)
	}

	for height := initialHeight; height <= initialHeight+numberHeights; height++ {

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
}
