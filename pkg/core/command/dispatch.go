package command

import (
	"errors"
	"fmt"
	"github.com/ijkcode/volumixer/pkg/core/entity"
)

var ErrCommandNotSupported = errors.New("command not supported")
var ErrEntityNotFound = errors.New("entity not found")
var ErrPanicNotAnError = errors.New("unknown panic")

func DispatchContext(ctx *entity.Context, entity entity.ID, command any) error {
	ent, ok := ctx.Get(entity)
	if !ok {
		return ErrEntityNotFound
	}
	return DispatchEntity(ent, command)
}

func DispatchEntity(ent *entity.Entity, command any) (err error) {
	handler, ok := ent.GetHandler(command)
	if !ok {
		return ErrCommandNotSupported
	}

	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(error)
			if ok {
				err = fmt.Errorf("panic: %w", e)
			}
			e = ErrPanicNotAnError
			err = fmt.Errorf("panic: %w, %v", ErrPanicNotAnError, r)
		}
	}()

	return handler(ent, command)
}
