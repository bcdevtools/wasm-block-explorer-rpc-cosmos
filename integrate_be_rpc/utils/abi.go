package utils

import "github.com/evmos/evmos/v12/contracts"

func UnpackAbiString(bz []byte, optionalMethodName string) (string, error) {
	methodName := "symbol"
	if len(optionalMethodName) > 0 {
		methodName = optionalMethodName
	}
	unpacked, err := contracts.ERC20MinterBurnerDecimalsContract.ABI.Methods[methodName].Outputs.Unpack(bz)
	if err != nil {
		return "", err
	}
	return unpacked[0].(string), nil
}
