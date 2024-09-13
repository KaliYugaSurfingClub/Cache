package transaction

import (
	"bytes"
	"cache/core"
	"errors"
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

func TestWriteAndReadEvents(t *testing.T) {
	for _, test := range cases {
		mockFile := bytes.NewBuffer(nil)

		wErr := writeEventTo(mockFile, test.event)
		copyMockFile := mockFile.Bytes()
		event, rErr := readEvent(mockFile)

		//fmt.Println(test.name)
		//fmt.Println(test.event, copyMockFile, wErr)
		//fmt.Println(event, rErr)

		if wErr != nil && !errors.Is(wErr, ErrLongField) {
			t.Errorf("\ntest: %s\nunexpected error while writing: %s", test.name, wErr)
		}

		if rErr != nil && !errors.Is(rErr, ErrEmptyFile) {
			t.Errorf("\ntest: %s\nunexpected error while reading: %s", test.name, rErr)
		}

		if wErr != nil && len(copyMockFile) != 0 {
			t.Errorf("\ntest: %s\nunempty file buffer after writing error \nbuffer: %v", test.name, copyMockFile)
		}

		if wErr == nil && !reflect.DeepEqual(test.event, event) {
			t.Errorf("\ntest: %s\nexpected event to be %v\ngot %v", test.name, test.event, event)
		}
	}
}
