//go:build be_json_rpc_evm

package evm

import (
	"encoding/hex"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v12/rpc/backend"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (m *EvmBackend) GetEvmTransactionByHash(hash common.Hash) (berpctypes.GenericBackendResponse, error) {
	rpcTx, err := m.evmJsonRpcBackend.GetTransactionByHash(hash)
	if err != nil {
		return nil, err
	}

	if rpcTx == nil {
		return nil, status.Error(codes.NotFound, "transaction not found")
	}

	evmTxResult, err := m.evmJsonRpcBackend.GetTxByEthHash(hash)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get transaction result")
	}

	blockNumber := rpcTx.BlockNumber.ToInt().Int64()

	block, err := m.clientCtx.Client.Block(m.ctx, &blockNumber)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get block").Error())
	}

	cosmosTx := block.Block.Txs[evmTxResult.TxIndex]
	cosmosTxResult, err := m.queryClient.ServiceClient.GetTx(m.ctx, &tx.GetTxRequest{
		Hash: strings.ToUpper(hex.EncodeToString(cosmosTx.Hash())),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get cosmos tx result").Error())
	}

	txRes := cosmosTxResult.TxResponse
	txEvents := berpctypes.ConvertTxEvent(txRes.Events).RemoveUnnecessaryEvmTxEvents()

	logs, err := backend.TxLogsFromEvents(txRes.Events, int(evmTxResult.MsgIndex))
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get transaction logs").Error())
	}

	receipt, err := m.evmJsonRpcBackend.GetTransactionReceipt(hash)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get transaction receipt").Error())
	}

	res := berpctypes.GenericBackendResponse{
		"hash":        hash,
		"height":      blockNumber,
		"evm_tx":      rpcTx,
		"evm_receipt": receipt,
		"result": map[string]any{
			"code":   txRes.Code,
			"events": txEvents,
			"gas": berpctypes.GenericBackendResponse{
				"limit": txRes.GasWanted,
				"used":  txRes.GasUsed,
			},
		},
	}

	if len(logs) > 0 {
		res["logs"] = logs
	}

	return res, nil
}
