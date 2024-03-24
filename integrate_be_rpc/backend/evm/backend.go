//go:build be_json_rpc_evm

package evm

import (
	"context"
	"github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/config"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	iberpctypes "github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
	"github.com/tendermint/tendermint/libs/log"
)

var _ EvmBackendI = (*EvmBackend)(nil)

type EvmBackendI interface {
	// Transactions

	// GetEvmTransactionByHash returns a transaction by its hash.
	GetEvmTransactionByHash(hash common.Hash) (berpctypes.GenericBackendResponse, error)

	// Misc

	GetEvmModuleParams() (*evmtypes.Params, error)
}

// EvmBackend implements the EvmBackendI interface
type EvmBackend struct {
	ctx         context.Context
	clientCtx   client.Context
	queryClient *iberpctypes.QueryClient // gRPC query client
	logger      log.Logger
	cfg         config.BeJsonRpcConfig
}

// NewEvmBackend creates a new EvmBackend instance for EVM Block Explorer
func NewEvmBackend(
	ctx *server.Context,
	logger log.Logger,
	clientCtx client.Context,
) *EvmBackend {
	appConf, err := config.GetConfig(ctx.Viper)
	if err != nil {
		panic(err)
	}

	return &EvmBackend{
		ctx:         context.Background(),
		clientCtx:   clientCtx,
		queryClient: iberpctypes.NewQueryClient(clientCtx),
		logger:      logger.With("module", "i_be_rpc"),
		cfg:         appConf,
	}
}
