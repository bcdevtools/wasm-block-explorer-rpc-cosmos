package wasm

import (
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
)

func (api *API) GetCw20ContractInfo(contractAddress string) (berpctypes.GenericBackendResponse, error) {
	api.logger.Debug("wasm_getCw20ContractInfo")
	return api.backend.GetCw20ContractInfo(contractAddress)
}

func (api *API) GetCw20Balance(accountAddress string, contractAddresses []string) (berpctypes.GenericBackendResponse, error) {
	api.logger.Debug("wasm_getCw20Balance")
	return api.backend.GetCw20Balance(accountAddress, contractAddresses)
}
