package component

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"reflect"
)

var ErrUnknownComponent = errors.New("unknown component")
var ErrUnknownProtoMessage = errors.New("unknown proto message")

func MarshalProto(cmp any) (proto.Message, MetadataAny, error) {
	typ := reflect.TypeOf(cmp)
	meta, ok := metadataMapC[typ]
	if !ok {
		return nil, nil, ErrUnknownComponent
	}
	msg, err := meta.MarshalProtoAny(cmp)
	return msg, meta, err
}

func UnmarshalProto(msg proto.Message) (any, MetadataAny, error) {
	typ := reflect.TypeOf(msg)
	meta, ok := metadataMapP[typ]
	if !ok {
		return nil, nil, ErrUnknownProtoMessage
	}
	cmp, err := meta.UnmarshalProtoAny(msg)
	return cmp, meta, err
}

func MarshalProtoAny(cmp any) (*anypb.Any, MetadataAny, error) {
	pro, meta, err := MarshalProto(cmp)
	if err != nil {
		return nil, nil, fmt.Errorf("marshaling component '%T: %w", cmp, err)
	}
	msg, err := anypb.New(pro)
	if err != nil {
		return nil, nil, fmt.Errorf("marshaling proto any, '%T': %w", msg, err)
	}
	return msg, meta, nil
}

func UnmarshalProtoAny(msg *anypb.Any) (any, MetadataAny, error) {
	pro, err := msg.UnmarshalNew()
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshaling proto any '%s': %w", msg.MessageName(), err)
	}
	cmp, meta, err := UnmarshalProto(pro)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshaling proto message '%T': %w", pro, err)
	}
	return cmp, meta, nil
}
