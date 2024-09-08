package transaction

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

func encodeFirstString(bytes []byte, buf *bytes.Buffer) error {
	length := len(bytes)
	if length > 255 {
		return errors.New("key length exceeds 255 characters")
	}
	if err := buf.WriteByte(byte(length)); err != nil {
		return err
	}
	if _, err := buf.Write(bytes); err != nil {
		return err
	}

	return nil
}

func encodeEvent(e *Event) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if err := binary.Write(buf, binary.LittleEndian, e.Sequence); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, e.Type); err != nil {
		return nil, err
	}

	if err := encodeFirstString([]byte(e.Key), buf); err != nil {
		return nil, err
	}

	if err := encodeFirstString(e.Value, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeFirstString(reader *bytes.Reader) ([]byte, error) {
	length, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, nil
	}

	str := make([]byte, length)
	if _, err := reader.Read(str); err != nil {
		return nil, err
	}

	return str, nil
}

func decodeEvents(data []byte, events chan<- Event, errs chan<- error) {
	buf := bytes.NewReader(data)

	for buf.Len() > 0 {
		e := Event{}

		//todo wrap errors
		if err := binary.Read(buf, binary.LittleEndian, &e.Sequence); err != nil {
			errs <- err
			return
		}

		typeByte, err := buf.ReadByte()
		if err != nil {
			errs <- err
			return
		}

		e.Type = EventType(typeByte)

		keyBytes, err := decodeFirstString(buf)
		if err != nil {
			errs <- err
			return
		}

		e.Key = string(keyBytes)

		if e.Value, err = decodeFirstString(buf); err != nil {
			errs <- err
			return
		}

		fmt.Println(e)

		events <- e
	}
}
