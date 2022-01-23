package command

import (
	"context"
	"errors"
	"github.com/StevenRojas/sharedWallet/cmd/command/api"
	"github.com/spf13/cobra"
)

// NewRunnerCommand creates the runner command
func NewRunnerCommand(ctx context.Context) *cobra.Command {
	runCommand := &cobra.Command{
		Use:   "run",
		Short: "Run contract methods in the blockchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("please specify a subcommand: [allowance, balance, ownership or transfer]")
		},
	}

	runCommand.AddCommand(api.NewAllowanceCommand(ctx))
	runCommand.AddCommand(api.NewTransferCommand(ctx))
	runCommand.AddCommand(api.NewBalanceCommand(ctx))
	runCommand.AddCommand(api.NewOwnershipCommand(ctx))
	return runCommand
}