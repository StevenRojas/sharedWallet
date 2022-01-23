package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type TransactionInfo struct {
	Operation string   `json:"operation"`
	Gas       uint64   `json:"gas"`
	GasPrice  *big.Int `json:"gas_price"`
	Cost      *big.Int `json:"cost"`
}

// processTransaction process the transaction in order to get stats
func processTransaction(ctx context.Context, tx *types.Transaction, operation string) {
	if tx == nil {
		return
	}
	_ = TransactionInfo{
		Operation: operation,
		Gas:       tx.Gas(),
		GasPrice:  tx.GasPrice(),
		Cost:      tx.Cost(),
	}
	// TODO: process the whole transaction and use if for stats
}
