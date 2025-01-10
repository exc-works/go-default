package go_default

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrNotPointer = errors.New("input must be a pointer to a struct")
)

type Config struct {
	Tag string
}

func Struct(input any) error {
	cfg := &Config{
		Tag: "default",
	}
	_ = cfg

	v := reflect.ValueOf(input)
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
		tag := field.Tag.Get(cfg.Tag)
		if tag == "" {
			continue
		}

		ok, err := isDefault(field, v.Elem().Field(i))
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			// TODO: support nested struct, maybe use custom setter
			continue
		}
		if err := setDefault(field, v.Elem().Field(i), tag); err != nil {
			return err
		}
	}
	return nil
}

func isDefault(field reflect.StructField, fieldValue reflect.Value) (bool, error) {
	switch field.Type.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool:
		return fieldValue.IsZero(), nil
	case reflect.Struct:
		return true, nil
	default:
		return false, fmt.Errorf("detect unsupported type for %s, type %s", field.Name, field.Type.Name())
	}
}

func setDefault(field reflect.StructField, fieldValue reflect.Value, value string) error {
	switch field.Type.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", field.Name, value, field.Type.Name())
		}
		fieldValue.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", field.Name, value, field.Type.Name())
		}
		fieldValue.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", field.Name, value, field.Type.Name())
		}
		fieldValue.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot set default value for %s, parse %s to %s failed", field.Name, value, field.Type.Name())
		}
		fieldValue.SetBool(b)
	default:
		return fmt.Errorf("unhandled type %s", field.Type.Name())
	}
	return nil
}
