package types

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// QueryClient defines a gRPC Client
type QueryClient struct {
	tx.ServiceClient

	BankQueryClient banktypes.QueryClient
	WasmQueryClient wasmtypes.QueryClient
}

// NewQueryClient creates a new gRPC query clients
func NewQueryClient(clientCtx client.Context) *QueryClient {
	queryClient := &QueryClient{
		ServiceClient:   tx.NewServiceClient(clientCtx),
		BankQueryClient: banktypes.NewQueryClient(clientCtx),
		WasmQueryClient: wasmtypes.NewQueryClient(clientCtx),
	}
	return queryClient
}
