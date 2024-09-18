package tests

import (
	"cache/core"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type TestingApp struct {
	location string
	root     string
	args     []string
	cmd      *exec.Cmd
}

const localhost = "http://127.0.0.1"
const defaultPort = "8080"

func NewApp(location string) *TestingApp {
	return &TestingApp{
		location: location,
		root:     localhost + ":" + defaultPort,
	}
}

func (a *TestingApp) WithBandwidth(bandwidth int) *TestingApp {
	a.args = append(a.args, fmt.Sprintf("-bandwidth=%d", bandwidth))
	return a
}

func (a *TestingApp) WithPort(port string) *TestingApp {
	a.root = localhost + ":" + port
	a.args = append(a.args, fmt.Sprintf("-port=%s", port))
	return a
}

func (a *TestingApp) Start() {
	//compile app file
	if err := exec.Command("go", "build", a.location).Run(); err != nil {
		log.Fatal(err)
	}

	executable := "./" + strings.TrimSuffix(filepath.Base(a.location), filepath.Ext(a.location)) + ".exe"

	a.cmd = exec.Command(executable, a.args...)

	//run app
	if err := a.cmd.Start(); err != nil {
		log.Fatal(err)
	}

	//time to start server
	time.Sleep(3 * time.Second)
}

func (a *TestingApp) Stop() {
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

func (a *TestingApp) CheckNoSuchKey(key string) error {
	return a.CheckGetRequest(key, core.ErrorNoSuchKey.Error()+"\n")
}

func (a *TestingApp) GetRequest(key string) (string, error) {
	url := a.root + "/v1/" + key

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(value), err
}

func (a *TestingApp) CheckGetRequest(key string, want string) error {
	value, err := a.GetRequest(key)
	if err != nil {
		return err
	}

	if value != want {
		return fmt.Errorf("for key: %q: got %q, want %q", key, value, want)
	}

	return nil
}

func (a *TestingApp) PutRequest(key string, value string) error {
	url := a.root + "/v1/" + key

	req, err := http.NewRequest("PUT", url, strings.NewReader(value))
	if err != nil {
		return err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("key: %q value %q not created", key, value)
	}

	defer resp.Body.Close()

	return nil
}

func (a *TestingApp) DeleteRequest(key string) error {
	url := a.root + "/v1/" + key

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (a *TestingApp) ClearRequest() error {
	url := a.root + "/v1/operation/clear"

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
