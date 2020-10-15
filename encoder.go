package rj

import (
	"bytes"
	"reflect"
	"strconv"
	"time"
)

type encoder struct {
	val interface{}
	*bytes.Buffer
}

func newEncoder(v interface{}) *encoder {
	return &encoder{v, new(bytes.Buffer)}
}

func (e *encoder) encode() []byte {
	rv := reflect.ValueOf(e.val)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Struct {
		vt := rv.Type()
		if vt.Name() == "Time" {
			e.encodeTime(rv)
		} else {
			e.encodeFields(rv, vt)
		}
	} else {
		e.encodeVal(rv)
	}
	return e.Bytes()
}

func (e *encoder) doEncode() []byte {
	rv := reflect.ValueOf(e.val)
	e.encodeVal(rv)
	return e.Bytes()
}

func (e *encoder) encodeVal(v reflect.Value) {
	switch v.Kind() {
	case reflect.String:
		e.WriteByte('"')
		e.WriteString(v.String())
		e.WriteByte('"')
	case reflect.Bool:
		if v.Bool() {
			e.WriteString("true")
		} else {
			e.WriteString("false")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b := strconv.AppendInt([]byte{}, v.Int(), 10)
		e.Write(b)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		b := strconv.AppendUint([]byte{}, v.Uint(), 10)
		e.Write(b)
	case reflect.Float32, reflect.Float64:
		b := strconv.AppendFloat([]byte{}, v.Float(), 'f', -1, 64)
		e.Write(b)
	case reflect.Struct:
		e.encodeStruct(v)
	case reflect.Slice, reflect.Array:
		e.encodeArray(v)
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			e.WriteString("null")
		} else {
			e.encodeVal(v.Elem())
		}
	}
}

func (e *encoder) encodeTime(v reflect.Value) {
	t := v.Interface().(time.Time)
	b := make([]byte, 0, len(time.RFC3339Nano))
	e.Write(t.AppendFormat(b, time.RFC3339Nano))
}

func (e *encoder) encodeStruct(v reflect.Value) {
	vt := v.Type()

	e.WriteString(vt.Name())
	e.WriteByte(':')
	e.WriteByte('{')

	e.encodeFields(v, vt)

	e.WriteByte('}')
}

func (e *encoder) encodeFields(v reflect.Value, vt reflect.Type) {
	n := v.NumField()
	for i := 0; i < n; i++ {
		f := v.Field(i)
		ft := f.Type()

		if ft.Name() == "" && ft.Kind() == reflect.Ptr {
			//ft = ft.Elem()
			f = f.Elem()
		}

		e.WriteString(vt.Field(i).Name)
		e.WriteString(": ")
		e.encodeVal(f)
		e.WriteByte('\n')

	}
}

func (e *encoder) encodeArray(v reflect.Value) {
	n := v.Len()
	e.WriteByte('[')
	for i := 0; i < n; i++ {
		if i != 0 {
			e.WriteByte(',')
		}
		e.encodeVal(v.Index(i))
	}
	e.WriteByte(']')
}
