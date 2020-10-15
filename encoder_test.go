package rj

import (
	"testing"
	"time"
)

type encodeTestCase struct {
	input    interface{}
	expected string
}

func TestEncodeVal(t *testing.T) {
	cases := []encodeTestCase{{input: "string value", expected: `"string value"`},
		{input: 123, expected: "123"},
		{input: true, expected: "true"},
		{input: 123.456, expected: "123.456"},
		{input: time.Date(2019, 10, 11, 12, 3, 4, 0, time.UTC),
			expected: "2019-10-11T12:03:04Z"},
	}

	for _, tc := range cases {
		e := newEncoder(tc.input)
		bts := e.encode()
		if string(bts) != tc.expected {
			t.Error("Test doEncode failed, input:", tc.input, ", expected:", tc.expected, ", got: ", string(bts))
		}
	}
}

func TestEncodeStruct(t *testing.T) {
	type stu struct {
		Name string
		Age  int
	}
	s := stu{Name: "Jimmy", Age: 12}
	e := newEncoder(s)
	bts := e.encode()
	out := `Name: "Jimmy"
Age: 12
`
	if string(bts) != out {
		t.Error("Test doEncode struct failed, expected: ", out, ", got: ", string(bts))
	}
}

func TestEncodeArray(t *testing.T) {
	in := []int{12, 23, 34}
	e := newEncoder(in)
	bts := e.doEncode()

	out := "[12,23,34]"
	if string(bts) != out {
		t.Error("Test doEncode struct failed, expected: ", out, ", got: ", string(bts))
	}
}
