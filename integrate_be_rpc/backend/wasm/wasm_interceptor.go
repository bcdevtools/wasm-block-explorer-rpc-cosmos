package wasm

import (
	berpcbackend "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/backend"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ berpcbackend.RequestInterceptor = (*DefaultRequestInterceptor)(nil)

type DefaultRequestInterceptor struct {
	beRpcBackend berpcbackend.BackendI
	backend      WasmBackendI
	bech32Cfg    berpctypes.Bech32Config
}

func NewDefaultRequestInterceptor(
	beRpcBackend berpcbackend.BackendI,
	backend WasmBackendI,
) *DefaultRequestInterceptor {
	return &DefaultRequestInterceptor{
		beRpcBackend: beRpcBackend,
		backend:      backend,
		bech32Cfg:    berpctypes.NewBech32Config(),
	}
}

func (m *DefaultRequestInterceptor) GetTransactionByHash(hashStr string) (intercepted bool, response berpctypes.GenericBackendResponse, err error) {
	// handle WASM txs, otherwise return false
	intercepted = false
	// TODO BE: implement
	return
}

func (m *DefaultRequestInterceptor) GetDenomsInformation() (intercepted, append bool, denoms map[string]string, err error) {
	intercepted = false
	// TODO BE: implement
	return
}

func (m *DefaultRequestInterceptor) GetModuleParams(moduleName string) (intercepted bool, res berpctypes.GenericBackendResponse, err error) {
	var params any

	switch moduleName {
	case "wasm", "cosmwasm":
		wasmParams, errFetch := m.backend.GetWasmModuleParams()
		if errFetch != nil {
			err = errors.Wrap(errFetch, "failed to get wasm params")
		} else {
			params = *wasmParams
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
	if !m.bech32Cfg.IsAccountAddr(accountAddressStr) {
		// not an account address, ignore
		intercepted = false
		append = false
		return
	}

	accAddr, err := sdk.AccAddressFromBech32(accountAddressStr)
	if err != nil {
		// not an account address, ignore
		intercepted = false
		append = false
		return
	}

	if len(accAddr.Bytes()) != 32 {
		// not a contract address, ignore
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

	response = make(berpctypes.GenericBackendResponse)

	codeId, err := m.backend.GetContractCodeId(accountAddressStr)
	if err != nil {
		err = status.Error(codes.Internal, errors.Wrap(err, "failed to check contract code id").Error())
		return
	}

	if codeId == 0 {
		// not a contract, ignore
		return
	}

	contractInfo := berpctypes.GenericBackendResponse{
		"codeId": codeId,
	}
	response["contract"] = contractInfo

	cw20TokenInfo, err := m.backend.GetCw20ContractInfo(accountAddressStr)
	if err == nil && len(cw20TokenInfo) > 0 {
		for k, v := range cw20TokenInfo {
			contractInfo[k] = v
		}
	}

	return
}
