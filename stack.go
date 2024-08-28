package main

import (
	"sync"
)

// simple stack impl
type stack struct {
	lock sync.Mutex
	s    []any
}

type EmptyStack struct{}

func (e *EmptyStack) Error() string {
	return "Empty Stack"
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
		return 0, &EmptyStack{}
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}
