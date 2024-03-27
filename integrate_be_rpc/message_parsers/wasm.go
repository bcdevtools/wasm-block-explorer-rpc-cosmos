//go:build be_json_rpc_wasm

package message_parsers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	berpc "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	berpcutils "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
)

func RegisterMessageParsersForWasm() {
	berpc.RegisterMessageParser(&wasmtypes.MsgStoreCode{}, ParseMsgStoreCode)
	berpc.RegisterMessageParser(&wasmtypes.MsgInstantiateContract{}, ParseMsgInstantiateContract)
	berpc.RegisterMessageParser(&wasmtypes.MsgInstantiateContract2{}, ParseMsgInstantiateContract2)
	berpc.RegisterMessageParser(&wasmtypes.MsgClearAdmin{}, ParseMsgClearAdmin)
	berpc.RegisterMessageParser(&wasmtypes.MsgExecuteContract{}, ParseMsgExecuteContract)
	berpc.RegisterMessageParser(&wasmtypes.MsgIBCCloseChannel{}, ParseMsgIBCCloseChannel)
	berpc.RegisterMessageParser(&wasmtypes.MsgIBCSend{}, ParseMsgIBCSend)
	berpc.RegisterMessageParser(&wasmtypes.MsgMigrateContract{}, ParseMsgMigrateContract)
	berpc.RegisterMessageParser(&wasmtypes.MsgUpdateAdmin{}, ParseMsgUpdateAdmin)
	berpc.RegisterMessageParser(&wasmtypes.MsgUpdateInstantiateConfig{}, ParseMsgUpdateInstantiateConfig)
}

func ParseMsgStoreCode(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgStoreCode)

	res = berpctypes.GenericBackendResponse{
		"sender": msg.Sender,
	}
	if msg.InstantiatePermission != nil {
		res["instantiatePermission"] = map[string]any{
			"permission": msg.InstantiatePermission.Permission.String(),
			"addresses":  msg.InstantiatePermission.Addresses,
		}
	}

	rb := berpctypes.NewFriendlyResponseContentBuilder().
		WriteAddress(msg.Sender).
		WriteText(" has stored new contract bytecode")

	var codeId, checksum string

	for _, event := range txResponse.Events {
		if event.Type != wasmtypes.EventTypeStoreCode {
			continue
		}

		for _, attr := range event.Attributes {
			if string(attr.Key) == wasmtypes.AttributeKeyCodeID {
				codeId = string(attr.Value)
			} else if string(attr.Key) == wasmtypes.AttributeKeyChecksum {
				checksum = string(attr.Value)
			}
		}

		break
	}

	if len(codeId) > 0 {
		res["codeId"] = codeId
		rb.WriteText(", code-id = ").WriteText(codeId)
	}

	if len(checksum) > 0 {
		res["checksum"] = checksum
		rb.WriteText(", checksum = ").WriteText(checksum)
	}

	rb.WriteText(" into chain").BuildIntoResponse(res)

	return
}

func ParseMsgInstantiateContract(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgInstantiateContract)

	res = berpctypes.GenericBackendResponse{
		"sender": msg.Sender,
		"codeId": msg.CodeID,
		"msg":    msg.Msg,
	}
	if msg.Admin != "" {
		res["admin"] = msg.Admin
	}
	if len(msg.Msg) > 0 {
		var unmarshalledMsg map[string]any
		err = json.Unmarshal(msg.Msg, &unmarshalledMsg)
		if err == nil {
			res["ctorMsg"] = unmarshalledMsg
		}
	}

	rb := berpctypes.NewFriendlyResponseContentBuilder().
		WriteAddress(msg.Sender).
		WriteText(" has deployed new contract")

	var contractAddress string
	for _, event := range txResponse.Events {
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
		res["contractAddress"] = contractAddress
		rb.WriteText(" ").WriteAddress(contractAddress)
	}

	rb.WriteText(" with code-id ").
		WriteText(fmt.Sprintf("%d", msg.CodeID)).
		BuildIntoResponse(res)

	return
}

func ParseMsgInstantiateContract2(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgInstantiateContract2)

	res = berpctypes.GenericBackendResponse{
		"sender": msg.Sender,
		"codeId": msg.CodeID,
		"msg":    msg.Msg,
	}
	if msg.Admin != "" {
		res["admin"] = msg.Admin
	}
	if len(msg.Msg) > 0 {
		var unmarshalledMsg map[string]any
		err = json.Unmarshal(msg.Msg, &unmarshalledMsg)
		if err == nil {
			res["ctorMsg"] = unmarshalledMsg
		}
	}

	rb := berpctypes.NewFriendlyResponseContentBuilder().
		WriteAddress(msg.Sender).
		WriteText(" has deployed new contract")

	var contractAddress string
	for _, event := range txResponse.Events {
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
		res["contractAddress"] = contractAddress
		rb.WriteText(" ").WriteAddress(contractAddress)
	}

	rb.WriteText(" with code-id ").
		WriteText(fmt.Sprintf("%d", msg.CodeID)).
		BuildIntoResponse(res)

	return
}

