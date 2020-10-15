package rj

import (
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const invalidTime = "invalid time value"

var (
	datetimeFormats = []string{"2006-01-02", "2006-01-02 15:04:05", "15:04:05", time.RFC3339}
	dateTimeRegs    = []*regexp.Regexp{regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`), //date only
		regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`),                            //datetime
		regexp.MustCompile(`^\d{2}:\d{2}:\d{2}(.\d+)?$`),                                       //time only
		regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(.\d+)?(Z|[+-]\d{2}:\d{2})?$`), //RFC3399
	}
)

func decodeDatetime(raw string) (val time.Time, err error) {
	for i := 0; i < 4; i++ {
		if dateTimeRegs[i].MatchString(raw) {
			return time.Parse(datetimeFormats[i], raw)
		}
	}

	return time.Time{}, newError(invalidTime)
}

//escape \uXXXX from
func escapeU4(in []byte) (out []byte, size int, err error) {
	var r rune

	if len(in) != 4 {
		err = newError(invalidUTF8StringValue)
		return
	}
	for i := 0; i < 4; i++ {
		c := in[i]
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			err = newError(invalidUTF8StringValue)
			return
		}
		r = r*16 + rune(c)
	}

	out = []byte{0, 0, 0, 0}
	size = utf8.EncodeRune(out, r)
	return
}

// decode a node to a field of type of struct or pointer to a struct
func decodeStructField(n *Node, ft reflect.Type) (v reflect.Value) {
	if ft.Kind() == reflect.Struct {
		v = reflect.New(ft).Elem()
		decodeNode(n, v)
	} else if ft.Kind() == reflect.Ptr {
		v = reflect.New(ft.Elem())
		decodeNode(n, v.Elem())
	}
	return
}

// decode node to a struct
func decodeNode(n *Node, rv reflect.Value) error {
	if rv.Kind() != reflect.Struct {
		return errTypeMismatch
	}

	for k, v := range n.dict {
		field := rv.FieldByName(k)
		if !field.IsValid() {
			field = rv.FieldByName(strings.Title(k))
		}
		if !field.IsValid() {
			continue
		}

		switch vt := v.(type) {
		case string:
			if field.Kind() != reflect.String {
				return errTypeMismatch
			}
			field.SetString(vt)
		case int:
			field.SetInt(int64(vt))
		case bool:
			field.SetBool(vt)
		case float64:
			field.SetFloat(vt)
		case time.Time:
			field.Set(reflect.ValueOf(vt))
		case *Node:
			s := decodeStructField(vt, field.Type())
			field.Set(s)
		case []string:
			l := len(vt)
			if l == 0 {
				return nil
			}

			arr := reflect.MakeSlice(field.Type(), l, l)
			for i := 0; i < l; i++ {
				arr.Index(i).Set(reflect.ValueOf(vt[i]))
			}
			field.Set(arr)
		case []int:
			l := len(vt)
			if l == 0 {
				return nil
			}

			arr := reflect.MakeSlice(field.Type(), l, l)
			for i := 0; i < l; i++ {
				arr.Index(i).Set(reflect.ValueOf(vt[i]))
			}
			field.Set(arr)
		case []bool:
			l := len(vt)
			if l == 0 {
				return nil
			}

			arr := reflect.MakeSlice(field.Type(), l, l)
			for i := 0; i < l; i++ {
				arr.Index(i).Set(reflect.ValueOf(vt[i]))
			}
			field.Set(arr)
		case []float64:
			l := len(vt)
			if l == 0 {
				return nil
			}

			arr := reflect.MakeSlice(field.Type(), l, l)
			for i := 0; i < l; i++ {
				arr.Index(i).Set(reflect.ValueOf(vt[i]))
			}
			field.Set(arr)
		case []time.Time:
			l := len(vt)
			if l == 0 {
				return nil
			}

			arr := reflect.MakeSlice(field.Type(), l, l)
			for i := 0; i < l; i++ {
				arr.Index(i).Set(reflect.ValueOf(vt[i]))
			}
			field.Set(arr)
		case []*Node:
			l := len(vt)
			if l == 0 {
				return nil
			}

			arr := reflect.MakeSlice(field.Type(), l, l)

			for i := 0; i < l; i++ {
				e := decodeStructField(vt[i], field.Type().Elem())
				arr.Index(i).Set(e)
			}
			field.Set(arr)
		}
	}
	return nil
}

func decode(n *Node, val interface{}) error {
	rv := reflect.ValueOf(val)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errTypeMismatch
	}

	rv = rv.Elem()

	return decodeNode(n, rv)
}
