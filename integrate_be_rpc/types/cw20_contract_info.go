package types

import "math/big"

type Cw20TokenInfo struct {
	Name           string   `json:"name,omitempty"`
	Symbol         string   `json:"symbol,omitempty"`
	Decimals       uint8    `json:"decimals"`
	TotalSupplyStr string   `json:"total_supply,omitempty"`
	TotalSupply    *big.Int `json:"-"`
}
