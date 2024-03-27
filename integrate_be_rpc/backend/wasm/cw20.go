//go:build be_json_rpc_wasm

package wasm

import (
	"encoding/json"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	berpcutils "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/utils"
	iberpctypes "github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
)

// GetCw20ContractInfo will return information of CW-20 contract by address.
//   - name (optional): the name of the CW-20 token.
//   - symbol (mandatory): the symbol of the CW-20 token.
//   - decimals (mandatory): the number of decimals the token uses.
//
// If failed to query any of the mandatory fields, it will return an error.
// If failed to query the optional field, it will continue.
func (m *WasmBackend) GetCw20ContractInfo(contractAddress string) (berpctypes.GenericBackendResponse, error) {
	tokenInfo, err := m.GetCw20TokenInfo(contractAddress)
	if err != nil {
		return nil, err
	}

	return berpctypes.GenericBackendResponse{
		"name":        tokenInfo.Name,
		"symbol":      tokenInfo.Symbol,
		"decimals":    tokenInfo.Decimals,
		"totalSupply": tokenInfo.TotalSupply,
	}, nil
}

func (m *WasmBackend) GetCw20Balance(accountAddress string, contractAddresses []string) (berpctypes.GenericBackendResponse, error) {
	// TODO BE: implement
	return nil, nil
}

func (m *WasmBackend) SmartContractState(input map[string]any, contract string, optionalBlockNumber *int64) ([]byte, error) {
	ctx := m.ctx
	if optionalBlockNumber != nil {
		height := *optionalBlockNumber
		if height > 0 {
			ctx = berpcutils.QueryContextWithHeight(*optionalBlockNumber)
		}
	}

	bz, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	resState, err := m.queryClient.WasmQueryClient.SmartContractState(ctx, &wasmtypes.QuerySmartContractStateRequest{
		Address:   contract,
		QueryData: bz,
	})

	if err != nil {
		return nil, err
	}

	return resState.Data, nil
}

func (m *WasmBackend) GetCw20TokenInfo(contractAddress string) (*iberpctypes.Cw20TokenInfo, error) {
	codeId, err := m.GetContractCodeId(contractAddress)
	if err != nil {
		return nil, err
	}

	if codeId == 0 {
		return nil, status.Error(codes.NotFound, errors.New(contractAddress+" is not a contract").Error())
	}

	state, err := m.SmartContractState(map[string]any{
		"token_info": map[string]any{},
	}, contractAddress, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get contract state").Error())
	}
	if len(state) < 1 {
		return nil, status.Error(codes.NotFound, errors.New("no response contract state").Error())
	}

	defer func() {
		m.logger.Error(string(state))
	}()

	var data iberpctypes.Cw20TokenInfo

	err = json.Unmarshal(state, &data)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to unmarshal response").Error())
	}

	if len(data.Name) == 0 && len(data.Symbol) == 0 {
		return nil, status.Error(codes.NotFound, errors.New("no token info found").Error())
	}

	if len(data.TotalSupplyStr) > 0 {
		totalSupply, ok := new(big.Int).SetString(data.TotalSupplyStr, 10)
		if !ok {
			return nil, status.Error(codes.Internal, errors.New("failed to parse total supply").Error())
		}
		data.TotalSupply = totalSupply
	}

	return &data, nil
}
