// Package util contains helper functions that haven't found a better home yet.
// Methods from this pkg should always be moved to more appropiate places over
// time.
package util

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
)

// Sdump returns a deep printed string dump of the values passed in.
func Sdump(vals ...interface{}) string {
	return spew.Sdump(vals...)
}

// DeepEqual performs a deep comparison between got and expected, and returns
// an error if they are different. It's a placeholder function until I find or
// create a pkg that performs this task and returns a more detailed diff.
func DeepEqual(got, expected interface{}) error {
	if reflect.DeepEqual(got, expected) {
		return nil
	}
	return spew.Errorf("\nGOT:\n\n%s\n\nEXPECTED:\n\n%s\n", got, expected)
}
