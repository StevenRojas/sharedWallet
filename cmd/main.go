package main

import (
	"context"
	"github.com/StevenRojas/sharedWallet/cmd/command"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const closeFallbackTime = 3

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	defer func() {
		log.Println("closing application")
		signal.Stop(signals)
		cancel()
	}()

	go func() {
		<-signals
		log.Println("propagating cancel signal")
		cancel()
		time.Sleep(closeFallbackTime * time.Second)
		log.Println("fallback exit")
		os.Exit(1)
	}()

	cmd := command.NewRootCommand(ctx)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
