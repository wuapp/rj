package rj

import (
	"reflect"
	"testing"
	"time"
)

func arrayEquals(a, b interface{}) bool {
	switch va := a.(type) {
	case []string:
		l := len(va)
		vb := b.([]string)
		if len(vb) != l {
			return false
		}
		for i := 0; i < l; i++ {
			if va[i] != vb[i] {
				return false
			}
		}
	case []int:
		l := len(va)
		vb := b.([]int)
		if len(vb) != l {
			return false
		}
		for i := 0; i < l; i++ {
			if va[i] != vb[i] {
				return false
			}
		}
	case []float64:
		l := len(va)
		vb := b.([]float64)
		if len(vb) != l {
			return false
		}
		for i := 0; i < l; i++ {
			if va[i] != vb[i] {
				return false
			}
		}
	case []bool:
		l := len(va)
		vb := b.([]bool)
		if len(vb) != l {
			return false
		}
		for i := 0; i < l; i++ {
			if va[i] != vb[i] {
				return false
			}
		}
	case []time.Time:
		l := len(va)
		vb := b.([]time.Time)
		if len(vb) != l {
			return false
		}
		for i := 0; i < l; i++ {
			if va[i] != vb[i] {
				return false
			}
		}
	}
	return true
}

func dictEquals(a, b map[string]interface{}) bool {
	for k, v := range a {
		ta := reflect.TypeOf(v)

		if ta.Name() == "Time" {
			vt := v.(time.Time)
			return vt.Equal(b[k].(time.Time))
		}

		kind := ta.Kind()
		if kind == reflect.Slice || kind == reflect.Array {
			return arrayEquals(v, b[k])
		}

		switch vt := v.(type) {
		case []string:
			l := len(vt)
			val := v.([]string)
			if len(val) != l {
				return false
			}
			for i := 0; i < l; i++ {
				if vt[i] != val[i] {
					return false
				}
			}
		case []int:
			l := len(vt)
			val := v.([]int)
			if len(val) != l {
				return false
			}
			for i := 0; i < l; i++ {
				if vt[i] != val[i] {
					return false
				}
			}
		case []float64:
			l := len(vt)
			val := v.([]float64)
			if len(val) != l {
				return false
			}
			for i := 0; i < l; i++ {
				if vt[i] != val[i] {
					return false
				}
			}
		case []bool:
			l := len(vt)
			val := v.([]bool)
			if len(val) != l {
				return false
			}
			for i := 0; i < l; i++ {
				if vt[i] != val[i] {
					return false
				}
			}
		case []time.Time:
			l := len(vt)
			val := v.([]time.Time)
			if len(val) != l {
				return false
			}
			for i := 0; i < l; i++ {
				if !vt[i].Equal(val[i]) {
					return false
				}
			}
		default:
			if b[k] != v {
				return false
			}
		}
	}

	return true
}

func TestIsComment(t *testing.T) {
	cases := []testCase{
		{input: `# This is a comment`, expected: true},
		{input: `// This is also a comment`, expected: true},
		{input: `Not a comment`, expected: false},
	}

	var sc *scanner
	for _, c := range cases {
		sc = newScanner([]byte(c.input))
		if sc.isComment() != c.expected {
			t.Error("TestIsComment failed, input: ", c.input, ",, expected:", c.expected)
		}
	}
}

func TestIsBlankLine(t *testing.T) {
	cases := []testCase{
		{input: `# This is a comment`, expected: true},
		{input: `// This is also a comment`, expected: true},
		{input: ` 	`, expected: true},
		{input: `Not a blank line`, expected: false},
		{input: ` 	Not a blank line either`, expected: false},
	}

	var sc *scanner
	for _, c := range cases {
		sc = newScanner([]byte(c.input))
		if sc.isBlankLine() != c.expected {
			t.Error("TestIsComment failed, input: ", c.input, ", expected:", c.expected)
		}
	}
}

func TestScanRaw(t *testing.T) {
	cases := []testCase{
		{input: `abc`, expected: "abc"},
		{input: `123 #b`, expected: "123"},
		{input: "de\n", expected: "de"},
		{input: "fg,", expected: "fg"},
		{input: "jk]", expected: "jk"},
	}

	var sc *scanner
	var ret string
	for _, tc := range cases {
		sc = newScanner([]byte(tc.input))
		ret = sc.scanRaw()
		if ret != tc.expected.(string) {
			t.Error("ScanRaw failed, input:", tc.input, "expected:", tc.expected, "actual:", ret)
		}
	}
}

