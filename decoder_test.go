package rj

import (
	"testing"
	"time"
)

type testCase struct {
	input    string
	expected interface{}
	err      bool
}

func TestDecodeDatetime(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/shanghai")
	tests := []testCase{
		{input: "2018-08-08", expected: time.Date(2018, time.August, 8, 0, 0, 0, 0, time.UTC)},
		{input: "2019-09-09 13:54:46", expected: time.Date(2019, time.September, 9, 13, 54, 46, 0, time.UTC)},
		{input: "2019-09-09T13:54:46.123Z", expected: time.Date(2019, time.September, 9, 13, 54, 46, 1.23e8, time.UTC)},
		{input: "2019-09-09T13:54:46.123+08:00", expected: time.Date(2019, time.September, 9, 13, 54, 46, 1.23e8, loc)},
		{input: "2019-09-09T13:54:46.123+a", expected: time.Time{}},
	}

	for _, tc := range tests {
		tm, _ := decodeDatetime(tc.input)
		if !tm.Equal(tc.expected.(time.Time)) {
			t.Error("Decode time failed:", tc.input, "expected:", tc.expected, "actual:", tm)
		}
	}
}

func TestEscapeU4(t *testing.T) {
	r, size, err := escapeU4([]byte("6C49"))
	s := string(r[:size])
	if err != nil || s != "汉" {
		t.Error("Decode unicode failed, expected: 汉， got:", s, ", err:", err)
	}
	r, size, err = escapeU4([]byte("6c49x"))

	if err == nil {
		t.Error("Should be error")
	}
}

func TestDecode(t *testing.T) {
	type ts struct {
		Fs string
		Fi int
		Ff float64
		Fb bool
		Ft time.Time
	}
	expected := ts{Fs: "ABC", Fi: 12, Ff: 12.34, Fb: true, Ft: time.Now()}
	node := NewNode()
	node.dict["Fs"] = expected.Fs
	node.dict["Fi"] = expected.Fi
	node.dict["Ff"] = expected.Ff
	node.dict["Fb"] = expected.Fb
	node.dict["Ft"] = expected.Ft

	decoded := new(ts)
	err := decode(node, decoded)

	if err != nil {
		t.Error("Decode struct failed, err:", err)
		return
	}

	if decoded.Fs != expected.Fs {
		t.Error("Decode string field failed, expected:", expected.Fs, ", got:", decoded.Fs)
	}

	if decoded.Fi != expected.Fi {
		t.Error("Decode int field failed, expected:", expected.Fi, ", got:", decoded.Fi)
	}

	if decoded.Ff != expected.Ff {
		t.Error("Decode float field failed, expected:", expected.Ff, ", got:", decoded.Ff)
	}

	if decoded.Fb != expected.Fb {
		t.Error("Decode bool field failed, expected:", expected.Fb, ", got:", decoded.Fb)
	}

	if decoded.Ft != expected.Ft {
		t.Error("Decode time field failed, expected:", expected.Ft, ", got:", decoded.Ft)
	}
}

func TestDecodeArray(t *testing.T) {
	type ts struct {
		Fs []string
		Fi []int
		Ff []float64
		Fb []bool
		Ft []time.Time
	}
	expected := ts{Fs: []string{"ABC", "def"},
		Fi: []int{12, 34},
		Ff: []float64{12.34, 56.78},
		Fb: []bool{true, false},
		Ft: []time.Time{time.Now(), time.Now()}}
	node := NewNode()
	node.dict["Fs"] = expected.Fs
	node.dict["Fi"] = expected.Fi
	node.dict["Ff"] = expected.Ff
	node.dict["Fb"] = expected.Fb
	node.dict["Ft"] = expected.Ft

	decoded := new(ts)
	err := decode(node, decoded)

	if err != nil {
		t.Error("Decode array field failed, err:", err)
		return
	}

	if len(decoded.Fs) != 2 || decoded.Fs[0] != expected.Fs[0] || decoded.Fs[1] != expected.Fs[1] {
		t.Error("Decode string array failed, expected:", expected.Fs, ", got:", decoded.Fs)
	}

	if len(decoded.Fi) != 2 || decoded.Fi[0] != expected.Fi[0] || decoded.Fi[1] != expected.Fi[1] {
		t.Error("Decode int array failed, expected:", expected.Fi, ", got:", decoded.Fi)
	}

	if len(decoded.Ff) != 2 || decoded.Ff[0] != expected.Ff[0] || decoded.Ff[1] != expected.Ff[1] {
		t.Error("Decode float array failed, expected:", expected.Ff, ", got:", decoded.Ff)
	}

	if len(decoded.Fb) != 2 || decoded.Fb[0] != expected.Fb[0] || decoded.Fb[1] != expected.Fb[1] {
		t.Error("Decode bool array failed, expected:", expected.Fb, ", got:", decoded.Fb)
	}

	if len(decoded.Ft) != 2 || decoded.Ft[0] != expected.Ft[0] || decoded.Ft[1] != expected.Ft[1] {
		t.Error("Decode time array failed, expected:", expected.Ft, ", got:", decoded.Ft)
	}
}

func TestDecodeNested(t *testing.T) {
	type child struct {
		Fs string
	}
	type parent struct {
		StrChild child
		PtrChild *child
		Children []child
	}

	cStr := child{"ABC"}
	cPtr := child{"DEF"}
	expected := parent{cStr, &cPtr, []child{cStr, cPtr}}

	nStr := NewNode()
	nStr.dict["Fs"] = cStr.Fs

	nPtr := NewNode()
	nPtr.dict["Fs"] = cPtr.Fs

	/*nChildren := NewNode()
	nChildren.dict[""]
	*/
	node := NewNode()
	node.dict["StrChild"] = nStr
	node.dict["PtrChild"] = nPtr
	node.dict["Children"] = []*Node{nStr, nPtr}

	decoded := new(parent)
	err := decode(node, decoded)

	if err != nil {
		t.Error("Decode nested struct failed, err:", err)
		return
	}

	if expected.StrChild.Fs != decoded.StrChild.Fs {
		t.Error("Decode struct filed, expected:", expected.StrChild.Fs, ", got:", decoded.StrChild.Fs)
		return
	}

	if decoded.PtrChild == nil || expected.PtrChild.Fs != decoded.PtrChild.Fs {
		t.Error("Decode pointer filed, expected:", expected.PtrChild, ", got:", decoded.PtrChild)
		return
	}

	if decoded.Children == nil || len(decoded.Children) != 2 ||
		decoded.Children[0].Fs != expected.StrChild.Fs ||
		decoded.Children[1].Fs != expected.PtrChild.Fs {
		t.Error("Decode array of sub notes failed. expected:", expected.Children, ", got: ", decoded)
	}
}
