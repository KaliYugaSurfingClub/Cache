package main

import (
	"cache/core"
	"cache/frontend"
	"cache/transaction"
	"errors"
	"fmt"
	"io"
)

func initializeTransactionLog(store *core.Store) (*transaction.Logger, error) {
	transact, err := transaction.NewFileLogger("logs.bin")
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction logger: %w", err)
	}

	events, errs := transact.ReadEvents()
	ok, e := true, core.Event{}

	for ok && err == nil {
		select {
		case err, ok = <-errs:

		case e, ok = <-events:
			switch e.Type {
			case core.EventDelete: // Got a DELETE event!
				store.Delete(e.Key)
			case core.EventPut: // Got a PUT event!
				store.Put(e.Key, e.Value)
			}
		}
	}

	if errors.Is(err, io.EOF) {
		return transact, nil
	}

	if err != nil {
		return nil, err
	}

	return transact, nil
}

// todo do not catch errors from tl.errs
func main() {
	//tl, err := transaction.NewFileLogger("logs.bin")
	//if err != nil {
	//	panic(err)
	//}

	//todo check is it work with hard shutdown
	//defer store.Close()
	store := core.NewStore()
	tl, err := initializeTransactionLog(store)
	if err != nil {
		fmt.Println(err)
	}
	store.WithTransactionLogger(tl)
	//if err := store.Restore(); err != nil {
	//	panic(err)
	//}

	tl.Start()

	front := frontend.NewRest(store)
	front.Run()
}
