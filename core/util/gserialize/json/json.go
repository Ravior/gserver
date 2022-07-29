package json

import (
	jsoniter "github.com/json-iterator/go"
)

type Serializer struct{}

// NewSerializer returns a new Serializer.
func NewSerializer() *Serializer {
	return &Serializer{}
}

// Marshal returns the JSON encoding of v.
func (s *Serializer) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

// Unmarshal parses the JSON-encoded game_config and stores the result
// in the value pointed to by v.
func (s *Serializer) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}
