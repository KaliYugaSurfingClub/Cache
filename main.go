package main

import (
	"cache/core"
	"cache/transactionLogger"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

var tl transactionLogger.TransactionLogger

func initTransactionLogger() error {
	var err error
	tl, err = transactionLogger.NewFileTransactionLogger("logs.log")

	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errs := tl.ReadEvents()
	e, err, ok := transactionLogger.Event{}, nil, true

	for ok && err == nil {
		select {
		//ok == false means channel was closed
		case err, ok = <-errs:
		case e, ok = <-events:
			switch e.Type {
			case transactionLogger.EventPut:
				err = core.Put(e.Key, e.Value)
			case transactionLogger.EventDelete:
				err = core.Delete(e.Key)
			}
		}
	}

	tl.Start()

	return err
}

//делаем путь host/v1/{key} чтобы get put delete запросы выглядели одинаково
//value передаем через body чтобы соответсвовать REST

func KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	value, err := core.Get(key)
	if errors.Is(err, core.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		fmt.Println(err)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Write([]byte(value))
}

func KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	err = core.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	tl.WritePut(key, string(value))

	w.WriteHeader(http.StatusCreated)
}

func KeyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	err := core.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	tl.WriteDelete(key)
}

// todo do not catch errors from tl.errs
func main() {
	router := mux.NewRouter()

	err := initTransactionLogger()
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/v1/{key}", KeyValuePutHandler).Methods(http.MethodPut)
	router.HandleFunc("/v1/{key}", KeyValueGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/v1/{key}", KeyValueDeleteHandler).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}
