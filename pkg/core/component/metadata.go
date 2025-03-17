package component

import (
	"errors"
	"google.golang.org/protobuf/proto"
)

var ErrMarshalInvalidType = errors.New("invalid component type")

type Marshaler[C any, P proto.Message] interface {
	MarshalProto(cmp C) (P, error)
	UnmarshalProto(msg P) (C, error)
}

type MetadataAny interface {
	MarshalProtoAny(any) (proto.Message, error)
	UnmarshalProtoAny(proto.Message) (any, error)
	ShouldReplicate() bool
	IsHidden() bool
}

type Metadata[C any, P proto.Message] struct {
	marshaler Marshaler[C, P]
	config    metadataConfig
}

type metadataConfig struct {
	replicate bool // disable component network replication
	hidden    bool // do not show component in visualizer
}

func (m Metadata[C, P]) MarshalProto(component C) (P, error) {
	return m.marshaler.MarshalProto(component)
}

func (m Metadata[C, P]) UnmarshalProto(msg P) (C, error) {
	return m.marshaler.UnmarshalProto(msg)
}

func (m Metadata[C, P]) MarshalProtoAny(component any) (proto.Message, error) {
	componentT, ok := component.(C)
	if !ok {
		return nil, ErrMarshalInvalidType
	}
	return m.marshaler.MarshalProto(componentT)
}

func (m Metadata[C, P]) UnmarshalProtoAny(msg proto.Message) (any, error) {
	msgP, ok := msg.(P)
	if !ok {
		var zero C
		return zero, ErrMarshalInvalidType
	}
	return m.marshaler.UnmarshalProto(msgP)
}

func (m Metadata[C, P]) ShouldReplicate() bool {
	return m.config.replicate
}

func (m Metadata[C, P]) IsHidden() bool {
	return m.config.hidden
}
