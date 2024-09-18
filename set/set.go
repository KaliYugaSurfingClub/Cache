package set

import (
	"math/rand"
)

type Set struct {
	elements map[string]int
	indexes  []string
}

func New() *Set {
	return &Set{
		elements: make(map[string]int),
		indexes:  []string{},
	}
}

func (s *Set) Add(item string) {
	if _, exists := s.elements[item]; exists {
		return
	}

	s.indexes = append(s.indexes, item)
	s.elements[item] = len(s.indexes) - 1
}

func (s *Set) Remove(item string) {
	index, exists := s.elements[item]
	if !exists {
		return
	}

	lastIndex := len(s.indexes) - 1
	lastItem := s.indexes[lastIndex]
	s.indexes[index] = lastItem
	s.elements[lastItem] = index

	s.indexes = s.indexes[:lastIndex]
	delete(s.elements, item)
}

func (s *Set) Random() string {
	if len(s.indexes) == 0 {
		return ""
	}

	randomIndex := rand.Intn(len(s.indexes))
	return s.indexes[randomIndex]
}

func (s *Set) Clear() {
	clear(s.elements)
	clear(s.indexes)
}
