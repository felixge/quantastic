// Package testuil contains helper methods for testing.
package testutil

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
)

// DeepEqual performs a deep comparison between got and expected, and returns
// an error if they are different. It's a placeholder function until I find or
// create a pkg that performs this task and returns a more detailed diff.
func DeepEqual(got, expected interface{}) error {
	if reflect.DeepEqual(got, expected) {
		return nil
	}
	return spew.Errorf("\nGOT:\n\n%s\n\nEXPECTED:\n\n%s\n", got, expected)
}
