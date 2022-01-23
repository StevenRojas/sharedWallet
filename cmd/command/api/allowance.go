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

var ErrInvalidAllowanceAction = errors.New("invalid allowance action")
var ErrInvalidAmountAction = errors.New("amount should be a positive value")

var allowanceActions = map[string]struct{}{
	"set": {},
	"get": {},
	"increase": {},
	"reduce": {},
}

// NewAllowanceCommand creates the allowance command
func NewAllowanceCommand(ctx context.Context) *cobra.Command {
	var (
		action string
		targetAddress string
		amount int64
	)
	allowanceCommand := &cobra.Command{
		Use:   "allowance",
		Short: "Change the allowance for a beneficiary",
		Run: func(cmd *cobra.Command, args []string) {
			err := runAllowance(ctx, action, targetAddress, amount)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
			}
		},
	}

	allowanceCommand.Flags().StringVar(&action, "action", "", "Action to perform: set, get, increase or reduce")
	allowanceCommand.Flags().Int64Var(&amount, "amount", 0, "Amount")
	allowanceCommand.Flags().StringVarP(&targetAddress, "target.address", "t", "", "Target address")
	_ = allowanceCommand.MarkFlagRequired("action")
	_ = allowanceCommand.MarkFlagRequired("target.address")
	return allowanceCommand
}

func runAllowance(ctx context.Context, action string, targetAddress string, amount int64) error {
	if _, ok := allowanceActions[action]; !ok {
		return ErrInvalidAllowanceAction
	}
	ctxCall, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()
	client, err := ethclient.DialContext(ctxCall, config.App.Blockchain.WS)
	if err != nil {
		return err
	}
	runner := blockchain.NewAllowanceRunner(config.App.Blockchain.PrivateKey, config.App.Contract.Address)

	switch action {
	case blockchain.GetAction:
		allowance, err := runner.GetAllowance(ctx, client, targetAddress)
		if err != nil {
			return err
		}
		fmt.Printf("Current allowance for address %s is %d\n", targetAddress, allowance)
	default:
		if amount <= 0 {
			return ErrInvalidAmountAction
		}
		err = runner.ChangeAllowance(ctx, client, action, targetAddress, amount)
		if err != nil {
			return err
		}
	}

	return nil
}
