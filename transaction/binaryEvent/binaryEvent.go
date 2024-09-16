package binaryEvent

import (
	"bufio"
	"cache/core"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"slices"
)

var ErrLongField = errors.New("field is too long")
var ErrEmptyFile = errors.New("file is empty")

func writeNum(buf *bufio.Writer, n uint64) error {
	binaryLen := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryLen, n)

	if _, err := buf.Write(binaryLen[0 : slices.Index(binaryLen, 0)+1]); err != nil {
		return err
	}

	return nil
}

func readNum(buf *bufio.Reader) (uint64, error) {
	var n uint64

	for i := 0; i < 8; i++ {
		b, err := buf.ReadByte()
		if err != nil {
			return 0, err
		}
		if b == 0 {
			break
		}

		n += uint64(b)
	}

	return n, nil
}

func writeString(buf *bufio.Writer, str string) error {
	if buf.Available() <= len(str) {
		return ErrLongField
	}

	if err := writeNum(buf, uint64(len(str))); err != nil {
		return err
	}

	if _, err := buf.WriteString(str); err != nil {
		return err
	}

	return nil
}

func readString(buf *bufio.Reader) (string, error) {
	length, err := readNum(buf)
	if err != nil {
		return "", err
	}

	str := make([]byte, length)
	if _, err := io.ReadFull(buf, str); err != nil {
		return "", err
	}

	return string(str), nil
}

func WriteTo(w io.Writer, e core.Event) error {
	tmp := "write %s of event was failed: %w"
	buf := bufio.NewWriter(w)

	if err := writeNum(buf, e.ID); err != nil {
		return fmt.Errorf(tmp, "ID", err)
	}

	if err := buf.WriteByte(e.Type); err != nil {
		return fmt.Errorf(tmp, "type", err)
	}

	if err := writeString(buf, e.Key); err != nil {
		return fmt.Errorf(tmp, "key", err)
	}

	if err := writeString(buf, e.Value); err != nil {
		return fmt.Errorf(tmp, "value", err)
	}

	return buf.Flush()
}

func Read(r io.Reader) (e core.Event, err error) {
	tmp := "read %s of event was failed: %w"
	buf := bufio.NewReader(r)

	if _, err = buf.Peek(1); err != nil {
		return e, ErrEmptyFile
	}

	if e.ID, err = readNum(buf); err != nil {
		return e, fmt.Errorf(tmp, "id", err)
	}

	if e.Type, err = buf.ReadByte(); err != nil {
		return e, fmt.Errorf(tmp, "type", err)
	}

	if e.Key, err = readString(buf); err != nil {
		return e, fmt.Errorf(tmp, "key", err)
	}

	if e.Value, err = readString(buf); err != nil {
		return e, fmt.Errorf(tmp, "value", err)
	}

	return e, nil
}
