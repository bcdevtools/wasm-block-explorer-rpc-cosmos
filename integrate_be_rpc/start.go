package integrate_be_rpc

import (
	"context"
	berpc "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc"
	berpcbackend "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/backend"
	berpccfg "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/config"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	berpcserver "github.com/bcdevtools/block-explorer-rpc-cosmos/server"
	wasmberpcbackend "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/wasm"
	bemsgparsers "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/message_parsers"
	wasmbeapi "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/namespaces/wasm"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/rpc"
	rpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"time"
)

type FuncRegister func(wasmberpcbackend.WasmBackendI)

func StartWasmBeJsonRPC(
	ctx *server.Context,
	clientCtx client.Context,
	chainId string,
	beRpcCfg berpccfg.BeJsonRpcConfig,
	externalServicesModifierFunc func(berpctypes.ExternalServices) berpctypes.ExternalServices,
	apiNamespaceRegisterFunc FuncRegister,
	customInterceptorCreationFunc func(berpcbackend.BackendI, wasmberpcbackend.WasmBackendI) berpcbackend.RequestInterceptor,
	tmRPCAddr, tmEndpoint string,
) (serverCloseDeferFunc func(), err error) {
	if err := beRpcCfg.Validate(); err != nil {
		return nil, err
	}

	clientCtx = clientCtx.WithChainID(chainId)

	externalServices := berpctypes.ExternalServices{
		ChainType: berpctypes.ChainTypeCosmWasm,
	}
	if externalServicesModifierFunc != nil {
		externalServices = externalServicesModifierFunc(externalServices)
	}

	wasmBeRpcBackend := wasmberpcbackend.NewWasmBackend(ctx, ctx.Logger, clientCtx, externalServices)

	berpc.RegisterAPINamespace(wasmbeapi.DymWasmBlockExplorerNamespace, func(ctx *server.Context,
		_ client.Context,
		_ *rpcclient.WSClient,
		_ map[string]berpctypes.MessageParser,
		_ map[string]berpctypes.MessageInvolversExtractor,
		_ func(berpcbackend.BackendI) berpcbackend.RequestInterceptor,
		_ berpctypes.ExternalServices,
	) []rpc.API {
		return []rpc.API{
			{
				Namespace: wasmbeapi.DymWasmBlockExplorerNamespace,
				Version:   wasmbeapi.ApiVersion,
				Service:   wasmbeapi.NewWasmBeAPI(ctx, wasmBeRpcBackend),
				Public:    true,
			},
		}
	}, false)

	if apiNamespaceRegisterFunc != nil {
		apiNamespaceRegisterFunc(wasmBeRpcBackend)
	}

	// register message parsers & message involvers extractor

	bemsgparsers.RegisterMessageParsersForWasm()

	var interceptorCreationFunc func(berpcbackend.BackendI) berpcbackend.RequestInterceptor
	if customInterceptorCreationFunc != nil {
		interceptorCreationFunc = func(backend berpcbackend.BackendI) berpcbackend.RequestInterceptor {
			return customInterceptorCreationFunc(backend, wasmBeRpcBackend)
		}
	}

	beJsonRpcHttpSrv, beJsonRpcHttpSrvDone, err := berpcserver.StartBeJsonRPC(
		ctx, clientCtx, tmRPCAddr, tmEndpoint,
		beRpcCfg,
		interceptorCreationFunc,
		externalServices,
	)
	if err != nil {
		return nil, err
	}

	return func() {
		shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFn()
		if err := beJsonRpcHttpSrv.Shutdown(shutdownCtx); err != nil {
			ctx.Logger.Error("Wasm Block Explorer Json-RPC HTTP server shutdown produced a warning", "error", err.Error())
		} else {
			ctx.Logger.Info("Wasm Block Explorer Json-RPC HTTP server shut down, waiting 5 sec")
			select {
			case <-time.Tick(5 * time.Second):
			case <-beJsonRpcHttpSrvDone:
			}
		}
	}, nil
}
