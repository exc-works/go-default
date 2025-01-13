package go_default

import (
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
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

	// Type implemented encoding.TextUnmarshaler
	BigInt    big.Int    `default:"1234567890"` // Not a pointer can not be set
	BigIntPtr *big.Int   `default:"1234567890987654321"`
	BigFloat  *big.Float `default:"1.234"`

	Nested       Nested `default:"dive"` // Use dive to indicate dive into the nested struct
	NestedIgnore Nested // not add default tag will be ignored

	NestedPtr *Nested `default:"dive"`
}

type Nested struct {
	String string `default:"world"`
}

func TestStruct(t *testing.T) {
	foo := &Foo{}

	err := Struct(foo)
	require.NoError(t, err)

	require.EqualValues(t, "hello", foo.String)
	require.EqualValues(t, 1, foo.Int)
	require.EqualValues(t, 0, foo.IntWithoutTag)
	require.EqualValues(t, 123, foo.Uint)
	require.EqualValues(t, 2.0, foo.Float32)
	require.EqualValues(t, 1.0, foo.Float64)
	require.EqualValues(t, true, foo.Bool)
	require.EqualValues(t, time.Second, foo.Duration)
	require.EqualValues(t, time.Date(2025, 1, 10, 17, 20, 0, 0, time.UTC), foo.Time)
	require.EqualValues(t, time.Date(2025, 1, 10, 17, 20, 0, 0, time.UTC), foo.TimeWithLayout)
	require.EqualValues(t, big.NewInt(0).String(), foo.BigInt.String())
	require.EqualValues(t, "1234567890987654321", foo.BigIntPtr.String())
	require.EqualValues(t, "1.234", foo.BigFloat.String())
	require.EqualValues(t, "world", foo.Nested.String)
	require.EqualValues(t, "", foo.NestedIgnore.String)
	require.EqualValues(t, "world", foo.NestedPtr.String)
}
