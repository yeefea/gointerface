# gointerface

`gointerface` is an interface extractor for Go.

## Installation

To install `gointerface`, just use the following command:

```bash
go install github.com/yeefea/gointerface@latest
```

## Usage

```
gointerface
  -i string
        input file
  -o string
        output file
```

For example, let's create a file `example.go`.

```go
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
// more comments...
func (x *SomeStruct) ToJson(sb strings.Builder) ([]byte, error) {
	return j.Marshal([]string{"1", "2", "3"})
}

// ComplexMethod is a complex method.
func (s *SomeStruct) ComplexMethod(c, d *map[string]int) (f1, f2 func(
	*map[string]int, *map[string]int), x *map[string]int, y *map[string]int) {
	return nil, nil, nil, nil
}
```

The interface will be extracted from the source code. The package statement, the `import` statements and the comments are preserved after the extraction.

```go
// Code generated from gointerface
package example

import (
        j "encoding/json"
        "strings"
)

type ISomeStruct interface {

        /*

           Some block comments.

        */

        // ToJson converts the object to a json bytes.
        // more comments...
        ToJson(sb strings.Builder) ([]byte, error)

        // ComplexMethod is a complex method.
        ComplexMethod(c, d *map[string]int) (f1, f2 func(
                *map[string]int, *map[string]int), x *map[string]int, y *map[string]int)
}
```

