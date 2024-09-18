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

type shutdownAble interface {
	Shutdown(ctx context.Context) error
}

func HandelShutdown(timeout time.Duration, services ...shutdownAble) {
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

	tl, err := transaction.NewLogger(cfg.LogsPath, cfg.Bandwidth)
	if err != nil {
		panic(err)
	}

	tl.Start()

	store := core.NewStore(tl)
	if err = store.Restore(); err != nil {
		panic(err)
	}

	server := frontend.NewRest(store, cfg.Port)

	go HandelShutdown(cfg.TimeForShutdown, server, tl)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
		panic(err)
	}
}
