package sreflect

import (
	"fmt"
	"testing"
)

// dumpAny dump any type of variable
func dumpAny(any interface{}) string {
	return fmt.Sprintf("%T : %+v", any, any)
}

// TestGetNameAndType
func TestGetNameAndType(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}
	var testStruct TestStruct
	testStruct.Name = "test"
	t.Log(dumpAny(testStruct))
}
