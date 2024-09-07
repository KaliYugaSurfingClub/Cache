package main

import (
	"cache/services/transactionLogger"
	"cache/transport"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	tl, err := transactionLogger.NewFileTransactionLogger("logs")
	if err != nil {
		log.Fatal(err)
	}

	tl.Start()

	go func() {
		for err := range tl.ErrCh() {
			fmt.Println(err)
		}
	}()

	router.HandleFunc("/v1/{key}", transport.KeyValueGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/v1/{key}", transport.KeyValuePutHandler(tl)).Methods(http.MethodPut)
	router.HandleFunc("/v1/{key}", transport.KeyValueDeleteHandler(tl)).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}
