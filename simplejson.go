package json

import (
	"encoding/json"
	"bytes"
	"errors"
)

// returns the current implementation version
func Version() string {
	return "0.5.0-alpha"
}

type Json struct {
	data interface{}
}

// NewJson returns a pointer to a new `Json` object
// after unmarshaling `body` bytes
func NewJson(body []byte) (*Json, error) {
	j := new(Json)
	err := j.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// Create new empty json
func NewEmptyJson() *Json {
	i, _ := NewJson([]byte("{}"))
	return i
}

func (j *Json) UnmarshalJSON(p []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(p))
	dec.UseNumber()
	return dec.Decode(&j.data)
}

// Implements the json.Marshaler interface.
func (j *Json) MarshalJSON() ([]byte, error) {
	return json.Marshal(&j.data)
}

// Get returns a pointer to a new `Json` object
// for `key` in its `map` representation
//
// useful for chaining operations (to traverse a nested JSON):
//    js.Get("top_level").Get("dict").Get("value").Int()
func (j *Json) Get(key string) *Json {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Json{val}
		}
	}
	return &Json{nil}
}

// Set modifies `Json` map by `key` and `value`
// Useful for changing single key/value in a `Json` object easily.
func (j *Json) Set(key string, val interface{}) {
	m, err := j.Map()
	if err != nil {
		return
	}
	m[key] = val
}

// Remove modifies `Json` map by `key`
func (j *Json) Remove(key string) {
	m, err := j.Map()
	if err != nil {
		return
	}
	delete(m, key)
}

func (j *Json) Has(key string) bool {
	m, err := j.Map()
	if err == nil {
		if _, ok := m[key]; ok {
			return true
		}
	}
	return false
}

func (j *Json) GetString(key string) (string, error) {
	if !j.Has(key) {
		return "", errors.New(key + "not existed")
	}
	return j.Get(key).String()
}

func (j *Json) GetBool(key string) (bool, error) {
	if !j.Has(key) {
		return false, errors.New(key + "not existed")
	}
	return j.Get(key).Bool()
}

func (j *Json) GetArray(key string) ([]interface{}, error) {
	if !j.Has(key) {
		return nil, errors.New(key + "not existed")
	}
	return j.Get(key).Array()
}

func (j *Json) GetMap(key string) (map[string]interface{}, error) {
	if !j.Has(key) {
		return nil, errors.New(key + "not existed")
	}
	return j.Get(key).Map()
}

func (j *Json) GetStringArray(key string) ([]string, error) {
	if !j.Has(key) {
		return nil, errors.New(key + "not existed")
	}
	return j.Get(key).StringArray()
}

// Map type asserts to `map`
func (j *Json) Map() (map[string]interface{}, error) {
	if m, ok := (j.data).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

// Array type asserts to an `array`
func (j *Json) Array() ([]interface{}, error) {
	if a, ok := (j.data).([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

// Bool type asserts to `bool`
func (j *Json) Bool() (bool, error) {
	if s, ok := (j.data).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}

// String type asserts to `string`
func (j *Json) String() (string, error) {
	if s, ok := (j.data).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}

// StringArray type asserts to an `array` of `string`
func (j *Json) StringArray() ([]string, error) {
	arr, err := j.Array()
	if err != nil {
		return nil, err
	}
	retArr := make([]string, 0, len(arr))
	for _, a := range arr {
		s, ok := a.(string)
		if !ok {
			return nil, err
		}
		retArr = append(retArr, s)
	}
	return retArr, nil
}