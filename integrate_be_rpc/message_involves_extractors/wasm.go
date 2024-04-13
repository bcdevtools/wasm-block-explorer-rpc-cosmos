package message_involves_extractors

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	berpc "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	berpcutils "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/utils"
	"github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/wasm"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	tmtypes "github.com/tendermint/tendermint/types"
)

func RegisterMessageInvolvesExtractorsForWasm(wasmBeRpcBackend wasm.WasmBackendI) {
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgStoreCode{}, ExtractFromMsgStoreCode)
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgInstantiateContract{}, func(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
		msg := sdkMsg.(*wasmtypes.MsgInstantiateContract)

		events, err := wasmBeRpcBackend.GetTmTxResult(tmTx)
		if err != nil {
			return nil, err
		}

		res = berpctypes.NewMessageInvolversResult()
		res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Sender)

		var contractAddress string
		for _, event := range events {
			if event.Type != wasmtypes.EventTypeInstantiate {
				continue
			}

			for _, attr := range event.Attributes {
				if string(attr.Key) == wasmtypes.AttributeKeyContractAddr {
					contractAddress = string(attr.Value)
					break
				}
			}

			break
		}

		if len(contractAddress) > 0 {
			res.AddGenericInvolvers(berpctypes.MessageInvolvers, contractAddress)
		}

		return
	})
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgInstantiateContract2{}, func(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
		msg := sdkMsg.(*wasmtypes.MsgInstantiateContract2)

		events, err := wasmBeRpcBackend.GetTmTxResult(tmTx)
		if err != nil {
			return nil, err
		}

		res = berpctypes.NewMessageInvolversResult()
		res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Sender)

		var contractAddress string
		for _, event := range events {
			if event.Type != wasmtypes.EventTypeInstantiate {
				continue
			}

			for _, attr := range event.Attributes {
				if string(attr.Key) == wasmtypes.AttributeKeyContractAddr {
					contractAddress = string(attr.Value)
					break
				}
			}

			break
		}

		if len(contractAddress) > 0 {
			res.AddGenericInvolvers(berpctypes.MessageInvolvers, contractAddress)
		}

		return
	})
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgClearAdmin{}, ExtractFromMsgClearAdmin)
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgExecuteContract{}, func(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
		msg := sdkMsg.(*wasmtypes.MsgExecuteContract)

		res = berpctypes.NewMessageInvolversResult()
		res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Sender)
		res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Contract)

		bech32Cfg := berpctypes.NewBech32Config()

		events, err := wasmBeRpcBackend.GetTmTxResult(tmTx)
		if err != nil {
			return nil, err
		}

		for _, event := range events {
			for _, attribute := range event.Attributes {
				if bech32Cfg.IsAccountAddr(string(attribute.Value)) {
					res.AddGenericInvolvers(berpctypes.MessageInvolvers, string(attribute.Value))
				}
			}
		}

		trackerCw20Contract := make(map[string]bool)
		for _, event := range events {
			match, kv := berpcutils.IsEventTypeWithAllAttributes(
				event,
				wasmtypes.WasmModuleEventType,
				wasmtypes.AttributeKeyContractAddr,
				"action",
				"from",
				"to",
				"amount",
			)
			if !match {
				continue
			}

			if kv["action"] != "transfer" {
				continue
			}

			contractAddr := kv[wasmtypes.AttributeKeyContractAddr]
			isCw20Contract, foundContract := trackerCw20Contract[contractAddr]
			if !foundContract {
				cw20TokenInfo, err := wasmBeRpcBackend.GetCw20TokenInfo(msg.Contract)
				if err == nil && len(cw20TokenInfo.Symbol) > 0 {
					isCw20Contract = true
				}
				trackerCw20Contract[contractAddr] = isCw20Contract
			}

			if isCw20Contract {
				res.AddContractInvolvers(
					berpctypes.Erc20Involvers,
					berpctypes.ContractAddress(contractAddr),
					kv["from"], kv["to"],
				)
			}
		}

		return
	})
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgIBCCloseChannel{}, ExtractFromMsgIBCCloseChannel)
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgIBCSend{}, ExtractFromMsgIBCSend)
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgMigrateContract{}, ExtractFromMsgMigrateContract)
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgUpdateAdmin{}, ExtractFromMsgUpdateAdmin)
	berpc.RegisterMessageInvolversExtractor(&wasmtypes.MsgUpdateInstantiateConfig{}, ExtractFromMsgUpdateInstantiateConfig)
}

func ExtractFromMsgStoreCode(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
	msg := sdkMsg.(*wasmtypes.MsgStoreCode)

	res = berpctypes.NewMessageInvolversResult()

	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Sender)
	if msg.InstantiatePermission != nil && len(msg.InstantiatePermission.Addresses) > 0 {
		res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.InstantiatePermission.Addresses...)
	}

	return
}

func ExtractFromMsgClearAdmin(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
	msg := sdkMsg.(*wasmtypes.MsgClearAdmin)

	res = berpctypes.NewMessageInvolversResult()
	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Sender)
	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Contract)

	// TODO BE: add the removed admin to the response

	return
}

func ExtractFromMsgIBCCloseChannel(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
	// msg := sdkMsg.(*wasmtypes.MsgIBCCloseChannel)

	res = berpctypes.NewMessageInvolversResult()

	return
}

func ExtractFromMsgIBCSend(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
	// msg := sdkMsg.(*wasmtypes.MsgIBCSend)

	res = berpctypes.NewMessageInvolversResult()

	// TODO BE: implement?

	return
}

func ExtractFromMsgMigrateContract(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
	msg := sdkMsg.(*wasmtypes.MsgMigrateContract)

	res = berpctypes.NewMessageInvolversResult()

	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Sender)
	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Contract)

	return
}

func ExtractFromMsgUpdateAdmin(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
	msg := sdkMsg.(*wasmtypes.MsgUpdateAdmin)

	res = berpctypes.NewMessageInvolversResult()

	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Sender)
	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Contract)
	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.NewAdmin)

	return
}

func ExtractFromMsgUpdateInstantiateConfig(sdkMsg sdk.Msg, tx *tx.Tx, tmTx tmtypes.Tx, clientCtx client.Context) (res berpctypes.MessageInvolversResult, err error) {
	msg := sdkMsg.(*wasmtypes.MsgUpdateInstantiateConfig)

	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Sender)

	return
}
