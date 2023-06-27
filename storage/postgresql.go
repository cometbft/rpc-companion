package storage

import (
	"context"
	"database/sql"
	"fmt"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	_ "github.com/lib/pq"
	"time"
)

type PostgresStorage struct {
	ConnectionString string
}

func (c *PostgresStorage) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.ConnectionString)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func (c *PostgresStorage) Disconnect(conn *sql.DB) error {
	err := conn.Close()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *PostgresStorage) Ping() error {
	conn, err := c.Connect()
	if err != nil {
		return err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = conn.PingContext(ctx)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (c *PostgresStorage) InsertBlock(resultBlock ctypes.ResultBlock) (bool, error) {
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return false, err
	} else {
		_, err := conn.Exec(
			"INSERT INTO comet.block_result"+
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
			_, err := c.InsertCommitSignature(resultBlock.Block.LastCommit.Height, signature)
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
				var lae *types.LightClientAttackEvidence
				lae = ev
				_, err = c.InsertLightClientAttackEvidence(resultBlock.Block.Header.Height, lae)
				if err != nil {
					return false, err
				}
			default:
				fmt.Printf("Evidence type not supported")
			}
		}

		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

func (c *PostgresStorage) InsertTransaction(height int64, tx types.Tx) (bool, error) {
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return false, err
	} else {
		_, err := conn.Exec("INSERT INTO comet.block_data (height, transaction) values ($1,$2)",
			height,
			tx)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

func (c *PostgresStorage) InsertDuplicateVoteEvidence(height int64, evidence *types.DuplicateVoteEvidence) (bool, error) {
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return false, err
	} else {
		//TODO: Find how to get the evidence type property e.g. 'tendermint/DuplicateVoteEvidence'
		_, err := conn.Exec("INSERT INTO comet.evidence_duplicate_vote ("+
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
			evidence.VoteA.Type,
			evidence.VoteA.Height,
			evidence.VoteA.Round,
			evidence.VoteA.BlockID.Hash,
			evidence.VoteA.BlockID.PartSetHeader.Hash,
			evidence.VoteA.BlockID.PartSetHeader.Total,
			evidence.VoteA.Timestamp,
			evidence.VoteA.ValidatorAddress,
			evidence.VoteA.ValidatorIndex,
			evidence.VoteA.Signature,
			evidence.VoteB.Type,
			evidence.VoteB.Height,
			evidence.VoteB.Round,
			evidence.VoteB.BlockID.Hash,
			evidence.VoteB.BlockID.PartSetHeader.Hash,
			evidence.VoteB.BlockID.PartSetHeader.Total,
			evidence.VoteB.Timestamp,
			evidence.VoteB.ValidatorAddress,
			evidence.VoteB.ValidatorIndex,
			evidence.VoteB.Signature,
			evidence.TotalVotingPower,
			evidence.ValidatorPower,
			evidence.Timestamp)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

func (c *PostgresStorage) InsertLightClientAttackEvidence(height int64, evidence *types.LightClientAttackEvidence) (bool, error) {
	evID := int64(0)
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return false, err
	} else {
		//TODO: Find how to get the evidence type property e.g. 'tendermint/LightClientAttackEvidence'
		err := conn.QueryRow("INSERT INTO comet.evidence_light_client_attack (height, evidence_type, common_height, total_voting_power, timestamp) VALUES ($1,$2,$3, $4, $5) RETURNING id",
			height,
			"tendermint/LightClientAttackEvidence",
			evidence.CommonHeight,
			evidence.TotalVotingPower,
			evidence.Timestamp).Scan(&evID)
		if err != nil {
			return false, err
		} else {
			// Insert byzantine validators
			for _, validator := range evidence.ByzantineValidators {
				valID, err := c.InsertValidator(validator)
				if err != nil {
					return false, err
				} else {
					//Insert an entry into a light client evidence to byzantine validator table (1-N)
					_, err := c.InsertEvidenceLCAByzantineValidator(evID, valID)
					if err != nil {
						return false, err
					}
				}
			}
			return true, nil
		}
	}
}

func (c *PostgresStorage) InsertValidator(v *types.Validator) (int64, error) {
	id := int64(0)
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return id, err
	} else {
		err := conn.QueryRow("INSERT INTO comet.validator (address, pub_key_type, pub_key_value, voting_power, proposer_priority) SELECT $1,$2,$3,$4,$5 WHERE NOT EXISTS (SELECT address FROM comet.validator WHERE address=$1 AND voting_power=$4) RETURNING id;",
			v.Address,
			v.PubKey.Type(),
			v.PubKey.Bytes(),
			v.VotingPower,
			v.ProposerPriority).Scan(&id)
		if err != nil {
			return id, err
		} else {
			return id, nil
		}
	}
}

func (c *PostgresStorage) InsertEvidenceLCAByzantineValidator(evID int64, valID int64) (sql.Result, error) {
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return nil, err
	} else {
		r, err := conn.Exec("INSERT INTO comet.evidence_light_client_attack_byzantine_validator(evidence_id, validator_id) VALUES ($1, $2);",
			evID,
			valID)
		if err != nil {
			return nil, err
		} else {
			return r, nil
		}
	}
}

func (c *PostgresStorage) InsertCommitSignature(height int64, commitSig types.CommitSig) (bool, error) {
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return false, err
	} else {
		_, err := conn.Exec("INSERT INTO comet.block_commit_sig (height, block_id_flag, validator_address, timestamp, signature) values ($1,$2,$3,$4,$5)",
			height,
			commitSig.BlockIDFlag,
			commitSig.ValidatorAddress,
			commitSig.Timestamp,
			commitSig.Signature)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

func (c *PostgresStorage) GetBlock(height int64) (ctypes.ResultBlock, error) {
	resultBlock := ctypes.ResultBlock{}
	conn, err := c.Connect()
	defer conn.Close()
	if err != nil {
		return resultBlock, err
	} else {
		lastCommit := types.Commit{}
		bId := types.BlockID{}
		b := new(types.Block)
		row := conn.QueryRow("SELECT "+
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
			"FROM comet.block_result WHERE block_header_height=$1", height)
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
		txs, err := conn.Query("SELECT transaction FROM comet.block_data WHERE height=$1", height)
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
		signatures, err := conn.Query("SELECT block_id_flag, validator_address, timestamp, signature FROM comet.block_commit_sig WHERE height=$1", resultBlock.Block.LastCommit.Height)
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
		dves, err := conn.Query("SELECT "+
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
			"evidence_timestamp "+
			"FROM comet.evidence_duplicate_vote "+
			"WHERE height=$1", height)
		if err != nil {
			return resultBlock, err
		}
		defer dves.Close()
		for dves.Next() {
			voteA := types.Vote{}
			voteB := types.Vote{}
			err := dves.Scan(
				&voteA.Type,
				&voteA.Height,
				&voteA.Round,
				&voteA.BlockID.Hash,
				&voteA.BlockID.PartSetHeader.Hash,
				&voteA.BlockID.PartSetHeader.Total,
				&voteA.Timestamp,
				&voteA.ValidatorAddress,
				&voteA.ValidatorIndex,
				&voteA.Signature,
				&voteB.Type,
				&voteB.Height,
				&voteB.Round,
				&voteB.BlockID.Hash,
				&voteB.BlockID.PartSetHeader.Hash,
				&voteB.BlockID.PartSetHeader.Total,
				&voteB.Timestamp,
				&voteB.ValidatorAddress,
				&voteB.ValidatorIndex,
				&voteB.Signature,
				&dve.TotalVotingPower,
				&dve.ValidatorPower,
				&dve.Timestamp)
			if err != nil {
				return resultBlock, err
			} else {
				dve.VoteA = &voteA
				dve.VoteB = &voteB
				resultBlock.Block.Evidence.Evidence = append(resultBlock.Block.Evidence.Evidence, &dve)
			}
		}

		// Retrieve light client attack evidences
		lcaev := types.LightClientAttackEvidence{}
		lcaevs, err := conn.Query("SELECT "+
			"id, "+
			"common_height, "+
			"total_voting_power, "+
			"timestamp "+
			"FROM comet.evidence_light_client_attack "+
			"WHERE height=$1", height)
		if err != nil {
			return resultBlock, err
		}
		defer lcaevs.Close()
		for lcaevs.Next() {
			var id int64
			err := lcaevs.Scan(
				&id,
				&lcaev.CommonHeight,
				&lcaev.TotalVotingPower,
				&lcaev.Timestamp)
			if err != nil {
				return resultBlock, err
			} else {
				//Retrieve byzantine validators
				// Retrieve commit signatures
				var validator types.Validator
				bValidators, err := conn.Query("SELECT V.address, v.voting_power, v.proposer_priority FROM comet.evidence_light_client_attack_byzantine_validator E JOIN comet.validator V ON E.validator_id = V.id WHERE E.evidence_id=$1;", id)
				if err != nil {
					return resultBlock, err
				}
				defer bValidators.Close()
				for bValidators.Next() {
					//TODO: Logic to retrieve pub_key
					err := bValidators.Scan(&validator.Address, &validator.VotingPower, &validator.ProposerPriority)
					if err != nil {
						return resultBlock, err
					} else {
						lcaev.ByzantineValidators = append(lcaev.ByzantineValidators, &validator)
					}
				}
				resultBlock.Block.Evidence.Evidence = append(resultBlock.Block.Evidence.Evidence, &lcaev)
			}
		}

		return resultBlock, err
	}
}
