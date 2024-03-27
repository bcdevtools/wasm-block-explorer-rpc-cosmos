//go:build !be_json_rpc_evm && !be_json_rpc_wasm

package backend

import (
	berpcbackend "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/backend"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	ieberpcbackend "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/evm"
	iwberpcbackend "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/wasm"
)

/**
This file is used to get rid of compile error in IDE or Non-EVM & Non-Wasm chains.
*/

type DefaultRequestInterceptor struct {
	beRpcBackend berpcbackend.BackendI
}

func NewDefaultRequestInterceptor(
	beRpcBackend berpcbackend.BackendI,
	_ ieberpcbackend.EvmBackendI,
	_ iwberpcbackend.WasmBackendI,
) *DefaultRequestInterceptor {
	return &DefaultRequestInterceptor{
		beRpcBackend: beRpcBackend,
	}
}

func (m *DefaultRequestInterceptor) GetTransactionByHash(hashStr string) (intercepted bool, response berpctypes.GenericBackendResponse, err error) {
	intercepted = false
	return
}

func (m *DefaultRequestInterceptor) GetDenomsInformation() (intercepted, append bool, denoms map[string]string, err error) {
	intercepted = false
	append = false
	return
}

func (m *DefaultRequestInterceptor) GetModuleParams(moduleName string) (intercepted bool, res berpctypes.GenericBackendResponse, err error) {
	intercepted = false
	return
}

func (m *DefaultRequestInterceptor) GetAccount(accountAddressStr string) (intercepted, append bool, response berpctypes.GenericBackendResponse, err error) {
	intercepted = false
	append = false
	return
}
