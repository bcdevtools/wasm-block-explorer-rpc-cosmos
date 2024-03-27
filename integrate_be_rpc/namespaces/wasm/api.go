package wasm

import (
	"fmt"
	iwberpcbackend "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/wasm"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/tendermint/tendermint/libs/log"
)

// RPC namespaces and API version
const (
	DymWasmBlockExplorerNamespace = "wasm"

	ApiVersion = "1.0"
)

// API is the Wasm Block Explorer JSON-RPC.
type API struct {
	ctx     *server.Context
	logger  log.Logger
	backend iwberpcbackend.WasmBackendI
}

// NewWasmBeAPI creates an instance of the Wasm Block Explorer API.
func NewWasmBeAPI(
	ctx *server.Context,
	backend iwberpcbackend.WasmBackendI,
) *API {
	return &API{
		ctx:     ctx,
		logger:  ctx.Logger.With("api", "wasm"),
		backend: backend,
	}
}

func (api *API) Echo(text string) string {
	api.logger.Debug("wasm_echo")
	return fmt.Sprintf("hello \"%s\" from Wasm Block Explorer API", text)
}
