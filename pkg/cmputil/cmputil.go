package cmputil

import (
	"github.com/google/go-cmp/cmp"
)

// Equal reports whether two values are equal according to cmp.Equal.
func Equal(a, b interface{}) bool {
	return cmp.Equal(a, b)
}

// Diff returns a human-friendly diff between two values using cmp.Diff.
func Diff(a, b interface{}) string {
	return cmp.Diff(a, b)
}
