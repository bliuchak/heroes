// Package json is a simple wrapper of package json from stdlib.
// This wrapper provides ability to mock methods like Marshal and Unmarshal
package json

import "encoding/json"

// Marshaler interface for JSON marshal method
type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
}

// JSON container for methods which are handling JSON
type JSON struct{}

// Marshal is a simply wrapper of json.Marshal()
// provide ability to create a mock over Marshal(v) method
func (j *JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
