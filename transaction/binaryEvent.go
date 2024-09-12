package transaction

import (
	"bytes"
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

func decodeString(r io.Reader, dest []byte) error {
	var length uint8
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return err
	}

	if length == 0 {
		return nil
	}

	if err := binary.Read(r, binary.LittleEndian, dest); err != nil {
		return err
	}

	return nil
}

func writeEventTo(w io.Writer, e core.Event) error {
	tmp := "write %s of event was failed: %w"

	buff := bytes.NewBuffer([]byte{})

	if err := binary.Write(buff, binary.LittleEndian, e.ID); err != nil {
		return fmt.Errorf(tmp, "ID", err)
	}

	if err := binary.Write(buff, binary.LittleEndian, e.Type); err != nil {
		return fmt.Errorf(tmp, "type", err)
	}

	if err := encodeString(buff, []byte(e.Key)); err != nil {
		return fmt.Errorf(tmp, "key", err)
	}

	if err := encodeString(buff, e.Value); err != nil {
		return fmt.Errorf(tmp, "value", err)
	}

	if _, err := buff.WriteTo(w); err != nil {
		return fmt.Errorf(tmp, "all fields", err)
	}

	return nil
}

func readEvent(r io.Reader) (core.Event, error) {
	tmp := "read %s of event was failed: %w"

	e := core.Event{}

	if err := binary.Read(r, binary.LittleEndian, &e.ID); err != nil {
		return e, fmt.Errorf(tmp, "id", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &e.Type); err != nil {
		return e, fmt.Errorf(tmp, "type", err)
	}

	keyBuff := make([]byte, 0)
	if err := decodeString(r, keyBuff); err != nil {
		return e, fmt.Errorf(tmp, "key", err)
	}

	e.Key = string(keyBuff)

	if err := decodeString(r, e.Value); err != nil {
		return e, fmt.Errorf(tmp, "value", err)
	}

	return e, nil
}
