package go_default

import (
	"github.com/stretchr/testify/require"
	"math/big"
	"net/url"
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
	URL            *url.URL      `default:"https://example.com"`
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
	require.EqualValues(t, "https://example.com", foo.URL.String())
	require.EqualValues(t, []byte{0x12, 0x34}, foo.HexBytes)
	require.EqualValues(t, []byte("Hello"), foo.Base64Bytes)
	require.EqualValues(t, big.NewInt(0).String(), foo.BigInt.String())
	require.EqualValues(t, "1234567890987654321", foo.BigIntPtr.String())
	require.EqualValues(t, "1.234", foo.BigFloat.String())
	require.EqualValues(t, "world", foo.Nested.String)
	require.EqualValues(t, "", foo.NestedIgnore.String)
	require.EqualValues(t, "world", foo.NestedPtr.String)
	require.EqualValues(t, "world", foo.Anonymous.String)
	require.EqualValues(t, "world", foo.AnonymousPtr.String)

	foo = &Foo{
		Time: time.Now(),
	}
	err = Struct(foo)
	require.NoError(t, err)
	require.NotEqualValues(t, time.Date(2025, 1, 10, 17, 20, 0, 0, time.UTC), foo.Time)
}

func TestStruct_String(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "hello", foo.String)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			String: "world",
		}
		err := Struct(foo)
		require.NoError(t, err)
	})
}

func TestStruct_Int(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, 1, foo.Int)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			Int: 2,
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, 2, foo.Int)
	})
	t.Run("should return error when failed to parse int", func(t *testing.T) {
		var foo struct {
			Int int `default:"not int"`
		}
		err := Struct(&foo)
		require.ErrorContains(t, err, "cannot set default value for Int, parse not int to int failed")
	})
}

func TestStruct_IntPtr(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		var foo struct {
			ValuePtr *int `default:"1"`
		}
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, 1, *foo.ValuePtr)
	})
	t.Run("not set", func(t *testing.T) {
		var foo struct {
			ValuePtr *int `default:"1"`
		}
		var i = 10
		foo.ValuePtr = &i
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, 10, *foo.ValuePtr)
	})
}

func TestStruct_Uint(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, 123, foo.Uint)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			Uint: 456,
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, 456, foo.Uint)
	})
}

func TestStruct_UintPtr(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		var foo struct {
			ValuePtr *uint `default:"1"`
		}
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, 1, *foo.ValuePtr)
	})
	t.Run("not set", func(t *testing.T) {
		var foo struct {
			ValuePtr *uint `default:"1"`
		}
		var i uint = 10
		foo.ValuePtr = &i
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, 10, *foo.ValuePtr)
	})
}

func TestStruct_Float32(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, 2.0, foo.Float32)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			Float32: 3.0,
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, 3.0, foo.Float32)
	})
}

func TestStruct_Float64(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, 1.0, foo.Float64)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			Float64: 4.0,
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, 4.0, foo.Float64)
	})
}

func TestStruct_Float64Ptr(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		var foo struct {
			ValuePtr *float64 `default:"1.432"`
		}
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, 1.432, *foo.ValuePtr)
	})
	t.Run("not set", func(t *testing.T) {
		var foo struct {
			ValuePtr *float64 `default:"1"`
		}
		var i = 1.234
		foo.ValuePtr = &i
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, 1.234, *foo.ValuePtr)
	})
}

func TestStruct_Bool(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		var boolT struct {
			Bool bool `default:"false"`
		}
		err := Struct(&boolT)
		require.NoError(t, err)
		require.EqualValues(t, false, boolT.Bool)
	})
	t.Run("not set", func(t *testing.T) {
		var boolT struct {
			Bool bool `default:"false"`
		}
		boolT.Bool = true
		err := Struct(&boolT)
		require.NoError(t, err)
		require.EqualValues(t, true, boolT.Bool)
	})
}

func TestStruct_Duration(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, time.Second, foo.Duration)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			Duration: 2 * time.Second,
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, 2*time.Second, foo.Duration)
	})
}

