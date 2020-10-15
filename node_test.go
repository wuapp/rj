package rj

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	node := NewNode()
	node.dict["name"] = "Zoe"

	s, err := node.Get("name")

	if err != nil {
		t.Error("Get failed, expected no error, got:", err)
	}

	if s.(string) != "Zoe" {
		t.Error("Get failed, expected: Zoe, got:", s)
	}

	_, err = node.Get("fake")

	if err == nil || err != errValueNotFound {
		t.Error("Get failed, should return error (value not found), got:", err)
	}
}

func TestGetString(t *testing.T) {
	node := NewNode()
	node.dict["name"] = "Zoe"
	node.dict["age"] = 12

	s, err := node.GetStringOrError("name")

	if err != nil {
		t.Error("GetStringOrError failed, expected no error, got:", err)
	}

	if s != "Zoe" {
		t.Error("GetStringOrError failed, expected: Zoe, got:", s)
	}

	_, err = node.GetStringOrError("age")

	if err == nil || err != errTypeMismatch {
		t.Error("GetStringOrError failed, should return error (type mismatch), got:", err)
	}

	s = node.GetStringOr("name", "Jim")

	if s != "Zoe" {
		t.Error("GetStringOr failed, expected: Zoe, got:", s)
	}

	s = node.GetStringOr("age", "Jim")

	if s != "Jim" {
		t.Error("GetStringOr failed, expected: Jim, got:", s)
	}

	s = node.GetString("name")

	if s != "Zoe" {
		t.Error("GetString failed, expected: Zoe, got:", s)
	}

	s = node.GetString("age")

	if s != "" {
		t.Error("GetString failed, expected empty string, got:", s)
	}

}

func TestGetInt(t *testing.T) {
	node := NewNode()
	node.dict["name"] = "Zoe"
	node.dict["age"] = 12

	i, err := node.GetIntOrError("age")

	if err != nil {
		t.Error("GetIntOrError failed, expected no error, got:", err)
	}

	if i != 12 {
		t.Error("GetIntOrError failed, expected: 12, got:", i)
	}

	_, err = node.GetIntOrError("name")

	if err == nil || err != errTypeMismatch {
		t.Error("GetIntOrError failed, should return error (type mismatch), got:", err)
	}

	i = node.GetIntOr("age", 10)

	if i != 12 {
		t.Error("GetIntOr failed, expected: 12, got:", i)
	}

	i = node.GetIntOr("name", 10)

	if i != 10 {
		t.Error("GetIntOr failed, expected: 10, got:", i)
	}

	i = node.GetInt("age")

	if i != 12 {
		t.Error("GetInt failed, expected: 12, got:", i)
	}

	i = node.GetInt("name")

	if i != 0 {
		t.Error("GetInt failed, expected empty value 0, got:", i)
	}
}

func TestGetBool(t *testing.T) {
	node := NewNode()
	node.dict["name"] = "Zoe"
	node.dict["passed"] = true

	b, err := node.GetBoolOrError("passed")

	if err != nil {
		t.Error("GetBoolOrError failed, expected no error, got:", err)
	}

	if !b {
		t.Error("GetBoolOrError failed, expected: true, got:", b)
	}

	_, err = node.GetBoolOrError("name")

	if err == nil || err != errTypeMismatch {
		t.Error("GetBoolOrError failed, should return error (type mismatch), got:", err)
	}

	b = node.GetBoolOr("passed", false)

	if !b {
		t.Error("GetBoolOr failed, expected: Zoe, got:", b)
	}

	b = node.GetBoolOr("name", false)

	if b {
		t.Error("GetBoolOr failed, expected: false, got:", b)
	}

	b = node.GetBool("passed")

	if !b {
		t.Error("GetBool failed, expected: true, got:", b)
	}

	b = node.GetBool("name")

	if b {
		t.Error("GetBool failed, expected false, got:", b)
	}

}

func TestGetFloat(t *testing.T) {
	node := NewNode()
	node.dict["name"] = "Zoe"
	node.dict["score"] = 82.5

	f, err := node.GetFloatOrError("score")

	if err != nil {
		t.Error("GetFloatOrError failed, expected no error, got:", err)
	}

	if f != 82.5 {
		t.Error("GetFloatOrError failed, expected: Zoe, got:", f)
	}

	_, err = node.GetFloatOrError("name")

	if err == nil || err != errTypeMismatch {
		t.Error("GetFloatOrError failed, should return error (type mismatch), got:", err)
	}

	f = node.GetFloatOr("score", 0)

	if f != 82.5 {
		t.Error("GetFloatOr failed, expected: Zoe, got:", f)
	}

	f = node.GetFloatOr("name", 23.4)

	if f != 23.4 {
		t.Error("GetFloatOr failed, expected: 23.4, got:", f)
	}

	f = node.GetFloat("score")

	if f != 82.5 {
		t.Error("GetFloat failed, expected: 82.5, got:", f)
	}

	f = node.GetFloat("name")

	if f != 0 {
		t.Error("GetFloat failed, expected 0, got:", f)
	}

}

func TestGetTime(t *testing.T) {
	node := NewNode()
	node.dict["name"] = "Zoe"

	expected := time.Now()
	node.dict["joined"] = expected

	tm, err := node.GetTimeOrError("joined")

	if err != nil {
		t.Error("GetTimeOrError failed, expected no error, got:", err)
	}

	if !tm.Equal(expected) {
		t.Error("GetTimeOrError failed, expected:", expected, ", got:", tm)
	}

	_, err = node.GetTimeOrError("name")

	if err == nil || err != errTypeMismatch {
		t.Error("GetTimeOrError failed, should return error (type mismatch), got:", err)
	}

	tm = node.GetTimeOr("joined", time.Time{})

	if !tm.Equal(expected) {
		t.Error("GetTimeOr failed, expected:", expected, ", got:", tm)
	}

	tm = node.GetTimeOr("name", time.Time{})

	if !tm.Equal(time.Time{}) {
		t.Error("GetFloatOr failed, expected:", time.Time{}, ", got:", tm)
	}

	tm = node.GetTime("joined")

	if !tm.Equal(expected) {
		t.Error("GetTime failed, expected:", expected, ", got:", tm)
	}

	tm = node.GetTime("name")

	if !tm.IsZero() {
		t.Error("GetTime failed, expected zero time, got:", tm)
	}

}

func TestNode_GetStringArray(t *testing.T) {
	node := NewNode()

	expected := []string{"ab", "cd"}
	node.dict["code"] = expected
	node.dict["score"] = []int{12, 34}

	arr, err := node.GetStringArrayOrError("code")

	if err != nil {
		t.Error("GetStringArrayOrError failed, expected no error, got:", err)
	} else {
		if !arrayEquals(arr, expected) {
			t.Error("GetStringArrayOrError failed, expected:", expected, ", got:", arr)
		}
	}

}
