// Package sets

// Copyright © 2021 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sets

import (
	"sort"
)

type String map[string]struct{}

// NewString creates a string set from a list of values.
func NewString(items ...string) String {
	ss := String{}
	ss.Add(items...)
	return ss
}

// Add adds items to the set.
func (s String) Add(items ...string) String {
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Remove removes all items from the set.
func (s String) Remove(items ...string) String {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if item is contained in the set.
func (s String) Has(item string) bool {
	_, ok := s[item]
	return ok
}

// HasAll returns true if all items are contained in the set.
func (s String) HasAll(items ...string) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any of the items are contained in the set.
func (s String) HasAny(items ...string) bool {
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
func (s String) NotIn(set String) String {
	result := NewString()
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
func (s String) Union(set String) String {
	result := NewString()
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
func (s String) Intersection(set String) String {
	result := NewString()
	for item := range s {
		if set.Has(item) {
			result.Add(item)
		}
	}
	return result
}

// Equal returns true if s is equal (as a set) to set.
func (s String) Equal(set String) bool {
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
func (s String) List() []string {
	res := make([]string, 0, len(s))
	for item := range s {
		res = append(res, item)
	}
	return res
}

// SortedList returns the contents as a sorted string slice.
func (s String) SortedList() []string {
	res := s.List()
	sort.Strings(res)
	return res
}

// Len returns the size of the set.
func (s String) Len() int {
	return len(s)
}