func TestStruct_DurationPtr(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		var foo struct {
			ValuePtr *time.Duration `default:"1s"`
		}
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, time.Second, *foo.ValuePtr)
	})
	t.Run("not set", func(t *testing.T) {
		var foo struct {
			ValuePtr *time.Duration `default:"1"`
		}
		var i = time.Hour
		foo.ValuePtr = &i
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, time.Hour, *foo.ValuePtr)
	})
}

func TestStruct_Time(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, time.Date(2025, 1, 10, 17, 20, 0, 0, time.UTC), foo.Time)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			Time: time.Now(),
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.NotEqualValues(t, time.Date(2025, 1, 10, 17, 20, 0, 0, time.UTC), foo.Time)
	})
}

func TestStruct_TimePtr(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		var foo struct {
			ValuePtr *time.Time `default:"2025-01-10T17:20:00Z"`
		}
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, time.Date(2025, 1, 10, 17, 20, 0, 0, time.UTC), *foo.ValuePtr)
	})
	t.Run("not set", func(t *testing.T) {
		var foo struct {
			ValuePtr *time.Time `default:"2025-01-10T17:20:00Z"`
		}
		var i = time.Now()
		foo.ValuePtr = &i
		err := Struct(&foo)
		require.NoError(t, err)
		require.EqualValues(t, i, *foo.ValuePtr)
	})
}

func TestStruct_TimeWithLayout(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, time.Date(2025, 1, 10, 17, 20, 0, 0, time.UTC), foo.TimeWithLayout)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			TimeWithLayout: time.Now(),
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.NotEqualValues(t, time.Date(2025, 1, 10, 17, 20, 0, 0, time.UTC), foo.TimeWithLayout)
	})
}

func TestStruct_URL(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "https://example.com", foo.URL.String())
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			URL: &url.URL{Scheme: "https", Host: "github.com"},
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "https://github.com", foo.URL.String())
	})
}

func TestStruct_HexBytes(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, []byte{0x12, 0x34}, foo.HexBytes)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			HexBytes: []byte{0x56, 0x78},
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, []byte{0x56, 0x78}, foo.HexBytes)
	})
}

func TestStruct_Base64Bytes(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, []byte("Hello"), foo.Base64Bytes)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			Base64Bytes: []byte("World"),
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, []byte("World"), foo.Base64Bytes)
	})
}

func TestStruct_BigIntPtr(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "1234567890987654321", foo.BigIntPtr.String())
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			BigIntPtr: big.NewInt(9876543210),
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "9876543210", foo.BigIntPtr.String())
	})
}

func TestStruct_BigFloat(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "1.234", foo.BigFloat.String())
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			BigFloat: big.NewFloat(5.678),
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "5.678", foo.BigFloat.String())
	})
}

func TestStruct_Nested(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "world", foo.Nested.String)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			Nested: Nested{String: "universe"},
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "universe", foo.Nested.String)
	})
}

func TestStruct_NestedPtr(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "world", foo.NestedPtr.String)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			NestedPtr: &Nested{String: "galaxy"},
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "galaxy", foo.NestedPtr.String)
	})
}

func TestStruct_Anonymous(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "world", foo.Anonymous.String)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			Anonymous: Nested{String: "cosmos"},
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "cosmos", foo.Anonymous.String)
	})
}

func TestStruct_AnonymousPtr(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		foo := &Foo{}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "world", foo.AnonymousPtr.String)
	})
	t.Run("not set", func(t *testing.T) {
		foo := &Foo{
			AnonymousPtr: &Nested{String: "multiverse"},
		}
		err := Struct(foo)
		require.NoError(t, err)
		require.EqualValues(t, "multiverse", foo.AnonymousPtr.String)
	})
}

func TestStruct_UnsupportedType(t *testing.T) {
	var foo struct {
		Unsupported chan int `default:"1"`
	}
	err := Struct(&foo)
	require.ErrorContains(t, err, "cannot set default value for Unsupported, no suitable default setter for chan int")
}
