package transport

import (
	"cache/core"
	"cache/transaction"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var tl transaction.Logger

func InitTransactionLogger() error {
	var err error
	tl, err = transaction.NewFileLogger("logs.bin")

	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errs := tl.ReadEvents()
	e, err, ok := transaction.Event{}, nil, true

	for ok && err == nil {
		select {
		//ok == false means channel was closed
		case err, ok = <-errs:
		case e, ok = <-events:
			switch e.Type {
			case transaction.EventPut:
				err = core.Put(e.Key, e.Value)
			case transaction.EventDelete:
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

	w.Write(value)
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

	err = core.Put(key, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	tl.WritePut(key, value)

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
