package gserialize

import (
	"github.com/Ravior/gserver/core/util/gserialize/json"
	"github.com/Ravior/gserver/core/util/gserialize/protobuf"
)

var (
	Protobuf = protobuf.NewSerializer()
	Json     = json.NewSerializer()
)

type Serializer interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}
