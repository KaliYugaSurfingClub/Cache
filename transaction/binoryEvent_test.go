package transaction

import (
	"bytes"
	"cache/core"
	"errors"
	"fmt"
	"math"
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
			Value: "cba",
		},
	},
	{
		name: "put with whitespaces",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "abc asdasd sadsad",
			Value: "cba asdsad asdasd",
		},
	},
	{
		name: "put with runes those bigger then byte",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "ASPOJSDA msapjodфыхв)(ЦЙУШО)шоыфывфтщ ыфвщтывфтфыъ-930032двй0-3939гцов-0392213ВЛТОВЫтштщNSOSDNu jsdnds==-2133221",
			Value: "ASPOIJ{SDAJВЫОЫВФльokDJJ*(0pl/\\	p9838j3838hnnosaISA09102312301mSuыфшвыфвыфзывыфтвфгышвфывы09=ёё1221ё",
		},
	},
	{
		name: "put too long",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   strings.Repeat("a", math.MaxInt32/300),
			Value: strings.Repeat("b", math.MaxInt32/300),
		},
	},
	{
		name: "put empty key",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "",
			Value: "avc",
		},
	},
	{
		name: "put empty value",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "avc",
			Value: "",
		},
	},
	{
		name: "put empty key and value",
		event: core.Event{
			ID:    14,
			Type:  core.EventPut,
			Key:   "avc",
			Value: "",
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
				Value: "",
			},
		}

		cases = append(cases, newCase)
	}
}

func summarizeString(str string, limit int) string {
	return summarizeSlice([]rune(str), limit)
}

func summarizeSlice[T any](slice []T, limit int) string {
	if len(slice) > limit {
		return fmt.Sprintf("%v", slice[0:limit/2]) + "..." + fmt.Sprintf("%v", slice[len(slice)-limit/2:])
	}

	return fmt.Sprintf("%v", slice)
}

func writeAndRead(t *testing.T, c Case) {
	mockFile := bytes.NewBuffer(nil)
	writeErr := writeEventTo(mockFile, c.event)

	mockFileAfterWriting := mockFile.Bytes()
	gotEvent, readErr := readEvent(mockFile)

	mockFileAfterWritingForPrint := summarizeSlice(mockFileAfterWriting, 20)

	expectedEventForPrint := core.Event{
		ID:    c.event.ID,
		Type:  c.event.Type,
		Key:   summarizeString(c.event.Key, 20),
		Value: summarizeString(c.event.Value, 20),
	}

	gotEventForPrint := core.Event{
		ID:    gotEvent.ID,
		Type:  gotEvent.Type,
		Key:   summarizeString(gotEvent.Key, 20),
		Value: summarizeString(gotEvent.Value, 20),
	}

	info := fmt.Sprintf(
		"\ncase: %s\nexpected event: %v\nmock file after writing%s\nwriting error %s\ngot event: %v\nreading error: %s",
		c.name, expectedEventForPrint,
		mockFileAfterWritingForPrint,
		writeErr, gotEventForPrint, readErr,
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

	//if we write no error, wrote gotEvent and read gotEvent should be equal
	if writeErr == nil && !reflect.DeepEqual(c.event, gotEvent) {
		t.Errorf("expected gotEvent and got gotEvent are different" + info)
	}
}

func FuzzWriteReadRestore(f *testing.F) {
	for _, test := range cases {
		f.Add(test.name, test.event.ID, test.event.Type, test.event.Key, test.event.Value)
	}

	f.Fuzz(func(t *testing.T, name string, ID uint64, eventType byte, key string, value string) {
		testCase := Case{
			name:  name,
			event: core.Event{ID: ID, Type: eventType, Key: key, Value: value},
		}

		writeAndRead(t, testCase)
	})
}
