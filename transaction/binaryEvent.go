package transaction

import (
	"cache/core"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var FieldIsTooLong = errors.New("field is too long")

func encodeString(w io.Writer, bytes []byte) error {
	if len(bytes) > 255 {
		return FieldIsTooLong
	}
	if err := binary.Write(w, binary.LittleEndian, uint8(len(bytes))); err != nil {
		return err
	}
	if _, err := w.Write(bytes); err != nil {
		return err
	}

	return nil
}

func decodeString(r io.Reader) ([]byte, error) {
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

func writingErrorf(field string, err error) error {
	return fmt.Errorf("write %s of event was faild: %w", field, err)
}

func writeEvent(w io.Writer, e core.Event) error {
	if err := binary.Write(w, binary.LittleEndian, e.Sequence); err != nil {
		return writingErrorf("id", err)
	}

	if err := binary.Write(w, binary.LittleEndian, e.Type); err != nil {
		return writingErrorf("type", err)
	}

	if err := encodeString(w, []byte(e.Key)); err != nil {
		return writingErrorf("key", err)
	}

	if err := encodeString(w, e.Value); err != nil {
		return writingErrorf("value", err)
	}

	return nil
}

func readingErrorf(field string, err error) error {
	return fmt.Errorf("read %s of event was faild: %w", field, err)
}

func readEvent(r io.Reader) (core.Event, error) {
	e := core.Event{}

	if err := binary.Read(r, binary.LittleEndian, &e.Sequence); err != nil {
		return e, readingErrorf("id", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &e.Type); err != nil {
		return e, readingErrorf("type", err)
	}

	keyBytes, err := decodeString(r)
	if err != nil {
		return e, readingErrorf("key", err)
	}

	e.Key = string(keyBytes)

	if e.Value, err = decodeString(r); err != nil {
		return e, readingErrorf("value", err)
	}

	return e, nil
}
