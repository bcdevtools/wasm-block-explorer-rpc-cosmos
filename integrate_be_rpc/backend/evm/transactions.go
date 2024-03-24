//go:build be_json_rpc_evm

package evm

import (
	"encoding/hex"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	berpcutils "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v12/rpc/backend"
	"github.com/pkg/errors"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (m *EvmBackend) GetEvmTransactionInvolversByHash(hash common.Hash, optionalTxResult *coretypes.ResultTx) (berpctypes.MessageInvolversResult, error) {
	evmTxResult, err := m.evmJsonRpcBackend.GetTxByEthHash(hash)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get transaction result")
	}

	blockNumber := evmTxResult.Height

	receipt, err := m.evmJsonRpcBackend.GetTransactionReceipt(hash)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get transaction receipt").Error())
	}

	var txAbciEvents []abcitypes.Event
	if optionalTxResult == nil {
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
		txAbciEvents = cosmosTxResult.TxResponse.Events
	} else {
		txAbciEvents = optionalTxResult.TxResult.Events
	}

	logs, err := backend.TxLogsFromEvents(txAbciEvents, int(evmTxResult.MsgIndex))
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get transaction logs").Error())
	}

	involvers := make(berpctypes.MessageInvolversResult)
	involvers.Add(berpctypes.MessageInvolvers, sdk.AccAddress(receipt["from"].(common.Address).Bytes()).String())
	to := receipt["to"]
	if to != nil {
		involvers.Add(berpctypes.MessageInvolvers, sdk.AccAddress(to.(*common.Address).Bytes()).String())
	} else {
		involvers.Add(berpctypes.MessageInvolvers, sdk.AccAddress(receipt["contractAddress"].(common.Address).Bytes()).String())
	}

	for _, log := range logs {
		involvers.Add(berpctypes.MessageInvolvers, sdk.AccAddress(log.Address.Bytes()).String())

		var involverType berpctypes.InvolversType
		var addrFromTopic1, addrFromTopic2, addrFromTopic3 bool

		involverType = berpctypes.MessageInvolvers // default

		if berpcutils.IsEvmEventMatch(
			log.Topics, log.Data, 3, /*3 topics is ERC-20 transfer*/
			berpctypes.EvmEvent_Erc20_Erc721_Transfer,
			true, true, false,
			true,
		) {
			involverType = berpctypes.Erc20Involvers
			addrFromTopic1 = true
			addrFromTopic2 = true
		} else if berpcutils.IsEvmEventMatch(
			log.Topics, log.Data, 4, /*4 topics is NFT transfer*/
			berpctypes.EvmEvent_Erc20_Erc721_Transfer,
			true, true, false,
			false,
		) {
			involverType = berpctypes.NftInvolvers
			addrFromTopic1 = true
			addrFromTopic2 = true
		} else if berpcutils.IsEvmEventMatch(
			log.Topics, log.Data, 3, /*3 topics is ERC-20 approvals*/
			berpctypes.EvmEvent_Erc20_Erc721_Approval,
			true, true, false,
			true,
		) {
			addrFromTopic1 = true
			addrFromTopic2 = true
		} else if berpcutils.IsEvmEventMatch(
			log.Topics, log.Data, 4, /*4 topics is NFT approvals*/
			berpctypes.EvmEvent_Erc20_Erc721_Approval,
			true, true, false,
			false,
		) {
			addrFromTopic1 = true
			addrFromTopic2 = true
		} else if berpcutils.IsEvmEventMatch(
			log.Topics, log.Data, 3,
			berpctypes.EvmEvent_Erc721_Erc1155_ApprovalForAll,
			true, true, false,
			true,
		) {
			addrFromTopic1 = true
			addrFromTopic2 = true
		} else if berpcutils.IsEvmEventMatch(
			log.Topics, log.Data, 4,
			berpctypes.EvmEvent_Erc1155_TransferSingle,
			true, true, true,
			true,
		) {
			addrFromTopic1 = true
			addrFromTopic2 = true
			addrFromTopic3 = true
		} else if berpcutils.IsEvmEventMatch(
			log.Topics, log.Data, 4,
			berpctypes.EvmEvent_Erc1155_TransferBatch,
			true, true, true,
			true,
		) {
			addrFromTopic1 = true
			addrFromTopic2 = true
			addrFromTopic3 = true
		} else if berpcutils.IsEvmEventMatch(
			log.Topics, log.Data, 2,
			berpctypes.EvmEvent_WDeposit,
			true, false, false,
			true,
		) {
			involverType = berpctypes.Erc20Involvers
			addrFromTopic1 = true
		} else if berpcutils.IsEvmEventMatch(
			log.Topics, log.Data, 2,
			berpctypes.EvmEvent_WWithdraw,
			true, false, false,
			true,
		) {
			involverType = berpctypes.Erc20Involvers
			addrFromTopic1 = true
		} else {
			continue
		}

		if addrFromTopic1 {
			involvers.Add(involverType, berpcutils.AccAddressFromTopic(log.Topics[1]).String())
		}
		if addrFromTopic2 {
			involvers.Add(involverType, berpcutils.AccAddressFromTopic(log.Topics[2]).String())
		}
		if addrFromTopic3 {
			involvers.Add(involverType, berpcutils.AccAddressFromTopic(log.Topics[3]).String())
		}
	}

	return involvers, nil
}

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

	receipt, err := m.evmJsonRpcBackend.GetTransactionReceipt(hash)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get transaction receipt").Error())
	}

	blockNumber := evmTxResult.Height

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
