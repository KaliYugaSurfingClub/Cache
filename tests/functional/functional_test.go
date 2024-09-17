package functional

//import only errors from core

import (
	"cache/tests"
	"testing"
)

type request struct {
	key   string
	value string
}

func TestBasicCases(t *testing.T) {
	a := tests.NewApp("../../main.go").WithPort("9989")

	a.Start()
	defer a.Stop()

	req1 := request{"first key", "some value"}
	req2 := request{"second key", "another value"}

	t.Run("no such key", func(t *testing.T) {
		if err := a.CheckNoSuchKey("aklsasdkpodsa"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("put", func(t *testing.T) {
		if err := a.PutRequest(req1.key, req1.value); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get", func(t *testing.T) {
		if err := a.GetRequest(req1.key, req1.value); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("delete", func(t *testing.T) {
		if err := a.DeleteRequest(req1.key); err != nil {
			t.Fatal(err)
		}

		if err := a.CheckNoSuchKey(req1.key); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("clear", func(t *testing.T) {
		if err := a.PutRequest(req1.key, req1.value); err != nil {
			t.Fatal(err)
		}

		if err := a.PutRequest(req2.key, req2.value); err != nil {
			t.Fatal(err)
		}

		if err := a.ClearRequest(); err != nil {
			t.Fatal(err)
		}

		if err := a.CheckNoSuchKey(req1.key); err != nil {
			t.Fatal(err)
		}

		if err := a.CheckNoSuchKey(req2.key); err != nil {
			t.Fatal(err)
		}
	})
}
