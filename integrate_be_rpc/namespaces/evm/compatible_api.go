//go:build !be_json_rpc_evm

package evm

import (
	iberpcbackend "github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/backend/evm"
	"github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/compatible"
	"github.com/cosmos/cosmos-sdk/server"
)

// API is the EVM Block Explorer JSON-RPC.
type API struct {
}

// NewEvmBeAPI creates an instance of the EVM Block Explorer API.
func NewEvmBeAPI(
	ctx *server.Context,
	backend iberpcbackend.EvmBackendI,
) *API {
	compatible.PanicInvalidBuildTag()
	return nil
}
