package ds

type Set map[string]bool

func (s Set) Add(key string) Set {
	s[key] = true
	return s
}

func (s Set) AddBulk(keys []string) Set {
	for _, key := range keys {
		s[key] = true
	}
	return s
}

func (s Set) Remove(key string) Set {
	delete(s, key)
	return s
}

func (s Set) Contains(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Equals(other Set) bool {
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

func (s Set) Union(other Set) Set {
	union := make(Set)
	for key, _ := range s {
		union.Add(key)
	}
	for key, _ := range other {
		union.Add(key)
	}
	return union
}

func (s Set) Intersection(other Set) Set {
	intersection := make(Set)
	for key, _ := range s {
		if other.Contains(key) {
			intersection.Add(key)
		}
	}
	return intersection
}

func (s Set) GetAll() []string {
	keys := make([]string, 0)
	for key, _ := range s {
		keys = append(keys, key)
	}
	return keys
}
