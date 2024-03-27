//go:build be_json_rpc_wasm

package wasm

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (m *WasmBackend) GetContractCodeId(contractAddress string) (uint64, error) {
	resContractInfo, err := m.queryClient.WasmQueryClient.ContractInfo(m.ctx, &wasmtypes.QueryContractInfoRequest{
		Address: contractAddress,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no such contract") {
			return 0, nil
		}
		return 0, status.Error(codes.Internal, errors.Wrap(err, "failed to get contract info").Error())
	}
	return resContractInfo.CodeID, nil
}
