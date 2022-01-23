package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

const (
	SendAction    = "send"
	ReceiveAction = "receive"
)

// Transfers interface
type Transfers interface {
	Receive(ctx context.Context, client *ethclient.Client, amount int64) error
	Send(ctx context.Context, client *ethclient.Client, target string, amount int64) error
}

type transfers struct {
	privateKey string
	contractAddress string
}

// NewTransfersRunner returns a new runner instance
func NewTransfersRunner(privateKey string, contractAddress string) Transfers {
	return &transfers{
		privateKey: privateKey,
		contractAddress: contractAddress,
	}
}

// Receive method to receive founds in the contract
func (t *transfers) Receive(ctx context.Context, client *ethclient.Client, amount int64) error {
	contract, err := getContract(ctx, client, t.contractAddress)
	if err != nil {
		return err
	}

	signer, err := getSigner(ctx, client)
	if err != nil {
		return err
	}

	signer.Value = etherToWei(big.NewInt(amount))
	tx, txErr := contract.Receive(signer)
	if txErr != nil {
		return txErr
	}
	receipt, err := bind.WaitMined(ctx, client, tx)
	if receipt.Status != types.ReceiptStatusSuccessful || err != nil {
		return err
	}
	processTransaction(ctx, tx, "receive")

	return nil
}

// Send method to send founds to a beneficiary
func (t *transfers) Send(ctx context.Context, client *ethclient.Client, target string, amount int64) error {
	contract, err := getContract(ctx, client, t.contractAddress)
	if err != nil {
		return err
	}

	signer, err := getSigner(ctx, client)
	if err != nil {
		return err
	}

	targetAddress := common.HexToAddress(target)
	tx, txErr := contract.SendMoney(signer, targetAddress, etherToWei(big.NewInt(amount)))
	if txErr != nil {
		return txErr
	}
	receipt, err := bind.WaitMined(ctx, client, tx)
	if receipt.Status != types.ReceiptStatusSuccessful || err != nil {
		return err
	}
	processTransaction(ctx, tx, "send")

	return nil
}