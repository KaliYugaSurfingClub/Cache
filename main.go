package main

import (
	"cache/core"
	"cache/frontend"
	"cache/transaction"
)

// todo do not catch errors from tl.errs
func main() {
	tl, err := transaction.NewFileLogger("logs.bin")
	if err != nil {
		panic(err)
	}

	//todo check is it work with hard shutdown
	//defer store.Close()

	store := core.NewStore().WithTransactionLogger(tl)

	if err := store.Restore(); err != nil {
		panic(err)
	}

	front := frontend.NewRest(store)
	front.Run()
}
