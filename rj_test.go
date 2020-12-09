package rj

import "testing"

func TestParse(t *testing.T) {
	in := `name: "a\n"
array: ["a","b"]`
	node, err := ParseString(in)
	if err != nil {
		t.Error("ParseString failed, expected no error, got:", err)
	} else {
		n := node.GetString("name")
		if n != "a\n" {
			t.Error("ParseString failed, expected: a\n, got: ", n)
		}

		array := node.GetStringArray("array")
		expected := []string{"a", "b"}
		if !arrayEquals(array, expected) {
			t.Error("ParseString failed, expected:", expected, ", got: ", array)
		}
	}

	in = `name: "a\c"
array: ["a","b"]

[ANode]
name: "Anna"`
	node, err = ParseString(in)
	if err != nil {
		name := node.GetString("name")
		if name != "" {
			t.Error("ParseString failed, expect empty string because of escape error, but got:", name)
		}

		arr := node.GetStringArray("array")
		if !arrayEquals(arr, []string{"a", "b"}) {
			t.Error("ParseString failed, expect correct array since we should skip error lines")
		}

		n, err := node.GetNode("ANode")
		if err != nil {
			t.Error("Parse node failed")
		}

		name = n.GetString("name")
		if name != "Anna" {
			t.Error("Get sub node failed, expect: Anna, got:", name)
		}
	} else {
		t.Error("ParseString failed, expect invalid escape error, but got nil")
	}
}

func TestMarshal(t *testing.T) {
	type s struct {
		Name string
		Age  int
	}
	in := `Name: "abc"
Age: 12
`
	s1 := new(s)
	err := Unmarshal([]byte(in), s1)

	if err == nil {
		if s1.Name != "abc" {
			t.Error("Test unmarshal failed, input:", in, ", expected s1.Name: abc, got:", s1.Name)
		}

		if s1.Age != 12 {
			t.Error("Test unmarshal failed, input:", in, ", expected s1.Age: 12, got:", s1.Age)
		}

		bytes := Marshal(s1)

		if string(bytes) != in {
			t.Error("Test marshal failed, expected: ", in, ", got:", string(bytes))
		}
	} else {
		t.Error("Test unmarshal failed, expected no error, got:", err)
	}

}
