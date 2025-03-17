package server

import (
	"net/http"
)

type HandlerBuilder interface {
	ServeMux() *http.ServeMux
	ReflectNames() *ReflectNames
}

type HandlerService interface {
	RegisterServerHandlers(builder HandlerBuilder) error
}

type ReflectNames struct {
	names []string
}

func (s *ReflectNames) Names() []string {
	return s.names
}

func (s *ReflectNames) Add(name string) {
	s.names = append(s.names, name)
}
