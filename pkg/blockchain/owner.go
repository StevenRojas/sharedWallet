package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Owner interface
type Owner interface {
	GetOwner(ctx context.Context, client *ethclient.Client) (string, error)
	TransferOwner(ctx context.Context, client *ethclient.Client, targetAddress string) error
}

type owner struct {
	privateKey string
	contractAddress string
}

// NewOwnerRunner returns a new runner instance
func NewOwnerRunner(privateKey string, contractAddress string) Owner {
	return &owner{
		privateKey: privateKey,
		contractAddress: contractAddress,
	}
}

// GetOwner returns the contract owner address
func (o *owner) GetOwner(ctx context.Context, client *ethclient.Client) (string, error) {
	contract, err := getContract(ctx, client, o.contractAddress)
	if err != nil {
		return "", err
	}

	ownerAddress, err := contract.Owner(&bind.CallOpts{Pending: false, Context: ctx})
	if err != nil {
		return "", err
	}
	return ownerAddress.Hex(), nil
}

// TransferOwner transfer the ownership to a target address
func (o *owner) TransferOwner(ctx context.Context, client *ethclient.Client, targetAddress string) error {
	contract, err := getContract(ctx, client, o.contractAddress)
	if err != nil {
		return err
	}
	signer, err := getSigner(ctx, client)
	if err != nil {
		return err
	}

	tx, txErr := contract.TransferOwnership(signer, common.HexToAddress(targetAddress))
	if txErr != nil {
		return txErr
	}
	receipt, err := bind.WaitMined(ctx, client, tx)
	if receipt.Status != types.ReceiptStatusSuccessful || err != nil {
		return err
	}
	processTransaction(ctx, tx, "transfer owner")
	return nil
}