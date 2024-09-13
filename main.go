package main

import (
	"cache/core"
	"cache/frontend"
	"cache/transaction"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type ShutdownAble interface {
	Shutdown(ctx context.Context) error
}

func handelShutdown(ctx context.Context, services ...ShutdownAble) {
	sigs := make(chan os.Signal)
	//todo notify with context
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigs
	fmt.Println("Shutting down ...")

	for _, service := range services {
		if err := service.Shutdown(ctx); err != nil {
			fmt.Println("Error closing service:", err)
		}
	}
}

// todo what happens if I terminate program while reading events from file
// the application simply closes the connection to the file, I suppose it is a safe behavior

//todo fix tests and check len of key and value in handlers or in core
//do new method of encoding events
//pass max len of key value thanks for config and do not use uint8 fo len

// todo cant read events from file_)))))))))
func main() {
	tl, err := transaction.NewFileLogger("logs.bin")
	if err != nil {
		panic(err)
	}

	store := core.NewStore().WithTransactionLogger(tl)
	store.Start()

	_ = frontend.NewRest(store).Start()

	for {

	}

	//ctx, _ := context.WithTimeout(context.Background(), 100000*time.Second)
	//do not change order, because the server needs open channels to complete all work
	//then we can close transactionLogger
	//handelShutdown(ctx, server, tl)
}
