package rj

import (
	"errors"
	"io/ioutil"
	"log"
)

// Load loads a file of given path and parse it into a node
func Load(path string) (node *Node, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	node, err = Parse(bytes)

	return
}

// Parse parses given bytes into a node
func Parse(input []byte) (node *Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			err = errors.New("parse RJ failed")
		}
	}()

	scanner := newScanner(input)
	scanner.scan()

	node = scanner.root
	if scanner.error != nil && scanner.error.msg != "" {
		err = scanner.error
	}
	return
}

// ParseString parses a string to a node
func ParseString(input string) (node *Node, err error) {
	return Parse([]byte(input))
}

// Marshal encodes a value to RJ bytes
func Marshal(v interface{}) []byte {
	e := newEncoder(v)
	return e.encode()
}

// Marshal encodes a value to RJ bytes
func MarshalToFile(v interface{}, filename string) error {
	bytes := Marshal(v)
	return ioutil.WriteFile(filename, bytes, 0644)
}

// Unmarshal decode RJ bytes to struct value
func Unmarshal(data []byte, v interface{}) (err error) {
	var node *Node
	node, err = Parse(data)
	if err != nil {
		return
	}

	/*rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errValueNotAssignable
	}*/
	return decode(node, v)
}

// UnmarshalFile decode a RJ file to struct
func UnmarshalFile(path string, v interface{}) (err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	return Unmarshal(bytes, v)
}
