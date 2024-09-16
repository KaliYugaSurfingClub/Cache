package main

import (
	"cache/core"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const port = "http://127.0.0.1:8080"

var cmd *exec.Cmd

func init() {
	//todo
	cmd = exec.Command("go", "run", "../main.go")

	if err := cmd.Start(); err != nil {
		os.Exit(1)
	}

	time.Sleep(3 * time.Second)
}

type request struct {
	key   string
	value string
}

func TestPutGet(t *testing.T) {
	req := request{"first", "some value"}

	if err := putRequest(req.key, req.value); err != nil {
		t.Fatal(err)
	}

	resp, err := getRequest(req.key)
	if err != nil {
		t.Fatal(err)
	}
	if resp != req.value {
		t.Fatal(resp)
	}
}

func TestNoSuchKey(t *testing.T) {
	resp, err := getRequest("do not exists")
	if err != nil {
		t.Fatal(err)
	}
	if resp != core.ErrorNoSuchKey.Error()+"\n" {
		t.Fatal(resp)
	}
}

func TestDelete(t *testing.T) {
	req := request{"first", "some value"}

	if err := putRequest(req.key, req.value); err != nil {
		t.Fatal(err)
	}

	if err := deleteRequest(req.key); err != nil {
		t.Fatal(err)
	}

	resp, err := getRequest(req.value)
	if err != nil {
		t.Fatal(err)
	}
	if resp != core.ErrorNoSuchKey.Error()+"\n" {
		t.Fatal(resp)
	}
}

func putRequest(key, value string) error {
	url := port + "/v1/" + key

	if _, err := http.NewRequest("PUT", url, strings.NewReader(value)); err != nil {
		return err
	}

	return nil
}

func getRequest(key string) (string, error) {
	url := port + "/v1/" + key

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	value, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	return string(value), nil
}

func deleteRequest(key string) error {
	url := port + "/v1/" + key

	if _, err := http.Get(url); err != nil {
		return err
	}

	return nil
}
