package main

import (
	"cache/core"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const port = "5678"
const root = "http://127.0.0.1:" + port

var client = &http.Client{}

func TestMain(m *testing.M) {
	//compile app file
	if err := exec.Command("go", "build", "../main.go").Run(); err != nil {
		log.Fatal(err)
	}

	//run app
	cmd := exec.Command("./main.exe", "-port="+port)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	//todo refactor
	//time to run server
	time.Sleep(2 * time.Second)

	exitCode := m.Run()

	if err := cmd.Process.Kill(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}

type request struct {
	key   string
	value string
}

func TestPutGet(t *testing.T) {
	req := request{"first key", "some value"}

	if err := putRequest(req.key, req.value); err != nil {
		t.Fatal(err)
	}

	resp, err := getRequest(req.key)
	if err != nil {
		t.Fatal(err)
	}
	if resp != req.value {
		t.Error(resp)
	}
}

func TestNoSuchKey(t *testing.T) {
	resp, err := getRequest("do not exists")
	if err != nil {
		t.Fatal(err)
	}
	if resp != core.ErrorNoSuchKey.Error()+"\n" {
		t.Error(resp)
	}
}

func TestDelete(t *testing.T) {
	req := request{"first key", "some value"}

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
		t.Error(resp)
	}
}

func getRequest(key string) (string, error) {
	url := root + "/v1/" + key

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

func putRequest(key, value string) error {
	url := root + "/v1/" + key

	req, err := http.NewRequest("PUT", url, strings.NewReader(value))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New("not created")
	}

	defer resp.Body.Close()

	return nil
}

func deleteRequest(key string) error {
	url := root + "/v1/" + key

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
