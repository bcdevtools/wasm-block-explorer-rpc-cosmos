//go:build be_json_rpc_wasm

package wasm

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

func (m *WasmBackend) GetContractCodeId(contractAddress string) (uint64, error) {
	resContractInfo, err := m.queryClient.WasmQueryClient.ContractInfo(m.ctx, &wasmtypes.QueryContractInfoRequest{
		Address: contractAddress,
	})
	if err != nil {
		return 0, err
	}
	return resContractInfo.CodeID, nil
}
