//go:build !be_json_rpc_wasm

package wasm

import (
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/compatible"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"
)

/**
This file is used to get rid of compile error in IDE.
*/

var _ WasmBackendI = (*WasmBackend)(nil)

type WasmBackendI interface {
	GetWasmTransactionInvolversByHash(hash common.Hash) (berpctypes.MessageInvolversResult, error)
}

type WasmBackend struct {
}

// NewWasmBackend creates a new WasmBackend instance for EVM Block Explorer.
// This method is for get rid of build error in IDE in final chains.
func NewWasmBackend(
	ctx *server.Context,
	logger log.Logger,
	clientCtx client.Context,
	externalServices berpctypes.ExternalServices,
) *WasmBackend {
	compatible.PanicInvalidBuildTag()
	return nil
}

func (m *WasmBackend) GetWasmTransactionInvolversByHash(hash common.Hash) (berpctypes.MessageInvolversResult, error) {
	return nil, nil
}
