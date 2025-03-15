package event

import "context"

type funcWrapper[E any] struct {
	handler func(ctx context.Context, event E)
}

func Func[E any](lambda func(ctx context.Context, event E)) Handler[E] {
	return &funcWrapper[E]{
		handler: lambda,
	}
}

func (h *funcWrapper[E]) Handle(ctx context.Context, event E) {
	h.handler(ctx, event)
}
