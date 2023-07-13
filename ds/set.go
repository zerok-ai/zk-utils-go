package ds

type Set[T comparable] map[T]bool

func (s Set[T]) Add(key T) Set[T] {
	s[key] = true
	return s
}

func (s Set[T]) AddBulk(keys []T) Set[T] {
	for _, key := range keys {
		s[key] = true
	}
	return s
}

func (s Set[T]) Remove(key T) Set[T] {
	delete(s, key)
	return s
}

func (s Set[T]) Contains(key T) bool {
	_, ok := s[key]
	return ok
}

func (s Set[T]) Equals(other Set[T]) bool {
	if len(s) != len(other) {
		return false
	}
	for key, _ := range s {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	union := make(Set[T])
	for key, _ := range s {
		union.Add(key)
	}
	for key, _ := range other {
		union.Add(key)
	}
	return union
}

func (s Set[T]) Intersection(other Set[T]) Set[T] {
	intersection := make(Set[T])
	for key, _ := range s {
		if other.Contains(key) {
			intersection.Add(key)
		}
	}
	return intersection
}

func (s Set[T]) GetAll() []T {
	keys := make([]T, 0)
	for key, _ := range s {
		keys = append(keys, key)
	}
	return keys
}
