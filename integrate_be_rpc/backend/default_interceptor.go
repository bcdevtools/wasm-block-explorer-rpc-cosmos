//go:build !be_json_rpc_evm

package backend

import (
	berpcbackend "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/backend"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	iberpcbackend "github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/backend/evm"
)

/**
This file is used to get rid of compile error in IDE or Non-EVM chains.
*/

type DefaultRequestInterceptor struct {
	beRpcBackend berpcbackend.BackendI
}

func NewDefaultRequestInterceptor(
	beRpcBackend berpcbackend.BackendI,
	_ iberpcbackend.EvmBackendI,
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
	return
}

func (m *DefaultRequestInterceptor) GetModuleParams(moduleName string) (intercepted bool, res berpctypes.GenericBackendResponse, err error) {
	intercepted = false
	return
}
