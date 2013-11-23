/**
 * User: jackong
 * Date: 11/12/13
 * Time: 9:44 AM
 */
package syntax

import "testing"

var slice []int


func TestAppendSlice(t *testing.T) {
	appendSlice(1)
	if len(slice) != 1 {
		t.Error("append one element failed")
	}

	appendSlice(45, 12, 3)
	if len(slice) != 4 {
		t.Error("append multi elements failed")
	}
}

func appendSlice(v ... int) {
	slice = append(slice, v...)
}
