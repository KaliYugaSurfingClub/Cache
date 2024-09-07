package transport

import (
	"cache/services"
	"cache/services/transactionLogger"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

//делаем путь host/v1/{key} чтобы get put delete запросы выглядели одинаково
//value передаем через body чтобы соответсвовать REST

func KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	value, err := services.Get(key)
	if errors.Is(err, services.ErrorNoSuchKey) {
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

func KeyValuePutHandler(tl transactionLogger.TransactionLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]

		value, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		err = services.Put(key, string(value), tl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func KeyValueDeleteHandler(tl transactionLogger.TransactionLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]

		err := services.Delete(key, tl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
	}
}
