package rj

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

var (
	errValueNotFound      = errors.New("value not found")
	errValueNotAssignable = errors.New("v is not assignable")
	errTypeMismatch       = errors.New("type mismatch")
	errNoName             = errors.New("no name provided")
)

// Node is the represent of a RJ Doc.
type Node struct {
	dict map[string]interface{}
}

// NewNode creates an empty RJ Node
func NewNode() *Node {
	return &Node{dict: make(map[string]interface{})}
}

// Get gets the value of the input name.
// It will return an error if there is anything wrong.
func (n *Node) Get(name string) (val interface{}, err error) {
	finalName, finalNode, err := getFinalNameAndNode(name, n)

	if err != nil {
		return
	}

	if finalNode == nil || finalNode.dict == nil {
		err = errValueNotFound
		return
	}

	val, ok := finalNode.dict[finalName]
	if !ok {
		err = errValueNotFound
	}
	return
}

// GetString gets a string value of the input name.
// It will return an empty string if there is any error.
func (n *Node) GetString(name string) string {
	return n.GetStringOr(name, "")
}

// GetStringOr gets a string value of the input name.
// It will return the input default value if there is any error.
func (n *Node) GetStringOr(name, defaultVal string) string {
	val, err := n.Get(name)

	if err == nil {
		switch v := val.(type) {
		case string:
			return v
		}
	}

	return defaultVal
}

// GetStringOrError gets a string value from the node.
// It will return the error directly if there is an error, return it.
func (n *Node) GetStringOrError(name string) (val string, err error) {
	v, err := n.Get(name)

	if err != nil {
		return
	}

	switch val := v.(type) {
	case string:
		return val, nil
	default:
		return "", errTypeMismatch
	}
}

// GetInt gets a string value of the input name.
// It will return 0 if there is any error.
func (n *Node) GetInt(name string) int {
	return n.GetIntOr(name, 0)
}

// GetIntOr gets a string value of the input name.
// It will return the input default value if there is any error.
func (n *Node) GetIntOr(name string, defaultVal int) int {
	val, err := n.Get(name)

	if err == nil {
		switch v := val.(type) {
		case int:
			return v
		}
	}

	return defaultVal
}

// GetIntOrError gets a string value from the node.
// It will return the error directly if there is an error, return it.
func (n *Node) GetIntOrError(name string) (val int, err error) {
	v, err := n.Get(name)

	if err != nil {
		return
	}

	switch val := v.(type) {
	case int:
		return val, nil
	default:
		return 0, errTypeMismatch
	}
}

// GetFloat gets a string value of the input name.
// It will return 0 if there is any error.
func (n *Node) GetFloat(name string) float64 {
	return n.GetFloatOr(name, 0)
}

// GetFloatOr gets a string value of the input name.
// It will return the input default value if there is any error.
func (n *Node) GetFloatOr(name string, defaultVal float64) float64 {
	val, err := n.Get(name)

	if err == nil {
		switch v := val.(type) {
		case float64:
			return v
		}
	}

	return defaultVal
}

// GetFloatOrError gets a string value from the node.
// It will return the error directly if there is an error, return it.
func (n *Node) GetFloatOrError(name string) (val float64, err error) {
	v, err := n.Get(name)

	if err != nil {
		return
	}

	switch val := v.(type) {
	case float64:
		return val, nil
	default:
		return 0, errTypeMismatch
	}
}

// GetBool gets a string value of the input name.
// It will return 0 if there is any error.
func (n *Node) GetBool(name string) bool {
	return n.GetBoolOr(name, false)
}

// GetBoolOr gets a string value of the input name.
// It will return the input default value if there is any error.
func (n *Node) GetBoolOr(name string, defaultVal bool) bool {
	val, err := n.Get(name)

	if err == nil {
		switch v := val.(type) {
		case bool:
			return v
		}
	}

	return defaultVal
}

// GetBoolOrError gets a string value from the node.
// It will return the error directly if there is one.
func (n *Node) GetBoolOrError(name string) (val bool, err error) {
	v, err := n.Get(name)

	if err != nil {
		return
	}

	switch val := v.(type) {
	case bool:
		return val, nil
	default:
		return false, errTypeMismatch
	}
}

// GetTime gets a string value of the input name.
// It will return an empty time if there is any error.
func (n *Node) GetTime(name string) time.Time {
	return n.GetTimeOr(name, time.Time{})
}

// GetTimeOr gets a string value of the input name.
// It will return the input default value if there is any error.
func (n *Node) GetTimeOr(name string, defaultVal time.Time) time.Time {
	val, err := n.Get(name)

	if err == nil {
		switch v := val.(type) {
		case time.Time:
			return v
		}
	}

	return defaultVal
}

// GetTimeOrError gets a string value from the node.
// It will return the error directly if there is an error, return it.
func (n *Node) GetTimeOrError(name string) (val time.Time, err error) {
	v, err := n.Get(name)

	if err != nil {
		return
	}

	switch val := v.(type) {
	case time.Time:
		return val, nil
	default:
		return time.Time{}, errTypeMismatch
	}
}

// GetArrayOrError gets an array from the node.
// It will return error if there is any.
func (n *Node) GetArrayOrError(name string) (array interface{}, err error) {
	val, err := n.Get(name)
	switch arr := val.(type) {
	case []string:
		array = arr
	case []bool:
		array = arr
	case []int:
		array = arr
	case []float64:
		array = arr
	case []time.Time:
		array = arr
	case nil:
		err = errValueNotFound
	default:
		err = errTypeMismatch
	}
	return
}

