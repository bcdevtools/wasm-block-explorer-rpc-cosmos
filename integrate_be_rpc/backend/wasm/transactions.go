package wasm

import (
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (m *WasmBackend) GetWasmTransactionInvolversByHash(hash string) (berpctypes.MessageInvolversResult, error) {
	// TODO BE: implement
	return nil, nil
}

func (m *WasmBackend) GetWasmTransactionByHash(hash string) (berpctypes.GenericBackendResponse, error) {
	// TODO BE: implement
	return nil, nil
}

func (m *WasmBackend) GetTmTxResult(tmTx tmtypes.Tx) ([]abci.Event, error) {
	resTxResult, errTxResult := m.clientCtx.Client.Tx(m.ctx, tmTx.Hash(), false)
	if errTxResult != nil {
		return nil, errTxResult
	}

	return resTxResult.TxResult.Events, nil
}
