package transaction

import (
	"bytes"
)

type mockFile struct {
	*bytes.Buffer
}

func (mf mockFile) Close() error {
	return nil
}