func ParseMsgClearAdmin(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgClearAdmin)

	res = berpctypes.GenericBackendResponse{
		"sender":   msg.Sender,
		"contract": msg.Contract,
	}

	berpctypes.NewFriendlyResponseContentBuilder().
		WriteAddress(msg.Sender).
		WriteText(" cleared admin of contract").
		WriteAddress(msg.Contract).
		BuildIntoResponse(res)

	return
}

func ParseMsgExecuteContract(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgExecuteContract)

	res = berpctypes.GenericBackendResponse{
		"sender":   msg.Sender,
		"contract": msg.Contract,
		"funds":    berpcutils.CoinsToMap(msg.Funds...),
	}

	if len(msg.Msg) > 0 {
		var unmarshalledMsg map[string]any
		err = json.Unmarshal(msg.Msg, &unmarshalledMsg)
		if err == nil {
			res["inputMsg"] = unmarshalledMsg
		}
	}

	// TODO BE: implement error message if any

	transfers := make([]berpctypes.GenericBackendResponse, 0)
	for _, event := range txResponse.Events {
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

		transfers = append(transfers, berpctypes.GenericBackendResponse{
			"from":   kv["from"],
			"to":     kv["to"],
			"amount": kv["amount"],
		})
	}

	action := make(berpctypes.GenericBackendResponse)
	if len(transfers) > 0 {
		action["transfers"] = transfers
	}

	res["action"] = action

	berpctypes.NewFriendlyResponseContentBuilder().
		WriteAddress(msg.Sender).
		WriteText(" executes contract ").
		WriteAddress(msg.Contract).
		BuildIntoResponse(res)

	return
}

func ParseMsgIBCCloseChannel(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgIBCCloseChannel)

	res = berpctypes.GenericBackendResponse{
		"channel": msg.Channel,
	}

	berpctypes.NewFriendlyResponseContentBuilder().
		WriteText("Wasm close IBC channel ").
		WriteText(msg.Channel).
		BuildIntoResponse(res)

	return
}

func ParseMsgIBCSend(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgIBCSend)

	res = berpctypes.GenericBackendResponse{
		"channel":               msg.Channel,
		"data":                  hex.EncodeToString(msg.Data),
		"timeoutHeight":         msg.TimeoutHeight,
		"timeoutTimestampNanos": msg.TimeoutTimestamp,
	}

	rb := berpctypes.NewFriendlyResponseContentBuilder().
		WriteText("Wasm IBC send via channel ").
		WriteText(msg.Channel)

	if msg.TimeoutHeight > 0 {
		rb.WriteText(" with timeout-block-height ").WriteText(fmt.Sprintf("%d", msg.TimeoutHeight))
	} else {
		rb.WriteText(" without timeout-block-height")
	}

	if msg.TimeoutTimestamp > 0 {
		rb.WriteText(" with timeout-timestamp ").WriteText(fmt.Sprintf("%d seconds", msg.TimeoutTimestamp/1_000_000_000))
	} else {
		rb.WriteText(" without timeout-timestamp")
	}

	rb.BuildIntoResponse(res)

	return
}

func ParseMsgMigrateContract(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgMigrateContract)

	res = berpctypes.GenericBackendResponse{
		"sender":   msg.Sender,
		"contract": msg.Contract,
		"codeId":   msg.CodeID,
	}

	if len(msg.Msg) > 0 {
		var unmarshalledMsg map[string]any
		err = json.Unmarshal(msg.Msg, &unmarshalledMsg)
		if err == nil {
			res["migrationMsg"] = unmarshalledMsg
		}
	}

	// TODO BE: implement error message if any

	berpctypes.NewFriendlyResponseContentBuilder().
		WriteAddress(msg.Sender).
		WriteText(" migrates contract").
		WriteAddress(msg.Contract).
		WriteText(" to new code-id ").
		WriteText(fmt.Sprintf("%d", msg.CodeID)).
		BuildIntoResponse(res)

	return
}

func ParseMsgUpdateAdmin(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgUpdateAdmin)

	res = berpctypes.GenericBackendResponse{
		"sender":   msg.Sender,
		"contract": msg.Contract,
		"newAdmin": msg.NewAdmin,
	}

	berpctypes.NewFriendlyResponseContentBuilder().
		WriteAddress(msg.Sender).
		WriteText(" updated admin for contract ").
		WriteAddress(msg.Contract).
		WriteText(" to ").
		WriteAddress(msg.NewAdmin).
		BuildIntoResponse(res)

	return
}

func ParseMsgUpdateInstantiateConfig(sdkMsg sdk.Msg, msgIdx uint, tx *tx.Tx, txResponse *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
	msg := sdkMsg.(*wasmtypes.MsgUpdateInstantiateConfig)

	res = berpctypes.GenericBackendResponse{
		"sender": msg.Sender,
		"codeId": msg.CodeID,
	}

	if msg.NewInstantiatePermission != nil {
		res["instantiatePermission"] = map[string]any{
			"permission": msg.NewInstantiatePermission.Permission.String(),
			"addresses":  msg.NewInstantiatePermission.Addresses,
		}
	}

	berpctypes.NewFriendlyResponseContentBuilder().
		WriteAddress(msg.Sender).
		WriteText(" updated init config for code-id ").
		WriteText(fmt.Sprintf("%d", msg.CodeID)).
		BuildIntoResponse(res)

	return
}
