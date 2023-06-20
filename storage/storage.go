package storage

import (
	"database/sql"
	"github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
)

type IStorage interface {
	Connect() (*sql.DB, error)
	Disconnect(db *sql.DB) error
	Ping() error
	InsertBlock(resultBlock coretypes.ResultBlock) (bool, error)
	InsertTransaction(height int64, tx types.Tx) (bool, error)
	InsertDuplicateVoteEvidence(height int64, evidence *types.DuplicateVoteEvidence) (bool, error)
	InsertLightClientAttackEvidence(height int64, evidence *types.LightClientAttackEvidence) (bool, error)
	InsertCommitSignature(height int64, commitSig types.CommitSig) (bool, error)
	GetBlock(height int64) (coretypes.ResultBlock, error)
}
