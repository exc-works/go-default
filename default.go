package go_default

import (
	"encoding"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotPointer = errors.New("input must be a pointer to a struct")
)

// DefaultSetter set the default value for a field
//
//   - path is the full path of the field, like "foo.bar.baz"
//   - fieldValue is the reflect.Value of the field
//   - value is the default value from the tag
type DefaultSetter func(path string, fieldValue reflect.Value, value string) (set bool, err error)

// DurationSetter set the default value for time.Duration
func DurationSetter(path string, fieldValue reflect.Value, value string) (set bool, err error) {
	if fieldValue.Type() != reflect.TypeOf(time.Duration(0)) {
		return false, nil
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return false, fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, fieldValue.Type().String())
	}
	fieldValue.Set(reflect.ValueOf(d))
	return true, nil
}

// TimeSetter set the default value for time.Time
//
// The default layout is time.RFC3339, you can specify a custom layout by separating the value with a semicolon.
// For example, "Fri, 10 Jan 2025 17:20:00 UTC;Mon, 02 Jan 2006 15:04:05 MST"
func TimeSetter(path string, fieldValue reflect.Value, value string) (set bool, err error) {
	if fieldValue.Type() != reflect.TypeOf(time.Time{}) {
		return false, nil
	}
	if !fieldValue.Interface().(time.Time).IsZero() {
		return true, nil // already set
	}
	values := strings.Split(value, ";")
	if len(values) == 1 {
		values = append(values, time.RFC3339)
	}
	t, err := time.Parse(values[1], values[0])
	if err != nil {
		return false, fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, fieldValue.Type().String())
	}
	fieldValue.Set(reflect.ValueOf(t))
	return true, nil
}

// URLSetter set the default value for *url.URL
func URLSetter(path string, fieldValue reflect.Value, value string) (set bool, err error) {
	if fieldValue.Type() != reflect.TypeOf(&url.URL{}) {
		return false, nil
	}
	if !fieldValue.IsNil() {
		return true, nil // already set
	}
	u, err := url.Parse(value)
	if err != nil {
		return false, fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, fieldValue.Type().String())
	}
	fieldValue.Set(reflect.ValueOf(u))
	return true, nil
}

// ByteSliceSetter set the default value for []byte
//
// The value can be a hex string or base64 string, like "0x1234" or "SGVsbG8="
func ByteSliceSetter(path string, fieldValue reflect.Value, value string) (set bool, err error) {
	if fieldValue.Type() != reflect.TypeOf([]byte{}) {
		return false, nil
	}
	if fieldValue.Len() > 0 {
		return true, nil // already set
	}
	if value == "" {
		return true, nil
	}
	if strings.HasPrefix(value, "0x") {
		b, err := hex.DecodeString(value[2:])
		if err != nil {
			return false, fmt.Errorf("cannot set default value for %s, decode %s to %s failed", path, value, fieldValue.Type().String())
		}
		fieldValue.Set(reflect.ValueOf(b))
	} else {
		b, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return false, fmt.Errorf("cannot set default value for %s, decode %s to %s failed", path, value, fieldValue.Type().String())
		}
		fieldValue.Set(reflect.ValueOf(b))
	}
	return true, nil
}

// TextUnmarshalerSetter set the default value for encoding.TextUnmarshaler
//
// The field must be a pointer to a type that implements encoding.TextUnmarshaler
func TextUnmarshalerSetter(path string, fieldValue reflect.Value, value string) (set bool, err error) {
	switch fieldValue.Type().Kind() {
	case reflect.Pointer:
		if !fieldValue.Type().Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
			return false, nil
		}
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		} else {
			return true, nil // already set
		}
		if err := fieldValue.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value)); err != nil {
			return false, fmt.Errorf("cannot set default value for %s, unmarshal %s failed", path, value)
		}
		return true, nil
	default:
		return false, nil
	}
}

type Config struct {
	TagName string          // default tag name
	Setters []DefaultSetter // default setters to convert string to specific type
}

type Option func(cfg *Config)

// WithTagName set the tag name to search for default value
func WithTagName(tagName string) Option {
	return func(cfg *Config) {
		cfg.TagName = tagName
	}
}

// WithSetters set the default setters to convert string to specific type
func WithSetters(setters ...DefaultSetter) Option {
	return func(cfg *Config) {
		cfg.Setters = setters
	}
}

// Struct set the default value for a struct
func Struct(input any, opts ...Option) error {
	cfg := &Config{
		TagName: "default",
		Setters: []DefaultSetter{
			DurationSetter,
			TimeSetter,
			URLSetter,
			ByteSliceSetter,
			TextUnmarshalerSetter,
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	v := reflect.ValueOf(input)
	t := v.Type()
	if t.Kind() != reflect.Pointer {
		return ErrNotPointer
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return ErrNotPointer
	}

	return fillStruct("", v, cfg)
}

func fillStruct(deepName string, v reflect.Value, cfg *Config) error {
	if v.Type().Kind() != reflect.Struct {

	}

	t := v.Type().Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get(cfg.TagName)
		if tagValue == "" {
			continue
		}
		fieldValue := v.Elem().Field(i)

		path := path(deepName, field.Name)
		if !isDefault(fieldValue) {
			continue
		}

		var (
			set bool
			err error
		)
		for _, setter := range cfg.Setters {
			set, err = setter(path, fieldValue, tagValue)
			if err != nil {
				return fmt.Errorf("cannot set default value for %s, err: %w", field.Name, err)
			}
			if set {
				break
			}
		}
		if set {
			continue
		}

		if field.Type.Kind() == reflect.Pointer {
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(field.Type.Elem())) // create a new instance
			}
			if err := fillStruct(path, fieldValue, cfg); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Struct {
			if err := fillStruct(path, fieldValue.Addr(), cfg); err != nil {
				return err
			}
		} else {
			if err := setDefault(path, fieldValue, tagValue); err != nil {
				return err
			}
		}
	}
	return nil
}

func isDefault(fieldValue reflect.Value) bool {
	switch fieldValue.Type().Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool:
		return fieldValue.IsZero()
	case reflect.Struct, reflect.Pointer:
		return true
	case reflect.Slice, reflect.Map:
		return fieldValue.Len() == 0
	default:
		return true
	}
}

func setDefault(path string, fieldValue reflect.Value, value string) error {
	switch fieldValue.Type().Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, fieldValue.Type().String())
		}
		fieldValue.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, fieldValue.Type().String())
		}
		fieldValue.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, fieldValue.Type().String())
		}
		fieldValue.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, fieldValue.Type().String())
		}
		fieldValue.SetBool(b)
	default:
		return fmt.Errorf("unhandled default value for %s, type %s", path, fieldValue.Type().String())
	}
	return nil
}

func path(deepName, name string) string {
	if deepName == "" {
		return name
	}
	return deepName + "." + name
}
