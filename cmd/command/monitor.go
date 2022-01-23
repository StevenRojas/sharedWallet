package command

import (
	"context"
	"github.com/StevenRojas/sharedWallet/config"
	"github.com/StevenRojas/sharedWallet/pkg/blockchain"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

// NewMonitorCommand creates the monitor command
func NewMonitorCommand(ctx context.Context) *cobra.Command {
	monitorCommand := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor events in the blockchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return monitoring(ctx)
		},
	}
	monitorCommand.Flags().StringP("contract.address", "c", "", "Contract address")
	return monitorCommand
}

func monitoring(ctx context.Context) error {
	ctxCall, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()
	client, err := ethclient.DialContext(ctxCall, config.App.Blockchain.WS)
	if err != nil {
		return err
	}
	monitor := blockchain.NewMonitor(config.App.Contract.Address)
	err = monitor.Start(ctx, client)
	if err != nil {
		return err
	}
	return nil
}