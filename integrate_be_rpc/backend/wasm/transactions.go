//go:build be_json_rpc_wasm

package wasm

import (
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/ethereum/go-ethereum/common"
)

func (m *WasmBackend) GetWasmTransactionInvolversByHash(hash common.Hash) (berpctypes.MessageInvolversResult, error) {
	// TODO BE: implement
	return nil, nil
}

func (m *WasmBackend) GetWasmTransactionByHash(hash common.Hash) (berpctypes.GenericBackendResponse, error) {
	// TODO BE: implement
	return nil, nil
}