func TestScanQuotedString(t *testing.T) {
	cases := []testCase{
		{input: `"abc"`, expected: "abc"},
		{input: `"123\na"`, expected: "123\na"},
		{input: `"de\u6C49"`, expected: "de汉"},
		{input: `"123\na`, err: true},
		{input: `"123`, err: true},
	}

	var sc *scanner
	var ret string
	var err error
	for _, tc := range cases {
		sc = newScanner([]byte(tc.input))
		ret, err = sc.scanQuotedString()
		if tc.err {
			if err == nil {
				t.Error("ScanQuotedString failed, input:", tc.input, "expected error")
			}
		} else if ret != tc.expected.(string) {
			t.Error("ScanQuotedString failed, input:", tc.input, "expected:", tc.expected, "actual:", ret)
		}
	}
}
func TestScanPair(t *testing.T) {
	cases := []testCase{
		/*{input: `name: "str \n value"`, expected: map[string]interface{}{"name": "str \n value"}},
		{input: `age: 123`, expected: map[string]interface{}{"age": 123}},
		{input: `raw: ` + "`" + `adf \n df` + "`", expected: map[string]interface{}{"raw": "adf \\n df"}},
		{input: `b: true`, expected: map[string]interface{}{"b": "true"}},
		{input: `arrInt: [1,2,3]`, expected: map[string]interface{}{"arrInt": []int{1, 2, 3}}},*/
		{input: `arrStr: ["a","b",
				# Test line breaks and comments
				"c"]`, expected: map[string]interface{}{"arrStr": []string{"a", "b", "c"}}},
	}

	var sc *scanner

	for _, tc := range cases {
		sc = newScanner([]byte(tc.input))
		sc.scanPair(sc.root)
		if !dictEquals(sc.root.dict, tc.expected.(map[string]interface{})) {
			t.Error("Scan pair failed, input:", tc.input, "expected:", tc.expected, "actual:", sc.root.dict)
		}
	}
}

func TestScanObject(t *testing.T) {
	cases := []testCase{
		{input: `{
name: "Jason"
age: 12
}
`, expected: map[string]interface{}{"name": "Jason", "age": 12}},
		{input: `{ #comment
name: "Jason" #comment
age: 12
}
`, expected: map[string]interface{}{"name": "Jason", "age": 12}},
		{input: `{name: "Jason"
age: 12}
`, expected: map[string]interface{}{"name": "Jason", "age": 12}},
		{input: `{name: "Jason"`, err: true},
	}

	var sc *scanner
	var node *Node
	var err error

	for _, tc := range cases {
		sc = newScanner([]byte(tc.input))
		node, err = sc.scanObject()
		if tc.err {
			if err == nil {
				t.Error("ScanNode failed, there should be error.")
			}
		} else {
			if !dictEquals(node.dict, tc.expected.(map[string]interface{})) {
				t.Error("ScanNode failed, input:", tc.input, "expected:", tc.expected, "actual:", sc.root.dict)
			}
		}
	}

}

func TestScanNodeList(t *testing.T) {
	cases := []testCase{
		{input: `- name: "Jason"
age: 12 # comment
-name: "Abby"
age: 13
`, expected: []map[string]interface{}{{"name": "Jason", "age": 12}, {"name": "Abby", "age": 13}}},
	}

	var sc *scanner
	var err error

	for _, tc := range cases {
		sc = newScanner([]byte(tc.input))
		list := sc.scanNodeList()
		if tc.err {
			if err == nil {
				t.Error("ScanNode failed, there should be error.")
			}
		} else {
			if len(list) != 2 ||
				!dictEquals(list[0].dict, tc.expected.([]map[string]interface{})[0]) ||
				!dictEquals(list[1].dict, tc.expected.([]map[string]interface{})[1]) {
				t.Error("ScanNodeList failed, input:", tc.input, "expected:", tc.expected, "actual:", list[0], list[1])
			}
		}
	}

}

func TestScanNode(t *testing.T) {
	cases := []testCase{
		{input: `[TestNode]
		name: "Jason"
		age: 12
		`, expected: map[string]interface{}{"name": "Jason", "age": 12}},
		{input: "[TestNode", err: true},
	}

	var sc *scanner

	for _, tc := range cases {
		sc = newScanner([]byte(tc.input))
		sc.scanNode(sc.root)
		if tc.err {
			if sc.error == nil {
				t.Error("ScanNode failed, there should be error.")
			}
		} else {
			if !dictEquals(sc.root.dict["TestNode"].(*Node).dict, tc.expected.(map[string]interface{})) {
				t.Error("ScanNode failed, input:", tc.input, "expected:", tc.expected, "actual:", sc.root.dict)
			}
		}
	}

}

func TestScan(t *testing.T) {
	in := `name: "a\b"
array: ["a","b"]`
	s := newScanner([]byte(in))
	s.scan()
}
