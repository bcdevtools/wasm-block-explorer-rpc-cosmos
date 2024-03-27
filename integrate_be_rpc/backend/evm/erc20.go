//go:build be_json_rpc_evm

package evm

import (
	"encoding/hex"
	"fmt"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	iberpcutils "github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmosrpctypes "github.com/evmos/evmos/v12/rpc/types"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"math/big"
	"strings"
)

// GetErc20ContractInfo will return information of ERC20 contract by address.
//   - name (optional): the name of the ERC20 token.
//   - symbol (mandatory): the symbol of the ERC20 token.
//   - decimals (mandatory): the number of decimals the token uses.
//
// If failed to query any of the mandatory fields, it will return an error.
// If failed to query the optional field, it will continue.
func (m *EvmBackend) GetErc20ContractInfo(contractAddress common.Address) (berpctypes.GenericBackendResponse, error) {
	resCode, err := m.queryClient.EvmQueryClient.Code(m.ctx, &evmtypes.QueryCodeRequest{
		Address: contractAddress.String(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get contract code").Error())
	}

	if len(resCode.Code) == 0 {
		return nil, status.Error(codes.NotFound, errors.New("not a contract").Error())
	}

	blockNumber, err := m.evmJsonRpcBackend.BlockNumber()
	if err != nil {
		return nil, err
	}

	chainId, err := m.evmJsonRpcBackend.ChainID()
	if err != nil {
		return nil, err
	}

	call := func(input string) ([]byte, error) {
		return m.EvmCall(input, contractAddress, chainId, &blockNumber, 0)
	}

	res := make(berpctypes.GenericBackendResponse)

	symbol, err := call("0x95d89b41") // symbol()
	if err != nil {
		return nil, err
	}
	if len(symbol) > 0 {
		unpackedSymbol, err := iberpcutils.UnpackAbiString(symbol, "symbol")
		if err != nil {
			return nil, err
		}
		res["symbol"] = unpackedSymbol
	}

	decimals, err := call("0x313ce567") // decimals()
	if err != nil {
		return nil, err
	}
	res["decimals"] = new(big.Int).SetBytes(decimals).Int64()

	name, err := call("0x06fdde03") // name()
	if err == nil {
		if len(name) > 0 {
			unpackedName, err := iberpcutils.UnpackAbiString(name, "name")
			if err != nil {
				return nil, err
			}
			res["name"] = unpackedName
		}
	}

	return res, nil
}

func (m *EvmBackend) GetErc20Balance(accountAddress common.Address, contractAddresses []common.Address) (berpctypes.GenericBackendResponse, error) {
	res := berpctypes.GenericBackendResponse{
		"account": accountAddress.String(),
	}

	blockNumber, err := m.evmJsonRpcBackend.BlockNumber()
	if err != nil {
		return nil, err
	}

	chainId, err := m.evmJsonRpcBackend.ChainID()
	if err != nil {
		return nil, err
	}

	resForContracts := make([]berpctypes.GenericBackendResponse, 0)

	for _, contractAddress := range contractAddresses {
		var display string
		var decimals int64
		var balance *big.Int

		resCode, err := m.queryClient.EvmQueryClient.Code(m.ctx, &evmtypes.QueryCodeRequest{
			Address: contractAddress.String(),
		})
		if err != nil {
			return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get contract code").Error())
		}

		if len(resCode.Code) == 0 {
			display = ""
			decimals = 0
			balance = big.NewInt(0)
		} else {
			call := func(input string) ([]byte, error) {
				return m.EvmCall(input, contractAddress, chainId, &blockNumber, 0)
			}

			displayBz, err := call("0x95d89b41") // symbol()
			if err != nil {
				// retry with name
				displayBz, err = call("0x06fdde03") // name()
			}
			if err != nil {
				return nil, status.Error(codes.Internal, errors.Wrapf(err, "failed to get name or symbol of %s", contractAddress).Error())
			}
			unpackedDisplay, err := iberpcutils.UnpackAbiString(displayBz, "symbol")
			if err != nil {
				display = fmt.Sprintf("(%s)", contractAddress.String()) // force value
			} else {
				display = unpackedDisplay
			}

			decimalsBz, err := call("0x313ce567") // decimals()
			if err != nil {
				return nil, status.Error(codes.Internal, errors.Wrapf(err, "failed to get decimals of %s", contractAddress).Error())
			}
			decimals = new(big.Int).SetBytes(decimalsBz).Int64()

			balanceBz, err := call("0x70a08231" /*balanceOf(address)*/ + hexutil.Encode(common.LeftPadBytes(accountAddress.Bytes(), 32))[2:])
			if err != nil {
				return nil, status.Error(codes.Internal, errors.Wrapf(err, "failed to get balance of %s on contract %s", accountAddress, contractAddress).Error())
			}
			balance = new(big.Int).SetBytes(balanceBz)
		}

		resForContracts = append(resForContracts, berpctypes.GenericBackendResponse{
			"contract": contractAddress.String(),
			"display":  display,
			"decimals": decimals,
			"balance":  balance.String(),
		})
	}

	res["erc20Balances"] = resForContracts

	return res, nil
}

func (m *EvmBackend) EvmCall(input string, contract common.Address, optionalChainId *hexutil.Big, optionalBlockNumber *hexutil.Uint64, optionalGas uint64) ([]byte, error) {
	var chainId *hexutil.Big
	if optionalChainId == nil {
		resChainId, err := m.evmJsonRpcBackend.ChainID()
		if err != nil {
			return nil, err
		}
		chainId = resChainId
	} else {
		chainId = optionalChainId
	}

	var blockNumber hexutil.Uint64
	if optionalBlockNumber == nil || (*optionalBlockNumber) == 0 {
		resBlockNumber, err := m.evmJsonRpcBackend.BlockNumber()
		if err != nil {
			return nil, err
		}
		blockNumber = resBlockNumber
	} else {
		blockNumber = *optionalBlockNumber
	}

	if strings.HasPrefix(input, "0x") {
		input = input[2:]
	}
	bz, err := hex.DecodeString(input)
	if err != nil {
		return nil, err
	}
	var gas uint64
	if optionalGas == 0 {
		gas = 300_000
	} else {
		gas = optionalGas
	}
	gasB := hexutil.Uint64(gas)
	inputB := hexutil.Bytes(bz)
	gasPriceB := hexutil.Big(*(new(big.Int).SetUint64(math.MaxUint64)))
	nonceB := hexutil.Uint64(0)
	res, err := m.evmJsonRpcBackend.DoCall(evmtypes.TransactionArgs{
		To:       &contract,
		Gas:      &gasB,
		GasPrice: &gasPriceB,
		Nonce:    &nonceB,
		Input:    &inputB,
		ChainID:  chainId,
	}, evmosrpctypes.BlockNumber(blockNumber))
	if err != nil {
		return nil, err
	}
	return res.Ret, nil
}
