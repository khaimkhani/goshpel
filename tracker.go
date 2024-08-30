package main

import (
	"errors"
)

type Tracker struct {
	// lock considerations ?
	vars map[string]any
}

func NewTracker() *Tracker {
	return &Tracker{}
}

func (t *Tracker) Add(key string, val any) error {

	if val, ok := t.vars[key]; !ok {
		t.vars[key] = val
		return nil
	}

	return errors.New("Already tracked")
}

func (t *Tracker) Exists(key string) bool {
	_, ok := t.vars[key]
	return ok
}
