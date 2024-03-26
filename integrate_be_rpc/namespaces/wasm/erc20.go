//go:build be_json_rpc_wasm

package wasm

import (
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/ethereum/go-ethereum/common"
)

func (api *API) GetErc20ContractInfo(contractAddress common.Address) (berpctypes.GenericBackendResponse, error) {
	api.logger.Debug("wasm_getErc20ContractInfo")
	return api.backend.GetErc20ContractInfo(contractAddress)
}

func (api *API) GetErc20Balance(accountAddress common.Address, contractAddresses []common.Address) (berpctypes.GenericBackendResponse, error) {
	api.logger.Debug("wasm_getErc20Balance")
	return api.backend.GetErc20Balance(accountAddress, contractAddresses)
}
