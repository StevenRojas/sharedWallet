package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	AddressBalance = "address"
	ContractBalance = "contract"
)

// Balance interface
type Balance interface {
	GetContractBalance(ctx context.Context, client *ethclient.Client) (int64, error)
	GetAddressBalance(ctx context.Context, client *ethclient.Client, address string) (int64, error)
}

type balance struct {
	privateKey string
	contractAddress string
}

// NewBalanceRunner returns a new runner instance
func NewBalanceRunner(privateKey string, contractAddress string) Balance {
	return &balance{
		privateKey: privateKey,
		contractAddress: contractAddress,
	}
}

// GetContractBalance returns the contract balance
func (b *balance) GetContractBalance(ctx context.Context, client *ethclient.Client) (int64, error) {
	value, err := client.BalanceAt(ctx, common.HexToAddress(b.contractAddress), nil)
	if err != nil {
		return 0, err
	}
	return weiToEther(value).Int64(), nil
}

// GetAddressBalance returns the balance of a given address
func (b *balance) GetAddressBalance(ctx context.Context, client *ethclient.Client, address string) (int64, error) {
	value, err := client.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return 0, err
	}
	return weiToEther(value).Int64(), nil
}