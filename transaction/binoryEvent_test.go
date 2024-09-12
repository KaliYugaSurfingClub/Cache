package transaction

import (
	"bytes"
	"cache/core"
	"errors"
	"io"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestWriteAndREadEvents(t *testing.T) {
	type Case struct {
		name     string
		badInput bool
		event    core.Event
	}

	cases := []Case{
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
			name:     "put with large runes",
			badInput: true,
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

	mockFile := bytes.NewBuffer(nil)

	for _, test := range cases {
		writeError := writeEventTo(mockFile, test.event)
		if writeError != nil && !errors.Is(writeError, FieldIsTooLong) {
			t.Errorf("\ntest: %s\nunexpected error while writing: %s", test.name, writeError)
		}

		event, readError := readEvent(mockFile)
		if readError != nil && !errors.Is(readError, io.EOF) {
			t.Errorf("\ntest: %s\nunexpected error while reading: %s", test.name, readError)
		}

		if writeError != nil && (mockFile.Len() != 0 || !errors.Is(readError, io.EOF)) {
			t.Errorf("\ntest: %s\nunempty file buffer after writing error", test.name)
		}

		if test.badInput && !reflect.DeepEqual(test.event, event) {
			t.Errorf("\ntest: %s\nexpected event to be %v\ngot %v", test.name, test.event, event)
		}
	}
}
