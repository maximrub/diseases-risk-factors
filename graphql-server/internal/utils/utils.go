package utils

import (
	"errors"
	"fmt"
)

// Deref dereference any type to its underlying type, if type is nil the zero value of the type returns
func Deref[V any](v *V) V {
	if v == nil {
		// Create zero value
		var res V
		return res
	}
	return *v
}

// DerefSlice dereference all types in the slice to their underlying type, if type is nil the zero value of the type returns
func DerefSlice[V any](list []*V) []V {
	var items []V

	for _, item := range list {
		var res V
		if item != nil {
			res = *item
		}

		items = append(items, res)
	}

	return items
}

func MapDeref[K comparable, V any](m map[K]*V) map[K]V {
	if m == nil {
		return nil
	}

	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = Deref(v)
	}

	return result
}

// Ptr return a pointer to any type
func Ptr[T any](i T) *T {
	return &i
}

// NewError create a new error
func NewError(format string, args ...any) error {
	return errors.New(fmt.Sprintf(format, args...))
}

// WrapError the given message with provided error
func WrapError(err error, msg string) error {
	return WrapErrorf(err, msg)
}

// WrapErrorf wraps the given message format (with args) with the given error
func WrapErrorf(err error, format string, args ...any) error {
	args = append(args, err)
	return fmt.Errorf(format+": %w", args...)
}

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

// Map maps values from type F to type T
func Map[F any, T any](values []F, f func(F) T) []T {
	if values == nil {
		return nil // Keep nil / empty slice semantics
	}

	res := make([]T, len(values))
	for i := range values {
		res[i] = f(values[i])
	}
	return res
}
