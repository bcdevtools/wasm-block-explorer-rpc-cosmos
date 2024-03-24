package types

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type basicQueryClient struct {
	tx.ServiceClient

	BankQueryClient banktypes.QueryClient
}

// newBasicQueryClient creates a new basic gRPC query clients
func newBasicQueryClient(clientCtx client.Context) *basicQueryClient {
	return &basicQueryClient{
		ServiceClient:   tx.NewServiceClient(clientCtx),
		BankQueryClient: banktypes.NewQueryClient(clientCtx),
	}
}
