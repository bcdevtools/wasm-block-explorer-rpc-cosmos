//go:build be_json_rpc_wasm

package wasm

import (
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GetErc20ContractInfo will return information of ERC20 contract by address.
//   - name (optional): the name of the ERC20 token.
//   - symbol (mandatory): the symbol of the ERC20 token.
//   - decimals (mandatory): the number of decimals the token uses.
//
// If failed to query any of the mandatory fields, it will return an error.
// If failed to query the optional field, it will continue.
func (m *WasmBackend) GetErc20ContractInfo(contractAddress common.Address) (berpctypes.GenericBackendResponse, error) {
	// TODO BE: implement
	return nil, nil
}

func (m *WasmBackend) GetErc20Balance(accountAddress common.Address, contractAddresses []common.Address) (berpctypes.GenericBackendResponse, error) {
	// TODO BE: implement
	return nil, nil
}

func (m *WasmBackend) EvmCall(input string, contract common.Address, optionalChainId *hexutil.Big, optionalBlockNumber *hexutil.Uint64, optionalGas uint64) ([]byte, error) {
	// TODO BE: implement
	return nil, nil
}
