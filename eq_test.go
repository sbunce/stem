// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package stem

import (
	"testing"
)

func TestStringEqual(t *testing.T) {
	a := map[string]interface{}{
		"num":   123,
		"slice": []byte{0x00, 0x01, 0x02},
		"map":   map[string]interface{}{"abc": 123},
		"float": float64(123),
	}
	b := map[string]interface{}{
		"num":   "123",
		"slice": []byte{0x00, 0x01, 0x02},
		"map":   map[string]interface{}{"abc": "123"},
		"float": 123,
	}
	if !StringEqual(a, b) {
		t.Fatal("string equal failed")
	}
}
