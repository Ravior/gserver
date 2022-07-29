package protobuf

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	googleProto "github.com/golang/protobuf/proto"
)

// ErrWrongValueType is the errors used for marshal the value with protobuf encoding.
var ErrWrongValueType = errors.New("protobuf: convert on wrong type value")

type Serializer struct {
}

// NewSerializer returns a new Serializer.
func NewSerializer() *Serializer {
	return &Serializer{}
}

// Marshal returns the protobuf encoding of v.
func (s *Serializer) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrWrongValueType
	}
	return proto.Marshal(pb)
}

// Unmarshal parses the protobuf-encoded game_config and stores the result
// in the value pointed to by v.
func (s *Serializer) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrWrongValueType
	}
	return proto.Unmarshal(data, pb)
}

func (s *Serializer) GetMessageName(msg proto.Message) string {
	reflect := googleProto.MessageV2(msg).ProtoReflect()
	return string(reflect.Descriptor().FullName())
}
