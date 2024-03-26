//go:build be_json_rpc_wasm

package types

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// QueryClient defines a gRPC Client
type QueryClient struct {
	*basicQueryClient
	WasmQueryClient wasmtypes.QueryClient
}

// NewQueryClient creates a new gRPC query clients
func NewQueryClient(clientCtx client.Context) *QueryClient {
	queryClient := &QueryClient{
		basicQueryClient: newBasicQueryClient(clientCtx),
		WasmQueryClient:  wasmtypes.NewQueryClient(clientCtx),
	}
	return queryClient
}
