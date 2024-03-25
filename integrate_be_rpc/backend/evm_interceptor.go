//go:build be_json_rpc_evm

package backend

import (
	berpcbackend "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/backend"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	berpcutils "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/utils"
	iberpcbackend "github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/backend/evm"
	iberpcutils "github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"strings"
)

type DefaultRequestInterceptor struct {
	beRpcBackend berpcbackend.BackendI
	backend      iberpcbackend.EvmBackendI
}

func NewDefaultRequestInterceptor(
	beRpcBackend berpcbackend.BackendI,
	backend iberpcbackend.EvmBackendI,
) *DefaultRequestInterceptor {
	return &DefaultRequestInterceptor{
		beRpcBackend: beRpcBackend,
		backend:      backend,
	}
}

func (m *DefaultRequestInterceptor) GetTransactionByHash(hashStr string) (intercepted bool, response berpctypes.GenericBackendResponse, err error) {
	// handle EVM txs, otherwise return false

	hashStr = strings.ToLower(hashStr)
	if !strings.HasPrefix(hashStr, "0x") {
		intercepted = false
		return
	}

	intercepted = true
	response, err = m.backend.GetEvmTransactionByHash(common.HexToHash(hashStr))
	return
}

func (m *DefaultRequestInterceptor) GetDenomsInformation() (intercepted, append bool, denoms map[string]string, err error) {
	evmParams, errFetchEvmParams := m.backend.GetEvmModuleParams()
	if errFetchEvmParams != nil {
		err = errors.Wrap(errFetchEvmParams, "failed to get evm params")
		return
	}

	intercepted = false // provide information for the account, so we don't need to ignore other response information
	defer func() {
		if err == nil {
			append = true
		}
	}()

	denoms = map[string]string{
		"evm": evmParams.EvmDenom,
	}

	return
}

func (m *DefaultRequestInterceptor) GetModuleParams(moduleName string) (intercepted bool, res berpctypes.GenericBackendResponse, err error) {
	var params any

	switch moduleName {
	case "evm":
		evmParams, errFetch := m.backend.GetEvmModuleParams()
		if errFetch != nil {
			err = errors.Wrap(errFetch, "failed to get evm params")
		} else {
			params = *evmParams
		}
		break
	default:
		intercepted = false
		return
	}

	if err != nil {
		return
	}

	res, err = berpctypes.NewGenericBackendResponseFrom(params)
	if err != nil {
		err = status.Error(codes.Internal, errors.Wrap(err, "module params").Error())
		return
	}

	intercepted = true
	return
}

// GetAccount returns the contract information if the account is a contract. Other-wise no-op.
func (m *DefaultRequestInterceptor) GetAccount(accountAddressStr string) (intercepted, append bool, response berpctypes.GenericBackendResponse, err error) {
	accAddrStr := berpcutils.ConvertToAccAddressIfHexOtherwiseKeepAsIs(accountAddressStr)

	if !strings.HasPrefix(accAddrStr, sdk.GetConfig().GetBech32AccountAddrPrefix()+"1") &&
		!strings.HasPrefix(accAddrStr, "0x") {
		// not an account address, ignore
		intercepted = false
		append = false
		return
	}

	accAddr, err := sdk.AccAddressFromBech32(accAddrStr)
	if err != nil {
		// not an account address, ignore
		intercepted = false
		append = false
		return
	}

	if len(accAddr.Bytes()) != common.AddressLength {
		// not an EVM account address, ignore
		intercepted = false
		append = false
		return
	}

	intercepted = false // provide information for the account, so we don't need to ignore other response information
	defer func() {
		if err == nil {
			append = true
		} else {
			response = nil // eraser
		}
	}()

	address := common.BytesToAddress(accAddr.Bytes())
	response = make(berpctypes.GenericBackendResponse)

	code, err := m.backend.GetContractCode(address)
	if err != nil {
		err = status.Error(codes.Internal, errors.Wrap(err, "failed to check contract code").Error())
		return
	}

	if len(code) == 0 {
		// not a contract, ignore
		return
	}

	call := func(input string) ([]byte, error) {
		return m.backend.EvmCall(input, address, nil, nil, 0)
	}

	contractInfo := make(berpctypes.GenericBackendResponse)
	response["contract"] = contractInfo

	symbol, err := call("0x95d89b41") // symbol()
	if err == nil {
		if len(symbol) > 0 {
			unpackedSymbol, err := iberpcutils.UnpackAbiString(symbol, "symbol")
			if err == nil {
				contractInfo["symbol"] = unpackedSymbol
			}
		}
	}

	decimals, err := call("0x313ce567") // decimals()
	if err == nil {
		contractInfo["decimals"] = new(big.Int).SetBytes(decimals).Int64()
	}

	name, err := call("0x06fdde03") // name()
	if err == nil {
		if len(name) > 0 {
			unpackedName, err := iberpcutils.UnpackAbiString(name, "name")
			if err == nil {
				contractInfo["name"] = unpackedName
			}
		}
	}

	return
}
