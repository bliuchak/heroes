package json

import "encoding/json"

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
}

type JSON struct{}

func (j *JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