// GetStringArrayOrError gets an array of string from the node.
func (n *Node) GetStringArrayOrError(name string) (arr []string, err error) {
	val, err := n.Get(name)
	if err != nil {
		return
	}

	switch arr := val.(type) {
	case []string:
		return arr, nil
	default:
		return nil, errTypeMismatch
	}
}

// GetStringArray gets an array of string from the node.
// It returns nil if the name does not exist or other error happens.
func (n *Node) GetStringArray(name string) []string {
	val, err := n.Get(name)
	if err == nil {
		switch arr := val.(type) {
		case []string:
			return arr
		}
	}

	return nil
}

// GetBoolArrayOrError gets an array of bool from the node.
func (n *Node) GetBoolArrayOrError(name string) (arr []bool, err error) {
	val, err := n.Get(name)
	if err != nil {
		return
	}

	switch arr := val.(type) {
	case []bool:
		return arr, nil
	default:
		return nil, errTypeMismatch
	}
}

// GetBoolArray gets an array of bool from the node.
// It returns nil if the key does not exist or other error happens.
func (n *Node) GetBoolArray(name string) []bool {
	val, err := n.Get(name)
	if err == nil {
		switch arr := val.(type) {
		case []bool:
			return arr
		}
	}

	return nil
}

// GetIntArrayOrError gets an array of int from the node.
func (n *Node) GetIntArrayOrError(name string) (arr []int, err error) {
	val, err := n.Get(name)
	if err != nil {
		return
	}

	switch arr := val.(type) {
	case []int:
		return arr, nil
	default:
		return nil, errTypeMismatch
	}
}

// GetIntArray gets an array of int from the node.
// It returns nil if the key does not exist or other error happens.
func (n *Node) GetIntArray(name string) []int {
	val, err := n.Get(name)
	if err == nil {
		switch arr := val.(type) {
		case []int:
			return arr
		}
	}

	return nil
}

// GetFloatArrayOrError gets an array of float from the node.
func (n *Node) GetFloatArrayOrError(name string) (arr []float64, err error) {
	val, err := n.Get(name)
	if err != nil {
		return
	}

	switch arr := val.(type) {
	case []float64:
		return arr, nil
	default:
		return nil, errTypeMismatch
	}
}

// GetFloatArray gets an array of float from the node.
// It returns nil if the key does not exist or other error happens.
func (n *Node) GetFloatArray(name string) []float64 {
	val, err := n.Get(name)
	if err == nil {
		switch arr := val.(type) {
		case []float64:
			return arr
		}
	}

	return nil
}

// GetDatetimeArray gets an array of time.Time from the node.
func (n *Node) GetTimeArrayOrError(name string) (arr []time.Time, err error) {
	val, err := n.Get(name)
	if err != nil {
		return
	}

	switch arr := val.(type) {
	case []time.Time:
		return arr, nil
	default:
		return nil, errTypeMismatch
	}
}

// GetDatetimeArray gets an array of time.Time from the node.
// It returns nil if the key does not exist or other error happens.
func (n *Node) GetTimeArray(name string) []time.Time {
	val, err := n.Get(name)
	if err == nil {
		switch arr := val.(type) {
		case []time.Time:
			return arr
		}
	}

	return nil
}

// GetStruct get a sub struct from the node
func (n *Node) GetStruct(name string, v interface{}) (err error) {
	node, err := n.GetNode(name)
	if err == nil {
		err = decode(node, v)
	}

	return
}

// GetNode gets a sub-node from the node
func (n *Node) GetNode(name string) (node *Node, err error) {
	val, err := n.Get(name)
	if err != nil {
		return
	}

	switch vt := val.(type) {
	case *Node:
		return vt, nil
	}

	return nil, errTypeMismatch
}

// GetStructList get a list of struct from the node
func (n *Node) GetStructList(name string, v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array:
		println("interface")
	case reflect.Slice:
		println("interface")
	case reflect.Interface:
		println("interface")
	}
	list, err := n.GetNodeList(name)
	if err == nil {
		return
	}

	l := len(list)
	vt := reflect.TypeOf(v)
	array := reflect.MakeSlice(vt, l, l)

	for i := 0; i < l; i++ {
		el := decodeStructField(list[i], vt.Elem())
		array.Index(i).Set(el)
	}

	v = array.Interface()

	return
}

// GetNodeList gets a node list from the node
func (n *Node) GetNodeList(name string) (list []*Node, err error) {
	val, err := n.Get(name)
	if err != nil {
		return
	}

	switch vt := val.(type) {
	case []*Node:
		return vt, nil
	}

	return nil, errTypeMismatch
}

// ToStruct decode the node itself to a struct
func (n *Node) ToStruct(val interface{}) error {
	return decode(n, val)
}

func getFinalNameAndNode(name string, node *Node) (finalName string, finalNode *Node, err error) {
	if len(name) == 0 {
		err = errNoName
		return
	}

	if node == nil {
		err = errValueNotFound
		return
	}

	names := strings.SplitN(name, ".", 2)
	//keys[0] is always there
	if len(names[0]) == 0 {
		err = errNoName
		return
	}
	if len(names) == 1 {
		finalName, finalNode = name, node
	} else {
		if len(names[1]) == 0 {
			err = errNoName
			return
		}
		switch v := node.dict[names[0]].(type) {
		case *Node:
			finalName, finalNode, err = getFinalNameAndNode(names[1], v)
		default:
			err = errValueNotFound
		}
	}

	return
}
