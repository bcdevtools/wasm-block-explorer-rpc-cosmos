//go:build be_json_rpc_evm

package evm

import (
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/ethereum/go-ethereum/common"
)

func (m *EvmBackend) GetEvmTransactionByHash(hash common.Hash) (berpctypes.GenericBackendResponse, error) {
	return berpctypes.GenericBackendResponse{
		"Hello": "Vietnam",
		"hash":  hash,
	}, nil
}
