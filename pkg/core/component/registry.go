package component

import (
	"errors"
	"google.golang.org/protobuf/proto"
	"log"
	"maps"
	"reflect"
	"slices"
	"sync"
)

var ErrMarshalerRequired = errors.New("marshaler is required")
var ErrComponentAlreadyRegistered = errors.New("component already registered")
var ErrProtoMessageAlreadyRegistered = errors.New("proto message already registered")

var metadataMut sync.RWMutex
var metadataMapC = make(map[reflect.Type]MetadataAny)
var metadataMapP = make(map[reflect.Type]MetadataAny)

func Register[T any, P proto.Message](marshaler Marshaler[T, P], options ...func(*metadataConfig)) error {

	// validate input and create metadata
	if marshaler == nil {
		return ErrMarshalerRequired
	}
	meta := &Metadata[T, P]{
		marshaler: marshaler,
		config: metadataConfig{
			replicate: true,
			hidden:    false,
		},
	}
	for _, option := range options {
		option(&meta.config)
	}

	// store metadata in map if it does not exist yet
	metadataMut.Lock()
	defer metadataMut.Unlock()

	typC := reflect.TypeFor[T]()
	_, existsC := metadataMapC[typC]
	if existsC {
		return ErrComponentAlreadyRegistered
	}

	typP := reflect.TypeFor[P]()
	_, existsP := metadataMapP[typP]
	if existsP {
		return ErrProtoMessageAlreadyRegistered
	}

	metadataMapC[typC] = meta
	metadataMapP[typP] = meta
	return nil
}

func RegisterFatal[T any, P proto.Message](marshaler Marshaler[T, P], options ...func(*metadataConfig)) {
	err := Register(marshaler, options...)
	if err != nil {
		var empty T
		log.Fatalf("component registration failed for %T: %v", empty, err)
	}
}

func WithHidden[T any](hidden bool) func(config *metadataConfig) {
	return func(config *metadataConfig) {
		config.hidden = hidden
	}
}

func WithReplicate[T any](replicate bool) func(config *metadataConfig) {
	return func(config *metadataConfig) {
		config.replicate = replicate
	}
}

func GetMetadata(component any) (MetadataAny, bool) {
	typ := reflect.TypeOf(component)
	return GetMetadataType(typ)
}

func GetMetadataType(typ reflect.Type) (MetadataAny, bool) {
	metadataMut.RLock()
	defer metadataMut.RUnlock()

	meta, exists := metadataMapC[typ]
	return meta, exists
}

func ListMetadata() []MetadataAny {
	metadataMut.RLock()
	defer metadataMut.RUnlock()

	return slices.Collect(maps.Values(metadataMapC))

}
