package typeset

import "iter"

func From(values ...any) TypeSet {
	return FromSlice(values)
}

func FromSlice(s []any) TypeSet {
	ts := Make()
	for _, v := range s {
		ts.Put(v)
	}
	return ts
}

func FromSeq(it iter.Seq[any]) TypeSet {
	ts := Make()
	for v := range it {
		ts.Put(v)
	}
	return ts
}

func FromMapKeys[V any](m map[any]V) TypeSet {
	ts := Make()
	for v := range m {
		ts.Put(v)
	}
	return ts
}

func FromMapValues[K comparable](m map[K]any) TypeSet {
	ts := Make()
	for _, v := range m {
		ts.Put(v)
	}
	return ts
}

func FromSeq2Keys[V any](it iter.Seq2[any, V]) TypeSet {
	ts := Make()
	for v, _ := range it {
		ts.Put(v)
	}
	return ts
}

func FromSeq2Values[K any](it iter.Seq2[K, any]) TypeSet {
	ts := Make()
	for _, v := range it {
		ts.Put(v)
	}
	return ts
}
