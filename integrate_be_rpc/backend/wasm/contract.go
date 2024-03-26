//go:build be_json_rpc_wasm

package wasm

import (
	"github.com/ethereum/go-ethereum/common"
)

func (m *WasmBackend) GetContractCode(contractAddress common.Address) ([]byte, error) {
	// TODO BE: implement
	return nil, nil
}
