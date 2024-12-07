package main

import (
	"errors"
	"sync"
)

// simple stack impl
type stack struct {
	lock sync.Mutex
	s    []any
}

func NewStack() *stack {
	return &stack{sync.Mutex{}, make([]any, 0)}
}

func (s *stack) Push(v any) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) Pop() (any, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Stack Empty")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}
