package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	contracts "github.com/StevenRojas/sharedWallet/contracts/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
	"log"
	"math/big"
	"time"
)

// Monitor interface
type Monitor interface {
	Start(ctx context.Context, client *ethclient.Client) error
}

type monitor struct {
	contractAddress string
}

// AllowanceChangedEvent struct
type AllowanceChangedEvent struct {
	Event string `json:"event_type"`
	Sender string `json:"sender"`
	Beneficiary string `json:"beneficiary"`
	PrevAmount *big.Int `json:"prev_amount"`
	NewAmount *big.Int `json:"new_amount"`
	Timestamp time.Time `json:"timestamp"`
}

// MoneyReceivedEvent struct
type MoneyReceivedEvent struct {
	Event string `json:"event_type"`
	Sender string `json:"sender"`
	BlockNumber uint64 `json:"block_number"`
	Amount *big.Int `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// MoneySentEvent struct
type MoneySentEvent struct {
	Event string `json:"event_type"`
	Beneficiary string `json:"beneficiary"`
	BlockNumber uint64 `json:"block_number"`
	Amount *big.Int `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// OwnershipTransferredEvent struct
type OwnershipTransferredEvent struct {
	Event string `json:"event_type"`
	PreviousOwner string `json:"previous_owner"`
	NewOwner string `json:"new_owner"`
	BlockNumber uint64 `json:"block_number"`
	Timestamp time.Time `json:"timestamp"`
}

// NewMonitor returns a new runner instance
func NewMonitor(contractAddress string) Monitor {
	return &monitor{
		contractAddress: contractAddress,
	}
}

// Start register to listen blockchain events
func (m *monitor) Start(ctx context.Context, client *ethclient.Client) error {
	log.Printf("start monitoring at %s\n", m.contractAddress)

	err := validateContractAddress(ctx, client, m.contractAddress)
	if err != nil {
		return err
	}
	contract, err := contracts.NewContract(common.HexToAddress(m.contractAddress), client)
	if err != nil {
		return err
	}

	eg := new(errgroup.Group)
	eg.Go(func() error {
		return m.watchAllowanceChanged(ctx, contract)
	})
	eg.Go(func() error {
		return m.watchMoneySent(ctx, contract)
	})
	eg.Go(func() error {
		return m.watchMoneyReceived(ctx, contract)
	})
	eg.Go(func() error {
		return m.watchOwnershipTransferred(ctx, contract)
	})
	if err = eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (m *monitor) watchAllowanceChanged(ctx context.Context, contract *contracts.Contract) error {
	events := make(chan *contracts.ContractAllowanceChanged)
	opts := &bind.WatchOpts{
		Start:   nil,
		Context: ctx,
	}
	subscription, err := contract.WatchAllowanceChanged(opts, events, nil, nil)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case errChan := <-subscription.Err():
			return errChan
		case event := <-events:
			j, _ := json.MarshalIndent(
				AllowanceChangedEvent{
					Event:       "AllowanceChanged",
					Sender:      event.Sender.Hex(),
					Beneficiary: event.Beneficiary.Hex(),
					PrevAmount:  weiToEther(event.PrevAmount),
					NewAmount:   weiToEther(event.NewAmount),
					Timestamp: time.Now(),
				},
				"",
				"  ",
			)
			fmt.Println(string(j))
		}
	}
}

func (m *monitor) watchMoneySent(ctx context.Context, contract *contracts.Contract) error {
	events := make(chan *contracts.ContractMoneySent)
	opts := &bind.WatchOpts{
		Start:   nil,
		Context: ctx,
	}
	subscription, err := contract.WatchMoneySent(opts, events, nil)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case errChan := <-subscription.Err():
			return errChan
		case event := <-events:
			j, _ := json.MarshalIndent(
				MoneySentEvent{
					Event:       "MoneySent",
					Beneficiary: event.Beneficiary.Hex(),
					BlockNumber: event.Raw.BlockNumber,
					Amount:      weiToEther(event.Amount),
					Timestamp: time.Now(),
				},
				"",
				"  ",
			)
			fmt.Println(string(j))
		}
	}
}

func (m *monitor) watchMoneyReceived(ctx context.Context, contract *contracts.Contract) error {
	events := make(chan *contracts.ContractMoneyReceived)
	opts := &bind.WatchOpts{
		Start:   nil,
		Context: ctx,
	}
	subscription, err := contract.WatchMoneyReceived(opts, events, nil)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case errChan := <-subscription.Err():
			return errChan
		case event := <-events:
			j, _ := json.MarshalIndent(
				MoneyReceivedEvent{
					Event:       "MoneyReceived",
					Sender:      event.From.Hex(),
					BlockNumber: event.Raw.BlockNumber,
					Amount:      weiToEther(event.Amount),
					Timestamp: time.Now(),
				},
				"",
				"  ",
			)
			fmt.Println(string(j))
		}
	}
}

func (m *monitor) watchOwnershipTransferred(ctx context.Context, contract *contracts.Contract) error {
	events := make(chan *contracts.ContractOwnershipTransferred)
	opts := &bind.WatchOpts{
		Start:   nil,
		Context: ctx,
	}
	subscription, err := contract.WatchOwnershipTransferred(opts, events, nil, nil)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case errChan := <-subscription.Err():
			return errChan
		case event := <-events:
			j, _ := json.MarshalIndent(
				OwnershipTransferredEvent{
					Event:         "OwnershipTransferred",
					PreviousOwner: event.PreviousOwner.Hex(),
					NewOwner: event.NewOwner.Hex(),
					BlockNumber:   event.Raw.BlockNumber,
					Timestamp: time.Now(),
				},
				"",
				"  ",
			)
			fmt.Println(string(j))
		}
	}
}