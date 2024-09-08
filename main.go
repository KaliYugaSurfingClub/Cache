package main

import (
	"cache/transport"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// todo do not catch errors from tl.errs
func main() {
	router := mux.NewRouter()

	err := transport.InitTransactionLogger()
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/v1/{key}", transport.KeyValuePutHandler).Methods(http.MethodPut)
	router.HandleFunc("/v1/{key}", transport.KeyValueGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/v1/{key}", transport.KeyValueDeleteHandler).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}
