//go:build be_json_rpc_evm

package backend

import (
	berpcbackend "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/backend"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	iberpcbackend "github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/backend/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	intercepted = false
	append = true
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
