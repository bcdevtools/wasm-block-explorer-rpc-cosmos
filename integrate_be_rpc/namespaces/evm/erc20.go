//go:build be_json_rpc_evm

package evm

import (
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/ethereum/go-ethereum/common"
)

func (api *API) GetErc20ContractInfo(address common.Address) (berpctypes.GenericBackendResponse, error) {
	api.logger.Debug("evm_getErc20ContractInfo")
	return api.backend.GetErc20ContractInfo(address)
}
