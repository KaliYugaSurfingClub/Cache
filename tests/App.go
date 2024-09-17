package tests

import (
	"cache/core"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type App struct {
	location  string
	root      string
	port      string
	bandwidth int
	client    *http.Client
	cmd       *exec.Cmd
}

const localhost = "http://127.0.0.1"

func NewApp(location string) *App {
	return &App{
		location: location,
		root:     localhost + ":" + "8080",
		client:   &http.Client{},
	}
}

func (a *App) WithBandwidth(bandwidth int) *App {
	a.bandwidth = bandwidth
	return a
}

func (a *App) WithPort(port string) *App {
	a.root = localhost + ":" + port
	a.port = port
	return a
}

func (a *App) Start() {
	//compile app file
	if err := exec.Command("go", "build", a.location).Run(); err != nil {
		log.Fatal(err)
	}

	args := ""
	if a.port != "" {
		args += "-port=" + a.port
	}
	if a.bandwidth > 0 {
		args += "-bandwidth=" + strconv.Itoa(a.bandwidth)
	}

	//run app
	a.cmd = exec.Command("./main.exe", args)
	if err := a.cmd.Start(); err != nil {
		log.Fatal(err)
	}

	//time to start server
	time.Sleep(3 * time.Second)
}

func (a *App) Stop() {
	//not supported by windows(

	//if err := a.cmd.Process.Signal(syscall.SIGINT); err != nil {
	//	log.Fatal(err)
	//}
	//if err := a.cmd.Wait(); err != nil {
	//	log.Fatal(err)
	//}

	//bad way but it works
	if err := a.cmd.Process.Kill(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) CheckNoSuchKey(key string) error {
	return a.GetRequest(key, core.ErrorNoSuchKey.Error()+"\n")
}

func (a *App) GetRequest(key string, want string) error {
	url := a.root + "/v1/" + key

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if string(value) != want {
		return fmt.Errorf("for key: %q: got %q, want %q", key, value, want)
	}

	return nil
}

func (a *App) PutRequest(key, value string) error {
	url := a.root + "/v1/" + key

	req, err := http.NewRequest("PUT", url, strings.NewReader(value))
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("key: %q value %q not created", key, value)
	}

	defer resp.Body.Close()

	return nil
}

func (a *App) DeleteRequest(key string) error {
	url := a.root + "/v1/" + key

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (a *App) ClearRequest() error {
	url := a.root + "/v1/operation/clear"

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
