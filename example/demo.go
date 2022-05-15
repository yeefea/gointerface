package example

import (
	"strings"

	j "encoding/json"
)

type A struct{}

/*

abcddd

*/

// ......
// comments
func (x *A) FFF(sb strings.Builder) {
	j.Marshal(123)

}

func (x A) String(sb strings.Builder) {

}

// a complex method
func (s *A) abc(c, d *map[string]int) (f1, f2 func(
	*map[string]int, *map[string]int), x *map[string]int, y *map[string]int) {
	return nil, nil, nil, nil
}
