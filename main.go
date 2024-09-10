package main

import (
	"cache/core"
	"cache/frontend"
	"cache/transaction"
)

func main() {
	tl, err := transaction.NewFileLogger("logs.bin")
	if err != nil {
		panic(err)
	}

	//todo check is it work with hard shutdown
	//defer store.Close()

	store := core.NewStore().WithTransactionLogger(tl)
	store.Restore()

	front := frontend.NewRest(store)
	front.Run()
}
