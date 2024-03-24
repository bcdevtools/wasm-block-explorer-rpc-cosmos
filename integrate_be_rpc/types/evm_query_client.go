//go:build be_json_rpc_evm

package types

import (
	"github.com/cosmos/cosmos-sdk/client"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
)

// QueryClient defines a gRPC Client
type QueryClient struct {
	*basicQueryClient
	EvmQueryClient evmtypes.QueryClient
}

// NewQueryClient creates a new gRPC query clients
func NewQueryClient(clientCtx client.Context) *QueryClient {
	queryClient := &QueryClient{
		basicQueryClient: newBasicQueryClient(clientCtx),
		EvmQueryClient:   evmtypes.NewQueryClient(clientCtx),
	}
	return queryClient
}
