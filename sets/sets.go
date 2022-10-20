// Created by zc on 2022/10/20.

package sets

type Set[T comparable] map[T]struct{}

func New[T comparable](items ...T) Set[T] {
	s := Set[T]{}
	s.Add(items...)
	return s
}

// Add adds items to the set.
func (s Set[T]) Add(items ...T) Set[T] {
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Remove removes all items from the set.
func (s Set[T]) Remove(items ...T) Set[T] {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if item is contained in the set.
func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

// HasAll returns true if all items are contained in the set.
func (s Set[T]) HasAll(items ...T) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any of the items are contained in the set.
func (s Set[T]) HasAny(items ...T) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

// NotIn returns a set of objects that are not in the current set.
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.NotIn(s2) = {a3}
// s2.NotIn(s1) = {a4, a5}
func (s Set[T]) NotIn(set Set[T]) Set[T] {
	result := New[T]()
	for item := range s {
		if !set.Has(item) {
			result.Add(item)
		}
	}
	return result
}

// Union returns a new set which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}
func (s Set[T]) Union(set Set[T]) Set[T] {
	result := New[T]()
	for item := range s {
		result.Add(item)
	}
	for item := range set {
		result.Add(item)
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2.
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}
func (s Set[T]) Intersection(set Set[T]) Set[T] {
	result := New[T]()
	for item := range s {
		if set.Has(item) {
			result.Add(item)
		}
	}
	return result
}

// Equal returns true if s is equal (as a set) to set.
func (s Set[T]) Equal(set Set[T]) bool {
	if len(s) != len(set) {
		return false
	}
	for item := range s {
		if !set.Has(item) {
			return false
		}
	}
	return true
}

// List returns the contents as a string slice.
func (s Set[T]) List() []T {
	res := make([]T, 0, len(s))
	for item := range s {
		res = append(res, item)
	}
	return res
}

// Len returns the size of the set.
func (s Set[T]) Len() int {
	return len(s)
}
