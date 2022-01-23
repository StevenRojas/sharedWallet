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

var ErrInvalidTransferAction = errors.New("invalid transfer action")

var transferActions = map[string]struct{}{
	"send": {},
	"receive": {},
}

// NewTransferCommand creates the transfers command
func NewTransferCommand(ctx context.Context) *cobra.Command {
	var (
		action string
		targetAddress string
		amount int64
	)
	transfersCommand := &cobra.Command{
		Use:   "transfer",
		Short: "Perform transfer operations",
		Run: func(cmd *cobra.Command, args []string) {
			err := runTransfers(ctx, action, targetAddress, amount)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
			}
		},
	}

	transfersCommand.Flags().StringVar(&action, "action", "", "Action to perform: send, receive")
	transfersCommand.Flags().Int64Var(&amount, "amount", 0, "Amount")
	transfersCommand.Flags().StringVarP(&targetAddress, "target.address", "t", "", "Target address")
	_ = transfersCommand.MarkFlagRequired("action")
	_ = transfersCommand.MarkFlagRequired("amount")
	return transfersCommand
}

func runTransfers(ctx context.Context, action string, targetAddress string, amount int64) error {
	if _, ok := transferActions[action]; !ok {
		return ErrInvalidTransferAction
	}
	ctxCall, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()
	client, err := ethclient.DialContext(ctxCall, config.App.Blockchain.WS)
	if err != nil {
		return err
	}
	runner := blockchain.NewTransfersRunner(config.App.Blockchain.PrivateKey, config.App.Contract.Address)

	if amount <= 0 {
		return ErrInvalidAmountAction
	}

	switch action {
	case blockchain.SendAction:
		if targetAddress == "" {
			return ErrInvalidOwnershipAddress
		}
		err = runner.Send(ctx, client, targetAddress, amount)
	case blockchain.ReceiveAction:
		err = runner.Receive(ctx, client, amount)
	}
	if err != nil {
		return err
	}

	return nil
}

