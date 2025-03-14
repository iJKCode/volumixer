package valset

type Set[T comparable] map[T]struct{}

func Make[T comparable]() Set[T] {
	return make(map[T]struct{})
}

func (s Set[T]) Has(val T) bool {
	_, ok := s[val]
	return ok
}

func (s Set[T]) Put(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Del(value T) {
	delete(s, value)
}
