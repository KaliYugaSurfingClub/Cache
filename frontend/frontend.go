package frontend

import (
	"cache/core"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

type Rest struct {
	store *core.Store
}

func NewRest(store *core.Store) *Rest {
	return &Rest{store: store}
}

func (f *Rest) Start() *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/v1/{key}", f.KeyValuePutHandler).Methods(http.MethodPut)
	router.HandleFunc("/v1/{key}", f.KeyValueGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/v1/{key}", f.KeyValueDeleteHandler).Methods(http.MethodDelete)

	s := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			panic(err)
		}
	}()

	return &s
}

func (f *Rest) KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
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
	if _, err = w.Write(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func (f *Rest) KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	//todo debug
	time.Sleep(1 * time.Minute)

	f.store.Put(key, value)
	w.WriteHeader(http.StatusCreated)
}

func (f *Rest) KeyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	f.store.Delete(key)
}
