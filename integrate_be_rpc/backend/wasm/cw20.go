//go:build be_json_rpc_wasm

package wasm

import (
	"encoding/json"
	"fmt"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	berpcutils "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/utils"
	iberpctypes "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
)

// GetCw20ContractInfo will return information of CW-20 contract by address.
//   - name: the name of the CW-20 token.
//   - symbol: the symbol of the CW-20 token.
//   - decimals: the number of decimals the token uses.
//   - totalSupply: the total supply of the token.
//
// If failed to query any of the mandatory fields, it will return an error.
// If failed to query the optional field, it will continue.
func (m *WasmBackend) GetCw20ContractInfo(contractAddress string) (berpctypes.GenericBackendResponse, error) {
	tokenInfo, err := m.GetCw20TokenInfo(contractAddress)
	if err != nil {
		return nil, err
	}

	res := berpctypes.GenericBackendResponse{
		"name":     tokenInfo.Name,
		"symbol":   tokenInfo.Symbol,
		"decimals": tokenInfo.Decimals,
	}

	if tokenInfo.TotalSupply != nil {
		res["totalSupply"] = tokenInfo.TotalSupply.String()
	}

	return res, nil
}

func (m *WasmBackend) GetCw20Balance(accountAddress string, contractAddresses []string) (berpctypes.GenericBackendResponse, error) {
	res := berpctypes.GenericBackendResponse{
		"account": accountAddress,
	}

	resForContracts := make([]berpctypes.GenericBackendResponse, 0)

	for _, contractAddress := range contractAddresses {
		var display string
		var decimals uint8
		var balance *big.Int

		codeId, err := m.GetContractCodeId(contractAddress)
		if err != nil {
			return nil, err
		}

		if codeId == 0 {
			display = ""
			decimals = 0
			balance = big.NewInt(0)
		} else {
			tokenInfo, err := m.GetCw20TokenInfo(contractAddress)
			if err != nil {
				return nil, err
			}

			if len(tokenInfo.Symbol) > 0 {
				display = tokenInfo.Symbol
			} else if len(tokenInfo.Name) > 0 {
				display = tokenInfo.Name
			} else {
				display = fmt.Sprintf("(%s)", contractAddress) // force value
			}

			decimals = tokenInfo.Decimals

			state, err := m.SmartContractState(map[string]any{
				"balance": map[string]any{
					"address": accountAddress,
				},
			}, contractAddress, nil)
			if err != nil {
				return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get contract state").Error())
			}
			if len(state) < 1 {
				return nil, status.Error(codes.NotFound, errors.New("no response contract state").Error())
			}

			var data struct {
				Balance string `json:"balance"`
			}

			err = json.Unmarshal(state, &data)
			if err != nil {
				return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to unmarshal response").Error())
			}

			var ok bool
			balance, ok = new(big.Int).SetString(data.Balance, 10)
			if !ok {
				return nil, status.Error(codes.Internal, errors.New("failed to parse balance "+data.Balance).Error())
			}
		}

		resForContracts = append(resForContracts, berpctypes.GenericBackendResponse{
			"contract": contractAddress,
			"display":  display,
			"decimals": decimals,
			"balance":  balance.String(),
		})
	}

	res["cw20Balances"] = resForContracts

	return res, nil
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
