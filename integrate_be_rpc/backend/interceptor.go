package backend

import (
	berpcbackend "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/backend"
)

var _ berpcbackend.RequestInterceptor = (*DefaultRequestInterceptor)(nil)
