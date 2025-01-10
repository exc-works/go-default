package go_default

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Foo struct {
	String string `default:"hello"`

	Int           int `default:"1"`
	IntWithoutTag int

	Uint uint `default:"123"`

	Float32 float32 `default:"2.0"`
	Float64 float64 `default:"1.0"`

	Bool bool `default:"true"`
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
}
