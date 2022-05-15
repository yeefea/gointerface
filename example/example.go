package example

import (
	"strings"

	j "encoding/json"
)

type SomeStruct struct{}

/*

Some block comments.

*/

// ToJson converts the object to a json bytes.
// some other comments
func (x *SomeStruct) ToJson(sb strings.Builder) ([]byte, error) {
	return j.Marshal([]string{"1", "2", "3"})
}

// ComplexMethod is a complex method
func (s *SomeStruct) ComplexMethod(c, d *map[string]int) (f1, f2 func(
	*map[string]int, *map[string]int), x *map[string]int, y *map[string]int) {
	return nil, nil, nil, nil
}
