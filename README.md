# gointerface

`gointerface` is an interface extractor for Go.

## Installation

To install `gointerface`, just use the following command:

```bash
go install github.com/yeefea/gointerface@latest
```

## Usage

```
Usage of gointerface:
  -i string
        Input file or package. By default, the program read from stdin.
  -o string
        Output file. By default, the program writes content to stdout.
  -p    Include private methods.
  -t string
        Specify the types. Multiple types are separated by comma(,). Extract all types if not specified.
```

For example, let's create a file `example.go`.

```go
package example

import (
	"fmt"
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

// ComplexMethod is a complex method.
func (s *SomeStruct) ComplexMethod(c, d *map[string]int) (f1, f2 func(
	*map[string]int, *map[string]int), x *map[string]int, y *map[string]int) {
	return nil, nil, nil, nil
}

// privateMethod is a private method.
func (x *SomeStruct) privateMethod(sb strings.Builder) ([]byte, error) {
	return j.Marshal([]string{"1", "2", "3"})
}

type IntArray []int

// Desc describes the array.
func (arr *IntArray) Desc() {
	fmt.Println(arr)
}
```

The interface will be extracted from the source code, with the `package` statement, the `import` statements and the comments preserved.

```go
// Code generated from gointerface
package example

import (
        j "encoding/json"
        "fmt"
        "strings"
)

type IIntArray interface {

        // Desc describes the array.
        Desc()
}

type ISomeStruct interface {

        // ComplexMethod is a complex method.
        ComplexMethod(c, d *map[string]int) (f1, f2 func(
                *map[string]int, *map[string]int), x *map[string]int, y *map[string]int)

        /*

           Some block comments.

        */

        // ToJson converts the object to a json bytes.
        // some other comments
        ToJson(sb strings.Builder) ([]byte, error)
}
```
