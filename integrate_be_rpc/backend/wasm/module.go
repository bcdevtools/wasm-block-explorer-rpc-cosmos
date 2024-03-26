//go:build be_json_rpc_wasm

package wasm

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

func (m *WasmBackend) GetWasmModuleParams() (*wasmtypes.Params, error) {
	res, err := m.queryClient.WasmQueryClient.Params(m.ctx, &wasmtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	return &res.Params, nil
}
