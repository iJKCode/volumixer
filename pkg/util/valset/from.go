package valset

import "iter"

func From[T comparable](values ...T) Set[T] {
	return FromSlice[T](values)
}

func FromSlice[T comparable](s []T) Set[T] {
	ts := Make[T]()
	for _, v := range s {
		ts.Put(v)
	}
	return ts
}

func FromSeq[T comparable](it iter.Seq[T]) Set[T] {
	ts := Make[T]()
	for v := range it {
		ts.Put(v)
	}
	return ts
}
