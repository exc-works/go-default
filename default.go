package go_default

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotPointer = errors.New("input must be a pointer to a struct")
)

type DefaultSetter func(field reflect.StructField, fieldValue reflect.Value, value string) (set bool, err error)

func DurationSetter(field reflect.StructField, fieldValue reflect.Value, value string) (set bool, err error) {
	if field.Type != reflect.TypeOf(time.Duration(0)) {
		return false, nil
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return false, fmt.Errorf("cannot set default value for %s, parse %s to %s failed", field.Name, value, field.Type.Name())
	}
	fieldValue.Set(reflect.ValueOf(d))
	return true, nil
}

func TimeSetter(field reflect.StructField, fieldValue reflect.Value, value string) (set bool, err error) {
	if field.Type != reflect.TypeOf(time.Time{}) {
		return false, nil
	}
	values := strings.Split(value, ";")
	if len(values) == 1 {
		values = append(values, time.RFC3339)
	}
	t, err := time.Parse(values[1], values[0])
	if err != nil {
		return false, fmt.Errorf("cannot set default value for %s, parse %s to %s failed", field.Name, value, field.Type.Name())
	}
	fieldValue.Set(reflect.ValueOf(t))
	return true, nil
}

func TextUnmarshalerSetter(field reflect.StructField, fieldValue reflect.Value, value string) (set bool, err error) {
	switch field.Type.Kind() {
	case reflect.Pointer:
		if !field.Type.Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
			return false, nil
		}
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(field.Type.Elem()))
		}
		if err := fieldValue.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value)); err != nil {
			return false, fmt.Errorf("cannot set default value for %s, unmarshal %s failed", field.Name, value)
		}
		return true, nil
	default:
		return false, nil
	}
}

type Config struct {
	TagName string
	Setters []DefaultSetter
}

func Struct(input any) error {
	cfg := &Config{
		TagName: "default",
		Setters: []DefaultSetter{
			DurationSetter,
			TimeSetter,
			TextUnmarshalerSetter,
		},
	}

	return fillStruct("", reflect.ValueOf(input), cfg)
}

func fillStruct(deepName string, v reflect.Value, cfg *Config) error {
	t := v.Type()
	if t.Kind() != reflect.Pointer {
		return ErrNotPointer
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return ErrNotPointer
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get(cfg.TagName)
		if tagValue == "" {
			continue
		}
		fieldValue := v.Elem().Field(i)

		path := path(deepName, field.Name)
		ok, err := isDefault(path, field, fieldValue)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		var set bool
		for _, setter := range cfg.Setters {
			set, err = setter(field, fieldValue, tagValue)
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
			if err := setDefault(path, field, fieldValue, tagValue); err != nil {
				return err
			}
		}
	}
	return nil
}

func isDefault(path string, field reflect.StructField, fieldValue reflect.Value) (bool, error) {
	switch field.Type.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool:
		return fieldValue.IsZero(), nil
	case reflect.Struct, reflect.Pointer:
		return true, nil
	default:
		return false, fmt.Errorf("detect unsupported type for %s, type %s", path, field.Type.Name())
	}
}

func setDefault(path string, field reflect.StructField, fieldValue reflect.Value, value string) error {
	switch field.Type.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, field.Type.Name())
		}
		fieldValue.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, field.Type.Name())
		}
		fieldValue.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, field.Type.Name())
		}
		fieldValue.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", path, value, field.Type.Name())
		}
		fieldValue.SetBool(b)
	default:
		return fmt.Errorf("unhandled default value for %s, type %s", path, field.Type.Name())
	}
	return nil
}

func path(deepName, name string) string {
	if deepName == "" {
		return name
	}
	return deepName + "." + name
}
