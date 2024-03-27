//go:build !be_json_rpc_wasm

package wasm

import (
	iwberpcbackend "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/wasm"
	"github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/compatible"
	"github.com/cosmos/cosmos-sdk/server"
)

// API is the Wasm Block Explorer JSON-RPC.
type API struct {
}

// NewWasmBeAPI creates an instance of the Wasm Block Explorer API.
func NewWasmBeAPI(
	ctx *server.Context,
	backend iwberpcbackend.WasmBackendI,
) *API {
	compatible.PanicInvalidBuildTag()
	return nil
}
