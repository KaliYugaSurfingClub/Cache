package frontend

import (
	"cache/core"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Rest struct {
	store *core.Store
}

func NewRest(store *core.Store, port string) *http.Server {
	router := mux.NewRouter()
	f := &Rest{store}

	router.HandleFunc("/v1/{key}", f.Put).Methods(http.MethodPut)
	router.HandleFunc("/v1/{key}", f.Get).Methods(http.MethodGet)
	router.HandleFunc("/v1/{key}", f.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/v1/operation/clear", f.Clear).Methods(http.MethodDelete)

	s := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return &s
}

func (f *Rest) Get(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	value, err := f.store.Get(key)
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

	if _, err = w.Write([]byte(value)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func (f *Rest) Put(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	f.store.Put(key, string(value))
	w.WriteHeader(http.StatusCreated)
}

func (f *Rest) Delete(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	f.store.Delete(key)
}

func (f *Rest) Clear(w http.ResponseWriter, r *http.Request) {
	f.store.Clear()
}
