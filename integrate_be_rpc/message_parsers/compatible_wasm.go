//go:build !be_json_rpc_wasm

package message_parsers

import (
	"github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/integrate_be_rpc/compatible"
)

func RegisterMessageParsersForWasm() {
	compatible.PanicInvalidBuildTag()
}
