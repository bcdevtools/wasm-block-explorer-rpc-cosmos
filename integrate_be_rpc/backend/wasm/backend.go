//go:build be_json_rpc_wasm

package wasm

import (
	"context"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/config"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	iberpctypes "github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/libs/log"
)

var _ WasmBackendI = (*WasmBackend)(nil)

type WasmBackendI interface {
	// Transactions

	// GetWasmTransactionByHash returns a transaction by its hash.
	GetWasmTransactionByHash(hash common.Hash) (berpctypes.GenericBackendResponse, error)

	GetWasmTransactionInvolversByHash(hash common.Hash) (berpctypes.MessageInvolversResult, error)

	// ERC-20

	GetErc20ContractInfo(contractAddress common.Address) (berpctypes.GenericBackendResponse, error)

	GetErc20Balance(accountAddress common.Address, contractAddresses []common.Address) (berpctypes.GenericBackendResponse, error)

	// EVM

	EvmCall(input string, contract common.Address, optionalChainId *hexutil.Big, optionalBlockNumber *hexutil.Uint64, optionalGas uint64) ([]byte, error)

	GetContractCode(contractAddress common.Address) ([]byte, error)

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
