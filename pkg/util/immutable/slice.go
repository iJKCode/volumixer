package immutable

import (
	"iter"
	"slices"
)

type Slice[T any] struct {
	values []T
}

func SliceFrom[T any](values []T) Slice[T] {
	return Slice[T]{
		values: slices.Clone(values),
	}
}

func (s Slice[T]) Len() int {
	return len(s.values)
}

func (s Slice[T]) Get(idx int) (T, bool) {
	if idx < 0 || idx >= len(s.values) {
		var emptyValue T
		return emptyValue, false
	}
	return s.values[idx], true
}

func (s Slice[T]) Values() iter.Seq[T] {
	return slices.Values(s.values)
}

func (s Slice[T]) Items() iter.Seq2[int, T] {
	return slices.All(s.values)
}

func (s Slice[T]) Slice() []T {
	return slices.Clone(s.values)
}

func SliceEqual[T comparable](a, b Slice[T]) bool {
	return slices.Equal(a.values, b.values)
}
