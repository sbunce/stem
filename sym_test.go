// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.
package stem

import (
	"testing"
)

func TestArray(t *testing.T) {
	st := newsymtab(map[string]interface{}{
		"a": []interface{}{},
	})
	if !st.Array("a").IsValid() {
		t.Fatal("invalid array")
	}
}

func TestIfdef(t *testing.T) {
	st := newsymtab(map[string]interface{}{
		"a": "b",
	})
	if !st.Ifdef("a") {
		t.Fatal("ifdef test failed")
	}
}

func TestIfndef(t *testing.T) {
	st := newsymtab(map[string]interface{}{
		"a": "b",
	})
	if st.Ifndef("a") {
		t.Fatal("ifndef test failed")
	}
}

func TestObject(t *testing.T) {
	st := newsymtab(map[string]interface{}{
		"a": map[string]interface{}{},
	})
	if !st.Object("a").IsValid() {
		t.Fatal("'a' is not a valid object")
	}
}

func TestPrint(t *testing.T) {
	st := newsymtab(map[string]interface{}{
		"a": "b",
	})
	if got, want := st.Print("a"), "b"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
