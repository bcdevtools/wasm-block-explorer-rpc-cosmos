//go:build be_json_rpc_evm

package evm

import (
	"encoding/hex"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/evmos/evmos/v12/contracts"
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

	call := func(_4bytes string) ([]byte, error) {
		return m.evmCall(_4bytes, contractAddress, chainId, blockNumber, 0)
	}

	res := make(berpctypes.GenericBackendResponse)

	symbol, err := call("0x95d89b41")
	if err != nil {
		return nil, err
	}
	unpackedSymbol, err := contracts.ERC20MinterBurnerDecimalsContract.ABI.Methods["symbol"].Outputs.Unpack(symbol)
	if err != nil {
		return nil, err
	}
	res["symbol"] = unpackedSymbol[0].(string)

	decimals, err := call("0x313ce567")
	if err != nil {
		return nil, err
	}
	res["decimals"] = new(big.Int).SetBytes(decimals).Int64()

	name, err := call("0x06fdde03")
	if err == nil {
		unpackedName, err := contracts.ERC20MinterBurnerDecimalsContract.ABI.Methods["name"].Outputs.Unpack(name)
		if err != nil {
			return nil, err
		}
		res["name"] = unpackedName[0].(string)
	}

	return res, nil
}

func (m *EvmBackend) GetErc20Balance(accountAddress common.Address, contractAddresses []common.Address) (berpctypes.GenericBackendResponse, error) {
	res := berpctypes.GenericBackendResponse{
		"account": accountAddress.String(),
	}

	for _, contractAddress := range contractAddresses {
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

		call := func(_4bytes string) ([]byte, error) {
			return m.evmCall(_4bytes, contractAddress, chainId, blockNumber, 0)
		}

		resPerContract := berpctypes.GenericBackendResponse{
			"contract": contractAddress.String(),
		}

		display, err := call("0x95d89b41") // symbol()
		if err != nil {
			// retry with name
			display, err = call("0x06fdde03") // name()
		}
		if err == nil && len(display) > 0 {
			unpackedDisplay, err := contracts.ERC20MinterBurnerDecimalsContract.ABI.Methods["symbol"].Outputs.Unpack(display)
			if err != nil {
				resPerContract["display_error"] = err.Error()
			} else {
				resPerContract["display"] = unpackedDisplay[0].(string)
			}
		} else if err == nil {
			resPerContract["display"] = ""
		} else {
			resPerContract["display_error"] = err.Error()
		}

		decimals, err := call("0x313ce567")
		if err != nil {
			resPerContract["decimals_error"] = err.Error()
		} else {
			resPerContract["decimals"] = new(big.Int).SetBytes(decimals).Int64()
		}

		balance, err := call("0x70a08231" + hexutil.Encode(common.LeftPadBytes(accountAddress.Bytes(), 32))[2:])
		if err != nil {
			return nil, err
		}
		resPerContract["balance"] = new(big.Int).SetBytes(balance).String()

		res[contractAddress.String()] = resPerContract
	}

	return res, nil
}

func (m *EvmBackend) evmCall(_4bytes string, contract common.Address, chainId *hexutil.Big, blockNumber hexutil.Uint64, overrideGas uint64) ([]byte, error) {
	if strings.HasPrefix(_4bytes, "0x") {
		_4bytes = _4bytes[2:]
	}
	bz, err := hex.DecodeString(_4bytes)
	if err != nil {
		return nil, err
	}
	var gas uint64
	if overrideGas == 0 {
		gas = 300_000
	} else {
		gas = overrideGas
	}
	gasB := hexutil.Uint64(gas)
	dataB := hexutil.Bytes(bz)
	gasPriceB := hexutil.Big(*(new(big.Int).SetUint64(math.MaxUint64)))
	nonceB := hexutil.Uint64(0)
	res, err := m.evmJsonRpcBackend.DoCall(evmtypes.TransactionArgs{
		From:                 nil,
		To:                   &contract,
		Gas:                  &gasB,
		GasPrice:             &gasPriceB,
		MaxFeePerGas:         nil,
		MaxPriorityFeePerGas: nil,
		Value:                nil,
		Nonce:                &nonceB,
		Input:                &dataB,
		AccessList:           nil,
		ChainID:              chainId,
	}, evmosrpctypes.BlockNumber(blockNumber))
	if err != nil {
		return nil, err
	}
	return res.Ret, nil
}
