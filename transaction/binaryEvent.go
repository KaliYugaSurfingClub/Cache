package transaction

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

func writeNum(w *bufio.Writer, n uint64) error {
	binaryLen := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryLen, n)

	if _, err := w.Write(binaryLen[0 : slices.Index(binaryLen, 0)+1]); err != nil {
		return err
	}

	return nil
}

func readNum(r *bufio.Reader) (uint64, error) {
	var n uint64

	for i := 0; i < 8; i++ {
		b, err := r.ReadByte()
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

func writeString(w *bufio.Writer, str string) error {
	if err := writeNum(w, uint64(len(str))); err != nil {
		return err
	}

	if _, err := w.WriteString(str); err != nil {
		return err
	}

	return nil
}

func readString(r *bufio.Reader) (string, error) {
	length, err := readNum(r)
	if err != nil {
		return "", err
	}

	str := make([]byte, length)
	if _, err := io.ReadFull(r, str); err != nil {
		return "", err
	}

	return string(str), nil
}

func writeEventTo(w io.Writer, e core.Event) error {
	tmp := "write %s of event was failed: %w"
	buf := bufio.NewWriter(w)

	if err := writeNum(buf, e.ID); err != nil {
		return fmt.Errorf(tmp, "ID", err)
	}

	if err := buf.WriteByte(e.Type); err != nil {
		return fmt.Errorf(tmp, "type", err)
	}

	if buf.Available() < len(e.Key) {
		return ErrLongField
	}
	if err := writeString(buf, e.Key); err != nil {
		return fmt.Errorf(tmp, "key", err)
	}

	if buf.Available() < len(e.Value) {
		return ErrLongField
	}
	if err := writeString(buf, e.Value); err != nil {
		return fmt.Errorf(tmp, "value", err)
	}

	//todo fix that if buffer is overflow buf.Flush() will be called auto
	//Available позволяет не писать в буфер лишнего и предотвратить флаш (вроде)
	return buf.Flush()
}

func readEvent(r io.Reader) (e core.Event, err error) {
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
