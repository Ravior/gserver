package json

import (
	jsonLib "github.com/json-iterator/go"
)

type Serializer struct{}

// NewSerializer returns a new Serializer.
func NewSerializer() *Serializer {
	return &Serializer{}
}

// Marshal returns the JSON encoding of v.
func (s *Serializer) Marshal(v interface{}) ([]byte, error) {
	return jsonLib.Marshal(v)
}

// Unmarshal parses the JSON-encoded game_config and stores the result
// in the value pointed to by v.
func (s *Serializer) Unmarshal(data []byte, v interface{}) error {
	return jsonLib.Unmarshal(data, v)
}
