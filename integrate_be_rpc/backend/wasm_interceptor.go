//go:build be_json_rpc_wasm

package backend

import (
	berpcbackend "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/backend"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	iwberpcbackend "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/wasm"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DefaultRequestInterceptor struct {
	beRpcBackend berpcbackend.BackendI
	backend      iwberpcbackend.WasmBackendI
	bech32Cfg    berpctypes.Bech32Config
}

func NewDefaultRequestInterceptor(
	beRpcBackend berpcbackend.BackendI,
	backend iwberpcbackend.WasmBackendI,
) *DefaultRequestInterceptor {
	return &DefaultRequestInterceptor{
		beRpcBackend: beRpcBackend,
		backend:      backend,
		bech32Cfg:    berpctypes.NewBech32Config(),
	}
}

func (m *DefaultRequestInterceptor) GetTransactionByHash(hashStr string) (intercepted bool, response berpctypes.GenericBackendResponse, err error) {
	// handle WASM txs, otherwise return false
	intercepted = false
	// TODO BE: implement
	return
}

func (m *DefaultRequestInterceptor) GetDenomsInformation() (intercepted, append bool, denoms map[string]string, err error) {
	intercepted = false
	// TODO BE: implement
	return
}

func (m *DefaultRequestInterceptor) GetModuleParams(moduleName string) (intercepted bool, res berpctypes.GenericBackendResponse, err error) {
	var params any

	switch moduleName {
	case "wasm", "cosmwasm":
		wasmParams, errFetch := m.backend.GetWasmModuleParams()
		if errFetch != nil {
			err = errors.Wrap(errFetch, "failed to get wasm params")
		} else {
			params = *wasmParams
		}
		break
	default:
		intercepted = false
		return
	}

	if err != nil {
		return
	}

	res, err = berpctypes.NewGenericBackendResponseFrom(params)
	if err != nil {
		err = status.Error(codes.Internal, errors.Wrap(err, "module params").Error())
		return
	}

	intercepted = true
	return
}

// GetAccount returns the contract information if the account is a contract. Other-wise no-op.
func (m *DefaultRequestInterceptor) GetAccount(accountAddressStr string) (intercepted, append bool, response berpctypes.GenericBackendResponse, err error) {
	intercepted = false
	// TODO BE: implement
	return
}
