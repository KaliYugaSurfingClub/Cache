package main

import (
	"cache/config"
	"cache/core"
	"cache/frontend"
	"cache/transaction"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type ShutdownAble interface {
	Shutdown(ctx context.Context) error
}

func handelShutdown(timeout time.Duration, services ...ShutdownAble) {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-ctx.Done()
	fmt.Println("Shutting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, service := range services {
		if err := service.Shutdown(ctx); err != nil {
			fmt.Println("Error closing service:", err)
		}
	}
}

func main() {
	cfg := config.Get()
	tl, err := transaction.NewLogger(cfg.LogsPath)
	if err != nil {
		panic(err)
	}

	store := core.NewStore(tl)
	store.Start()

	server := frontend.NewRest(store, cfg.Port)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			panic(err)
		}
	}()

	handelShutdown(cfg.TimeToShutdown, server, tl)
}
