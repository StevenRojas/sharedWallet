package blockchain

import (
	"context"
	contracts "github.com/StevenRojas/sharedWallet/contracts/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Deployer interface
type Deployer interface {
	Deploy(ctx context.Context, client *ethclient.Client) error
	ContractAddress() string
}

type deployer struct {
	address common.Address
	transaction *types.Transaction
	contract *contracts.Contract
}

// NewDeployer returns a new runner instance
func NewDeployer() Deployer {
	return &deployer{}
}

// Deploy deploys a new Ethereum contract
func (d *deployer) Deploy(ctx context.Context, client *ethclient.Client) error {
	signer, err := getSigner(ctx, client)
	if err != nil {
		return err
	}
	address, tx, contract, err := contracts.DeployContract(signer, client)
	if err != nil {
		return err
	}
	_, err = bind.WaitDeployed(ctx, client, tx)
	if err != nil {
		return err
	}
	d.address = address
	d.transaction = tx
	d.contract = contract
	return nil
}

// ContractAddress returns the contract address
func (d *deployer) ContractAddress() string {
	return d.address.Hex()
}