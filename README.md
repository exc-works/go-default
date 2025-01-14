# Go Default Value Setter

This project provides a utility to set default values for struct fields in Go using struct tags. It supports various
types including `time.Time`, `*url.URL`, `[]byte`, and more.

## Installation

To install the package, use `go get`:

```sh
go get github.com/exc-works/go-default
```

## Usage

#### Basic Usage

To use the default value setter, define your struct with default tags and call the `Struct` function:

```go
package main

import (
	"fmt"
	"net"
	"time"
	"math/big"
	"net/url"
	godefault "github.com/exc-works/go-default"
)

type Foo struct {
	String string `default:"hello"`

	Int           int `default:"1"`
	IntWithoutTag int

	Uint uint `default:"123"`

	Float32 float32 `default:"2.0"`
	Float64 float64 `default:"1.0"`

	Bool bool `default:"true"`

	// Complex types

	Duration       time.Duration `default:"1s"`
	Time           time.Time     `default:"2025-01-10T17:20:00Z"`                                        // use time.RFC3339 as default layout
	TimeWithLayout time.Time     `default:"Fri, 10 Jan 2025 17:20:00 UTC;Mon, 02 Jan 2006 15:04:05 MST"` // use custom layout time.RFC1123
	URL            *url.URL      `default:"https://example.com"`
	IPV4           *net.IPAddr   `default:"156.33.241.5"`
	IPV6           *net.IPAddr   `default:"2600:1400:a::1743:fa93"`
	HexBytes       []byte        `default:"0x1234"`
	Base64Bytes    []byte        `default:"SGVsbG8="`

	// Type implemented encoding.TextUnmarshaler
	BigInt    big.Int    `default:"1234567890"` // Not a pointer can not be set
	BigIntPtr *big.Int   `default:"1234567890987654321"`
	BigFloat  *big.Float `default:"1.234"`

	Nested       Nested `default:"dive"` // Use dive to indicate dive into the nested struct
	NestedIgnore Nested // not add default tag will be ignored

	NestedPtr *Nested `default:"dive"`

	Anonymous     `default:"dive"`
	*AnonymousPtr `default:"dive"`
}

type Nested struct {
	String string `default:"world"`
}

type Anonymous = Nested
type AnonymousPtr = Nested

func main() {
	var foo Foo
	err := godefault.Struct(&foo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Default Value:", foo.Time)
}
```

#### Supported Types

The following types are supported out of the box:

- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `bool`
- `string`
- `time.Duration`
- `time.Time`

> Note: The default layout for `time.Time` is `time.RFC3339`. To use a custom layout, specify the layout in the tag.
> For example, `default:"Fri, 10 Jan 2025 17:20:00 UTC;Mon, 02 Jan 2006 15:04:05 MST"`.

- `*url.URL`
- `*net.IPAddr`
- `[]byte`
- any type that implements `encoding.TextUnmarshaler`, e.g. `*big.Int`, `*big.Float`

> Note: The pointer types are supported for all the above types.

#### Nested Structs

The default value setter supports nested structs. To set default values for nested structs, use the `dive` tag:

```go
type Foo struct {
Nested Nested `default:"dive"`
}
```

#### Custom Tag Name

You can configure the tag name using options:

```go
package main

import (
	"fmt"
	godefault "github.com/exc-works/go-default"
)

type Foo struct {
	Value string `custom:"default_value"`
}

func main() {
	var foo Foo
	err := godefault.Struct(&foo, godefault.WithTagName("custom"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Value:", foo.Value)
}
```

#### Custom setters

You can define custom setters for additional types. A custom setter is a function that matches the `DefaultSetter` type:

```go
package main

import (
	"fmt"
	"reflect"
	godefault "github.com/exc-works/go-default"
)

func CustomSetter(path string, fieldValue reflect.Value, value string) (bool, error) {
	if fieldValue.Type().Kind() != reflect.String {
		return false, nil
	}
	fieldValue.SetString("custom:" + value)
	return true, nil
}

func main() {
	var foo struct {
		CustomField string `default:"example"`
	}
	err := godefault.Struct(&foo,
		godefault.WithSetters(append(
			godefault.DefaultSetters(),
			func(path string, fieldValue reflect.Value, value string) (set bool, err error) {
				if fieldValue.Type().Kind() != reflect.String {
					return false, nil
				}
				fieldValue.SetString(value + " world")
				return true, nil
			},
		)...,
		),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Custom Field:", foo.CustomField)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
