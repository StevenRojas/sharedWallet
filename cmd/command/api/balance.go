package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/StevenRojas/sharedWallet/config"
	"github.com/StevenRojas/sharedWallet/pkg/blockchain"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var (
	ErrInvalidBalanceAction = errors.New("invalid balance action")
	ErrInvalidBalanceAddress = errors.New("invalid balance address")
)

var balanceOfList = map[string]struct{}{
	"owner": {},
	"address": {},
	"contract": {},
}

// NewBalanceCommand creates the balance command
func NewBalanceCommand(ctx context.Context) *cobra.Command {
	var (
		of string
		targetAddress string
	)
	balanceCommand := &cobra.Command{
		Use:   "balance",
		Short: "Get the balance of an address or contract",
		Run: func(cmd *cobra.Command, args []string) {
			err := runBalance(ctx, of, targetAddress)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
			}
		},
	}

	balanceCommand.Flags().StringVar(&of, "of", "", "Balance of: address, contract")
	balanceCommand.Flags().StringVarP(&targetAddress, "target.address", "t", "", "Target address")
	_ = balanceCommand.MarkFlagRequired("of")
	return balanceCommand
}

func runBalance(ctx context.Context, of string, targetAddress string) error {
	if _, ok := balanceOfList[of]; !ok {
		return ErrInvalidBalanceAction
	}

	ctxCall, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()
	client, err := ethclient.DialContext(ctxCall, config.App.Blockchain.WS)
	if err != nil {
		return err
	}
	runner := blockchain.NewBalanceRunner(config.App.Blockchain.PrivateKey, config.App.Contract.Address)

	switch of {
	case blockchain.ContractBalance:
		balance, err := runner.GetContractBalance(ctx, client)
		if err != nil {
			return err
		}
		fmt.Printf("The contract balance is %d\n", balance)
	case blockchain.AddressBalance:
		if targetAddress == "" {
			return ErrInvalidBalanceAddress
		}
		balance, err := runner.GetAddressBalance(ctx, client, targetAddress)
		if err != nil {
			return err
		}
		fmt.Printf("The contract balance is %d\n", balance)
	}


	return nil
}

