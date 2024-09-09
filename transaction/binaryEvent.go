package transaction

import (
	"cache/core"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

func encodeString(w *os.File, bytes []byte) error {
	length := uint8(len(bytes))
	if length > 255 {
		return errors.New("key length exceeds 255 characters")
	}
	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return err
	}
	if _, err := w.Write(bytes); err != nil {
		return err
	}

	return nil
}

func encodeEvent(e core.Event, w *os.File) error {
	if err := binary.Write(w, binary.LittleEndian, e.Sequence); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, e.Type); err != nil {
		return err
	}

	if err := encodeString(w, []byte(e.Key)); err != nil {
		return err
	}

	if err := encodeString(w, e.Value); err != nil {
		return err
	}

	return nil
}

func decodeFirstString(r io.Reader) ([]byte, error) {
	var length uint8
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, nil
	}

	str := make([]byte, length)
	if err := binary.Read(r, binary.LittleEndian, &str); err != nil {
		return nil, err
	}

	return str, nil
}

func decodeEvent(r io.Reader) (core.Event, error) {
	e := core.Event{}

	//todo wrap errors
	if err := binary.Read(r, binary.LittleEndian, &e.Sequence); err != nil {
		return e, err
	}

	if err := binary.Read(r, binary.LittleEndian, &e.Type); err != nil {
		return e, err
	}

	//todo refactor
	keyBytes, err := decodeFirstString(r)
	if err != nil {
		return e, err
	}

	e.Key = string(keyBytes)

	if e.Value, err = decodeFirstString(r); err != nil {
		return e, err
	}

	fmt.Println(e)

	return e, nil
}
