package transaction

import (
	"bytes"
	"cache/core"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type Case struct {
	name  string
	event core.Event
}

var cases = []Case{
	{
		name: "default put",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "abc",
			Value: []byte("cba"),
		},
	},
	{
		name: "put with whitespaces",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "abc asdasd sadsad",
			Value: []byte("cba asdsad asdasd"),
		},
	},
	{
		name: "put with runes those bigger then byte",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "ASPOJSDA msapjodфыхв)(ЦЙУШО)шоыфывфтщ ыфвщтывфтфыъ-930032двй0-3939гцов-0392213ВЛТОВЫтштщNSOSDNu jsdnds==-2133221",
			Value: []byte("ASPOIJ{SDAJВЫОЫВФльokDJJ*(0pl/\\	p9838j3838hnnosaISA09102312301mSuыфшвыфвыфзывыфтвфгышвфывы09=ёё1221ё"),
		},
	},
	{
		name: "put too long",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   strings.Repeat("ASPOJSDA msapjodфыхв)(ЦЙУШО)шоыфывфтщ ыфвщтывфтфыъ-930032двй0-3939гцов-0392213ВЛТОВЫтштщNSOSDNu jsdnds==-2133221", 10),
			Value: []byte(strings.Repeat("ASPOJSDA msapjodфыхв)(ЦЙУШО)шоыфывфтщ ыфвщтывфтфыъ-930032двй0-3939гцов-0392213ВЛТОВЫтштщNSOSDNu jsdnds==-2133221", 10)),
		},
	},
	{
		name: "put empty key",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "",
			Value: []byte("avc"),
		},
	},
	{
		name: "put empty value",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "avc",
			Value: nil,
		},
	},
	{
		name: "put empty key and value",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "avc",
			Value: nil,
		},
	},
}

func init() {
	//transform all put cases to delete cases and add to slice of cases
	for _, c := range slices.Clone(cases) {
		newCase := Case{
			name: strings.Replace(c.name, "put", "delete", 1),
			event: core.Event{
				ID:    c.event.ID,
				Type:  core.EventDelete,
				Key:   c.event.Key,
				Value: nil,
			},
		}

		cases = append(cases, newCase)
	}
}

func writeAndRead(t *testing.T, c Case) {
	mockFile := bytes.NewBuffer(nil)
	writeErr := writeEventTo(mockFile, c.event)

	mockFileAfterWriting := mockFile.Bytes()
	event, readErr := readEvent(mockFile)

	info := fmt.Sprintf(
		"\ncase: %s\noriginal event: %v\nmock file after writing%v\nwriting error %s\ngot event: %v\n reading error %s",
		c.name, c.event, mockFileAfterWriting, writeErr, event, readErr,
	)

	if writeErr != nil && !errors.Is(writeErr, ErrLongField) {
		t.Errorf("unexpected error while writing" + info)
	}

	//we will get ErrEmptyFile, if we can not write
	if readErr != nil && !errors.Is(readErr, ErrEmptyFile) {
		t.Errorf("unexpected error while reading" + info)
	}

	//if we write with error, we should not mutate destination
	if writeErr != nil && len(mockFileAfterWriting) != 0 {
		t.Errorf("unempty file buffer after writing error" + info)
	}

	//if we write no error, wrote event and read event should be equal
	if writeErr == nil && !reflect.DeepEqual(c.event, event) {
		t.Errorf("expected event and got event are different" + info)
	}
}

func FuzzWriteReadRestore(f *testing.F) {
	for _, test := range cases {
		f.Add(test.name, test.event.ID, byte(test.event.Type), test.event.Key, test.event.Value)
	}

	f.Fuzz(func(t *testing.T, name string, ID uint64, eventType byte, key string, value []byte) {
		testCase := Case{
			name:  name,
			event: core.Event{ID: ID, Type: core.EventType(eventType), Key: key, Value: value},
		}

		writeAndRead(t, testCase)
	})
}
