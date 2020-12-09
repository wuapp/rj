package rj

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	delimiter              = ':'
	invalidNodeName        = "invalid node name"
	invalidName            = "invalid name"
	invalidValue           = "invalid value"
	invalidArray           = "invalid array"
	invalidStringValue     = "invalid string value"
	invalidUTF8StringValue = "invalid utf-8 string value"
	invalidBoolValue       = "invalid bool value"
	invalidNullValue       = "invalid null value"
	invalidEscape          = "invalid escape"
	invalidObject          = "invalid object"
)

type scanState int

const (
	scanOk scanState = iota
	scanError
	eof
	endOfArray
	endOfItem
)

type scanner struct {
	data   []byte
	len    int
	offset int
	root   *Node
	error  *RJError
	state  scanState
	line   int
}

func newScanner(in []byte) *scanner {
	node := NewNode()
	return &scanner{data: in, len: len(in), root: node, line: 1 /*,currentNode:node*/}
}

func (s *scanner) char() byte {
	return s.data[s.offset]
}

func (s *scanner) addErrorMsg(msg string) {
	if s.error == nil {
		s.error = newError(msg)
	} else {
		s.error.addMsg(msg)
	}
}

func (s *scanner) addError(err error) {
	s.addErrorMsg(err.Error())
}

func (s *scanner) scan() {
	for s.offset < s.len-1 {
		s.skip()

		if s.data[s.offset] == '[' {
			s.scanNode(s.root)
		} else if s.isComment() {
			s.skipRestOfLine()
		} else {
			s.scanPair(s.root)
		}

		s.skipRestOfLine()
	}
}

func (s *scanner) scanPair(parent *Node) {
	name := s.scanName()
	if name == "" {
		s.addErrorMsg(invalidName)
		s.skipRestOfLine()
		return
	}

	s.skipSpace()
	val, err := s.scanValue()
	if err != nil {
		s.addErrorMsg(err.Error() + ", name: " + name)
	} else {
		//s.addPair(name,val)
		parent.dict[name] = val
	}
}

func (s *scanner) scanName() (name string) {
	for i := s.offset; i < s.len; i++ {
		c := s.data[i]
		if isSpace(c) || c == delimiter || c == ']' {
			name = string(s.data[s.offset:i])
			s.offset = i + 1
			return
		} else if isLineEnd(c) {
			return ""
		}
	}

	return ""
}

func (s *scanner) scanValue() (val interface{}, err error) {
	//s.skip()
	c := s.data[s.offset]
	switch {
	case c == '"':
		return s.scanQuotedString()
	case c == '`':
		return s.scanRawString()
	case (c >= '0' && c <= '9') || c == '+' || c == '-':
		raw := s.scanRaw()

		val, err = strconv.Atoi(raw)
		if err == nil {
			return
		}

		val, err = strconv.ParseFloat(raw, 64)
		if err == nil {
			return
		}

		return decodeDatetime(raw)

	case c == '[':
		val, err = s.scanArray()

	case c == 't':
		ret := s.scanExact([]byte{'r', 'u', 'e'})
		if ret {
			val = true
		} else {
			err = newError(invalidBoolValue)
		}
	case c == 'f':
		ret := s.scanExact([]byte{'a', 'l', 's', 'e'})
		if ret {
			val = false
		} else {
			err = newError(invalidBoolValue)
		}
	case c == 'n':
		ret := s.scanExact([]byte{'u', 'l', 'l'})
		if ret {
			val = nil
		} else {
			err = newError(invalidNullValue)
		}
	case c == '{':
		val, err = s.scanObject()
	default:
		err = newError(invalidValue)
	}
	return
}

