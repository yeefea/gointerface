# GoInterface

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

### Extract interfaces from a file

For example, let's create a file `example.go` and define a few types as follows:

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

To read from a Go file and print the interfaces to `stdout`, use the following command:

```bash
gointerface -i example.go
```

The interfaces will be extracted from the source code, with the `package` statement, the `import` statements and the comments preserved.

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

### Write interfaces to a file

By default, `gointerface` writes the output to `stdout`. To write the output to a file, use the `-o` option:

```bash
gointerface -i example.go -o interface.go
```

The program will write the output to `interface.go`.


### Extract private methods

Private methods, which start with a lower case letter, are not extracted by default. To extract private methods, use the `-p` option:

```bash
gointerface -i example.go -p
```

Then the private method `privateMethod` will included in the output interface.

```go
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

        // privateMethod is a private method.
        privateMethod(sb strings.Builder) ([]byte, error)
}
```



### Extract interfaces from a package

Let's create a new directory `example` and put the `example.go` file in it. Then create a new file `example2.go` in that directory. Now we have the following directory structure:

```
└── example
    ├── example.go
    └── example2.go
```

Use `-i` option to specify the package:

```bash
gointerface -i example
```

The program will analyze all go files in the `example` directory and extract the interfaces.


### Process both receivers and pointer receivers

A method can have both receivers and pointer receivers. For example, the following `struct`:

```go
// example2.go
package example

type MixedReceiver struct{}

func (r *MixedReceiver) PointerReceiver() {}

func (r MixedReceiver) ValueReceiver() {}
```

If a `struct` have methods with both receivers and pointer receivers, `gointerface` will generate two interfaces for it. One interface corresponds to the pointer receiver methods and is named `I{TypeName}`. The other interface corresponds to the value receiver methods and is named `I{TypeName}Value`. The above example will have the following two interfaces:

```go
// Code generated from gointerface
package example

type IMixedReceiver interface {
        PointerReceiver()
}

type IMixedReceiverValue interface {
        ValueReceiver()
}
```


## License

[MIT Licence](https://github.com/yeefea/gointerface/blob/main/LICENSE)