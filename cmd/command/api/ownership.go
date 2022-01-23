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
	ErrInvalidOwnershipAction = errors.New("invalid ownership action")
	ErrInvalidOwnershipAddress = errors.New("invalid ownership address")
)

var ownershipActions = map[string]struct{}{
	"get": {},
	"transfer": {},
}

// NewOwnershipCommand creates the ownership command
func NewOwnershipCommand(ctx context.Context) *cobra.Command {
	var (
		action string
		targetAddress string
	)
	ownershipCommand := &cobra.Command{
		Use:   "ownership",
		Short: "Get or transfer contract ownership",
		Run: func(cmd *cobra.Command, args []string) {
			err := runOwnership(ctx, action, targetAddress)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
			}
		},
	}

	ownershipCommand.Flags().StringVar(&action, "action", "", "Ownership action: get, transfer")
	ownershipCommand.Flags().StringVarP(&targetAddress, "target.address", "t", "", "Target address")
	_ = ownershipCommand.MarkFlagRequired("action")
	return ownershipCommand
}

func runOwnership(ctx context.Context, action string, targetAddress string) error {
	if _, ok := ownershipActions[action]; !ok {
		return ErrInvalidOwnershipAction
	}

	ctxCall, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()
	client, err := ethclient.DialContext(ctxCall, config.App.Blockchain.WS)
	if err != nil {
		return err
	}
	runner := blockchain.NewOwnerRunner(config.App.Blockchain.PrivateKey, config.App.Contract.Address)

	switch action {
	case blockchain.GetAction:
		owner, err := runner.GetOwner(ctx, client)
		if err != nil {
			return err
		}
		fmt.Printf("Current owner address is %s\n", owner)
	default:
		if targetAddress == "" {
			return ErrInvalidOwnershipAddress
		}
		err = runner.TransferOwner(ctx, client, targetAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