func (s *scanner) scanArray() (val interface{}, err error) {
	s.offset++
	v, err := s.scanValue()
	if err != nil {
		return
	}

	switch v0 := v.(type) {
	case string:
		arr := []string{v0}

		state := s.skipRestOfArrayItem()
		for {
			switch state {
			case endOfItem:
				s.skip()
				v, err := s.scanString()
				if err != nil {
					goto ERR
				} else {
					arr = append(arr, v)
					state = s.skipRestOfArrayItem()
				}
			case endOfArray:
				return arr, nil
			default:
				return nil, newError(invalidArray)
			}
		}
	case int:
		arr := []int{v0}

		state := s.skipRestOfArrayItem()
		for {
			switch state {
			case endOfItem:
				raw := s.scanRaw()
				v, err := strconv.Atoi(raw)
				if err != nil {
					goto ERR
				} else {
					arr = append(arr, v)
					state = s.skipRestOfArrayItem()
				}
			case endOfArray:
				return arr, nil
			default:
				return nil, newError(invalidArray)
			}
		}
	case bool:
		arr := []bool{v0}

		state := s.skipRestOfArrayItem()
		for {
			switch state {
			case endOfItem:
				v, err := s.scanBool()
				if err != nil {
					goto ERR
				} else {
					arr = append(arr, v)
					state = s.skipRestOfArrayItem()
				}
			case endOfArray:
				return arr, nil
			default:
				return nil, newError(invalidArray)
			}
		}
	case float64:
		arr := []float64{v0}

		for state := s.skipRestOfArrayItem(); state == endOfItem; {
			switch state {
			case endOfItem:
				raw := s.scanRaw()
				v, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					goto ERR
				} else {
					arr = append(arr, v)
					state = s.skipRestOfArrayItem()
				}
			case endOfArray:
				return arr, nil
			default:
				return nil, newError(invalidArray)
			}
		}
	case time.Time:
		arr := []time.Time{v0}

		for state := s.skipRestOfArrayItem(); state == endOfItem; {
			switch state {
			case endOfItem:
				raw := s.scanRaw()
				v, err := decodeDatetime(raw)
				if err != nil {
					goto ERR
				} else {
					arr = append(arr, v)
					state = s.skipRestOfArrayItem()
				}
			case endOfArray:
				return arr, nil
			default:
				return nil, newError(invalidArray)
			}
		}
	case *Node:
		arr := []*Node{v0}
		for state := s.skipRestOfArrayItem(); state == endOfItem; {
			switch state {
			case endOfItem:
				v, err := s.scanObject()
				if err != nil {
					goto ERR
				} else {
					arr = append(arr, v)
					state = s.skipRestOfArrayItem()
				}
			case endOfArray:
				return arr, nil
			default:
				return nil, newError(invalidArray)
			}
		}
	}

ERR:
	s.skipUntilChar(']')
	s.offset++

	return
}

func (s *scanner) scanRaw() (val string) {
	start := s.offset
	for ; s.offset < s.len; s.offset++ {
		c := s.data[s.offset]
		if isLineEnd(c) || c == ',' || c == ']' || c == '}' || s.isComment() {
			val = string(s.data[start:s.offset])
			break
		}
	}

	if s.offset == s.len {
		val = string(s.data[start:])
	}

	return strings.TrimSpace(val)
}

func (s *scanner) scanString() (val string, err error) {
	c := s.data[s.offset]
	switch c {
	case '"':
		return s.scanQuotedString()
	case '`':
		return s.scanRawString()
	default:
		return "", newError(invalidStringValue)
	}
}

func (s *scanner) scanQuotedString() (val string, err error) {
	s.offset++
	var ret []byte
	for i := s.offset; i < s.len; {
		switch c := s.data[i]; {
		case c == '"':
			if ret != nil {
				val = string(ret)
			} else {
				val = string(s.data[s.offset:i])
			}

			s.offset = i + 1
			return
		case c == '\\':
			//escape
			//i++
			if i == s.len-1 {
				return "", newError(invalidEscape)
			}
			if ret == nil {
				ret = s.data[s.offset:i]
			}

			switch s.data[i+1] {
			case '"', '\\', '/', '\'':
				ret = append(ret, s.data[i])
				i += 2
			case 'b':
				ret = append(ret, '\b')
				i += 2
			case 'f':
				ret = append(ret, '\f')
				i += 2
			case 'n':
				ret = append(ret, '\n')
				i += 2
			case 'r':
				ret = append(ret, '\r')
				i += 2
			case 't':
				ret = append(ret, '\t')
				i += 2
			case 'u':
				i += 2
				j := i + 4
				if j >= s.len {
					return "", newError(invalidUTF8StringValue)
				}

				ub, size, err := escapeU4(s.data[i:j])
				if err != nil {
					return "", err
				}

				ret = append(ret, ub[0:size]...)
				i = j
			default:
				return "", newError(invalidEscape)
			}
		case c < utf8.RuneSelf:
			// ASCII
			if ret != nil {
				ret = append(ret, c)
			}
			i++

			// Coerce to well-formed UTF-8.
		default:
			r, size := utf8.DecodeRune(s.data[i:])
			if r == utf8.RuneError || size == 1 {
				return "", newError(invalidUTF8StringValue)
			} else {
				j := i + size
				if ret != nil {
					ret = append(ret, s.data[i:j]...)
				}
				i = j
			}
		}
	}

	return "", newError(invalidStringValue)
}

func (s *scanner) scanRawString() (val string, err error) {
	s.offset++
	val, err = s.scanUntilChar('`')
	if err == nil {
		s.offset++
	} else {
		err = newError(invalidStringValue)
	}

	return
}

func (s *scanner) scanBool() (val bool, err error) {
	raw := s.scanRaw()
	if raw == "true" {
		val = true
	} else if raw != "false" {
		err = newError(invalidBoolValue)
	}
	return
}

func (s *scanner) scanObject() (val *Node, err error) {
	s.offset++ //skip '{'
	val = NewNode()

	for s.offset < s.len-1 {
		s.skip()
		if s.data[s.offset] == '}' {
			return
		}
		s.scanPair(val)
	}

	return val, newError(invalidObject)
}

