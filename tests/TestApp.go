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

type AppStarter struct {
	location string
	root     string
	args     []string
	client   *http.Client
	cmd      *exec.Cmd
}

const localhost = "http://127.0.0.1"
const defaultPort = "8080"

func NewApp(location string) *AppStarter {
	return &AppStarter{
		location: location,
		root:     localhost + ":" + defaultPort,
		client:   &http.Client{},
	}
}

func (a *AppStarter) WithBandwidth(bandwidth int) *AppStarter {
	a.args = append(a.args, fmt.Sprintf("-bandwidth=%d", bandwidth))
	return a
}

func (a *AppStarter) WithPort(port string) *AppStarter {
	a.root = localhost + ":" + port
	a.args = append(a.args, fmt.Sprintf("-port=%s", port))
	return a
}

func (a *AppStarter) Start() {
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

func (a *AppStarter) Stop() {
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

func (a *AppStarter) CheckNoSuchKey(key string) error {
	return a.GetRequest(key, core.ErrorNoSuchKey.Error()+"\n")
}

func (a *AppStarter) GetRequest(key string, want string) error {
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

func (a *AppStarter) PutRequest(key, value string) error {
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

func (a *AppStarter) DeleteRequest(key string) error {
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

func (a *AppStarter) ClearRequest() error {
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
