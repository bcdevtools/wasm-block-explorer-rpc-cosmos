//go:build !be_json_rpc_wasm

package types

import (
	"github.com/cosmos/cosmos-sdk/client"
)

/**
This file is used to get rid of compile error in IDE or Non-EVM & Non-Wasm chains.
*/

// QueryClient defines a gRPC Client
type QueryClient struct {
	*basicQueryClient
}

// NewQueryClient creates a new gRPC query clients
func NewQueryClient(clientCtx client.Context) *QueryClient {
	queryClient := &QueryClient{
		basicQueryClient: newBasicQueryClient(clientCtx),
	}
	return queryClient
}