func (s *scanner) scanNode(parent *Node) {
	s.offset++
	name := s.scanName()
	if name == "" {
		s.addErrorMsg(invalidNodeName)
		s.skipRestOfLine()
		return
	}

	s.skip()

	if s.data[s.offset] == '-' {
		parent.dict[name] = s.scanNodeList()
	} else {
		parent.dict[name] = s.scanSingleNode()
	}
}

func (s *scanner) scanLine(parent *Node) {
	s.skipSpace()
	s.scanPair(parent)
	s.skipRestOfLine()
}

func (s *scanner) scanSingleNode() (node *Node) {
	node = NewNode()
	for !s.isBlankLine() {
		s.scanLine(node)
	}

	return
}

func (s *scanner) scanNodeList() (list []*Node) {
	list = []*Node{}

	var node *Node
	for !s.isBlankLine() {
		if s.data[s.offset] == '-' {
			node = NewNode()
			list = append(list, node)
			s.offset++
		}
		s.scanLine(node)
	}
	return
}

// skip to next meaningful byte
func (s *scanner) skip() {
	for s.offset < s.len-1 {
		c := s.data[s.offset]
		if s.isComment() || isLineEnd(c) {
			s.skipRestOfLine()
		} else if isSpace(c) {
			s.offset++
		} else {
			return
		}
	}
	/*for i := s.offset; i < s.len; {
		if s.isComment() || isLineEnd(s.data[i]) {
			s.skipRestOfLine()
			i = s.offset
		} else if isSpace(s.data[i]) {
			s.offset++
			i++
		} else {
			return
		}
	}*/
}

func (s *scanner) skipLineEnd() {
	if s.offset >= s.len-1 {
		return
	}

	s.line++
	if s.data[s.offset] == '\r' {
		i := s.offset + 1
		if i < s.len && s.data[i] == '\n' {
			s.offset += 2
			return
		}
	}
	s.offset++
}

func (s *scanner) skipRestOfLine() {
	s.skipUntil(isLineEnd)
	s.skipLineEnd()
}

func (s *scanner) skipSpace() {
	for i := s.offset; i < s.len; i++ {
		if isSpace(s.data[i]) {
			s.offset++
		} else {
			return
		}
	}
}

func (s *scanner) skipRestOfArrayItem() scanState {
	for i := s.offset; i < s.len; {
		c := s.data[i]
		switch {
		case c == ',':
			s.offset = i + 1
			return endOfItem
		case isSpace(c):
			i++
			// continue
		case c == ']':
			s.offset = i + 1
			return endOfArray
		case s.isComment() || isLineEnd(c):
			s.skipRestOfLine()
			i = s.offset
			// continue
		default:
			return scanError
		}
	}
	return eof
}

func (s *scanner) skipUntil(fn func(byte) bool) {
	for s.offset < s.len-1 {
		if fn(s.data[s.offset]) {
			return
		}
		s.offset++
	}
}

func (s *scanner) skipUntilChar(c byte) {
	for i := s.offset; i < s.len; i++ {
		if c == s.data[i] {
			s.offset = i
			return
		}
	}
}

func (s *scanner) findPosOf(c byte) int {
	for i := s.offset; i < s.len; i++ {
		if s.data[i] == c {
			return i
		}
	}
	return -1
}

func (s *scanner) scanUntilChar(c byte) (v string, err error) {
	for i := s.offset; i < s.len; i++ {
		if s.data[i] == c {
			v = string(s.data[s.offset:i])
			s.offset = i
			return
		}
	}

	return "", newError(invalidValue)
}

func (s *scanner) scanLineUntilChar(c byte) (v string, err error) {
	for i := s.offset; i < s.len; i++ {
		if s.data[i] == c {
			v = string(s.data[s.offset:i])
			s.offset = i
			return
		}

		if isLineEnd(s.data[i]) {
			goto ERR
		}
	}

ERR:
	err = errors.New(invalidValue)
	return
}

func (s *scanner) scanExact(expect []byte) bool {
	l := len(expect)
	for i := 0; i < l; i++ {
		if s.offset+i > s.len {
			return false
		}
		if s.data[i] != s.data[s.offset+i+1] {
			return false
		}
	}

	s.offset += l + 1
	return true
}

// RJ support both '#' and `//` to start a comment
func isComment(a, b byte) bool {
	if a == '#' {
		return true
	}

	if a == '/' && b == '/' {
		return true
	}

	return false
}

func (s *scanner) isComment() bool {
	if s.data[s.offset] == '#' {
		return true
	}

	if s.data[s.offset] == '/' {
		next := s.offset + 1
		if next < s.len && s.data[next] == '/' {
			return true
		}
	}

	return false
}

//A node end by blank lines.
//A blank line means it contains only spaces or comments
func (s *scanner) isBlankLine() bool {
	for i := s.offset; i < s.len; i++ {
		c := s.data[i]
		if s.isComment() || isLineEnd(c) {
			s.skipRestOfLine()
			return true
		}

		if !isSpace(c) {
			return false
		}
	}
	return true
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t'
}

func isLineEnd(c byte) bool {
	return c == '\r' || c == '\n'
}
