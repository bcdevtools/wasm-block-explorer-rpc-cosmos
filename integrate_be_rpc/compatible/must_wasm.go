//go:build !be_json_rpc_wasm

package compatible

func PanicInvalidBuildTag() {
	panic("invalid build tag, require `be_json_rpc_wasm`")
}
