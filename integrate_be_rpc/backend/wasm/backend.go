package wasm

import (
	"context"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/config"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	iberpctypes "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

var _ WasmBackendI = (*WasmBackend)(nil)

type WasmBackendI interface {
	// Transactions

	// GetWasmTransactionByHash returns a transaction by its hash.
	GetWasmTransactionByHash(hash string) (berpctypes.GenericBackendResponse, error)

	GetWasmTransactionInvolversByHash(hash string) (berpctypes.MessageInvolversResult, error)

	GetTmTxResult(tmTx tmtypes.Tx) ([]abci.Event, error)

	// CW-20

	GetCw20ContractInfo(contractAddress string) (berpctypes.GenericBackendResponse, error)

	GetCw20TokenInfo(contractAddress string) (*iberpctypes.Cw20TokenInfo, error)

	GetCw20Balance(accountAddress string, contractAddresses []string) (berpctypes.GenericBackendResponse, error)

	// Wasm

	SmartContractState(input map[string]any, contract string, optionalBlockNumber *int64) ([]byte, error)

	GetContractCodeId(contractAddress string) (uint64, error)

	// Misc

	GetWasmModuleParams() (*wasmtypes.Params, error)
}

// WasmBackend implements the WasmBackendI interface
type WasmBackend struct {
	ctx         context.Context
	clientCtx   client.Context
	queryClient *iberpctypes.QueryClient // gRPC query client
	logger      log.Logger
	cfg         config.BeJsonRpcConfig
}

// NewWasmBackend creates a new WasmBackend instance for Wasm Block Explorer
func NewWasmBackend(
	ctx *server.Context,
	logger log.Logger,
	clientCtx client.Context,
	_ berpctypes.ExternalServices,
) *WasmBackend {
	appConf, err := config.GetConfig(ctx.Viper)
	if err != nil {
		panic(err)
	}

	return &WasmBackend{
		ctx:         context.Background(),
		clientCtx:   clientCtx,
		queryClient: iberpctypes.NewQueryClient(clientCtx),
		logger:      logger.With("module", "wasm_be_rpc"),
		cfg:         appConf,
	}
}
