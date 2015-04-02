// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.
package stem

import (
	"fmt"
	"reflect"
)

func stringEqual(arv, brv reflect.Value) bool {
	arv = indirect(arv)
	brv = indirect(brv)
	switch arv.Kind() {
	case reflect.Map:
		if brv.Kind() != reflect.Map {
			return false
		}
		if arv.IsNil() != brv.IsNil() {
			return false
		}
		if arv.Len() != brv.Len() {
			return false
		}
		for _, key := range arv.MapKeys() {
			if !stringEqual(arv.MapIndex(key), brv.MapIndex(key)) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if brv.Kind() != reflect.Slice {
			return false
		}
		if arv.IsNil() != brv.IsNil() {
			return false
		}
		if arv.Len() != brv.Len() {
			return false
		}
		for i := 0; i < arv.Len(); i++ {
			if !stringEqual(arv.Index(i), brv.Index(i)) {
				return false
			}
		}
		return true
	}
	return fmt.Sprint(arv.Interface()) == fmt.Sprint(brv.Interface())
}

// StringEqual returns true if both maps would render the same in a template.
// All values are converted to strings before comparison.
func StringEqual(a, b map[string]interface{}) bool {
	return stringEqual(reflect.ValueOf(a), reflect.ValueOf(b))
}
