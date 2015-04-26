// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package stem

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

// Set of templates which can include eachother.
type Set struct {
	rwm   sync.RWMutex
	cache map[string]*Template

	// Design Note:
	// A Set is meant to be used by multiple goroutines concurrently. To minimize
	// locking we assume the template will never be modified after it is added.
	// For this reason we don't have a exported function to get a pointer to a
	// template in the Set.
}

// Create new set.
func NewSet() *Set {
	return &Set{
		cache: make(map[string]*Template),
	}
}

// Add template to set or replace existing template. Once a template is added it
// must never be used outside the set because it wouldn't be threadsafe.
func (s *Set) Add(t *Template) {
	s.rwm.Lock()
	defer s.rwm.Unlock()
	s.cache[t.name] = t
}

// Del template from set.
func (s *Set) Del(name string) {
	s.rwm.Lock()
	defer s.rwm.Unlock()
	delete(s.cache, name)
}

// template returns a template in the Set. We do not export this func because it
// would not be threadsafe.
func (s *Set) template(name string) *Template {
	s.rwm.Lock()
	defer s.rwm.Unlock()
	if t, ok := s.cache[name]; ok {
		return t
	}
	return nil
}

// Execute template with specified data.
func (s *Set) Execute(wr io.Writer, name string, data map[string]interface{}) error {
	t := s.template(name)
	if t == nil {
		return fmt.Errorf("template %q not found", name)
	}
	return executeRecurse(wr, s, newsymtab(data), t.tree)
}

// ExecuteJSON executes template with specified JSON data.
func (s *Set) ExecuteJSON(wr io.Writer, name, JSON string) error {
	t := s.template(name)
	if t == nil {
		return fmt.Errorf("template %q not found", name)
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(JSON), data); err != nil {
		return fmt.Errorf("couldn't unmarshal json, error: %v", err)
	}
	return executeRecurse(wr, s, newsymtab(data), t.tree)
}
