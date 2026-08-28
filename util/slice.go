package util

import (
	"cmp"
	"reflect"
	"slices"
	"unsafe"
)

func IsNil(v any) bool {
	if v == nil {
		return true
	}
	reflectValue := reflect.ValueOf(v)
	switch reflectValue.Kind() {
	case reflect.Pointer, reflect.Chan, reflect.Func, reflect.Map, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return reflectValue.IsNil()
	}
	return false
}

func StrToSlice(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func UniqueSlice[S ~[]E, E comparable](s S) S {
	if len(s) < 2 {
		return s
	}
	picked := make(map[E]struct{}, len(s))
	ret := make([]E, 0, len(s))
	for _, v := range s {
		if _, ok := picked[v]; ok {
			continue
		}
		picked[v] = struct{}{}
		ret = append(ret, v)
	}
	return ret
}

func EqualSlice[T cmp.Ordered](x, y []T) bool {
	if (x == nil) != (y == nil) {
		return false
	}
	if len(x) != len(y) {
		return false
	}
	copyX := append([]T(nil), x...)
	copyY := append([]T(nil), y...)
	slices.Sort(copyX)
	slices.Sort(copyY)
	for index, value := range copyX {
		if value != copyY[index] {
			return false
		}
	}
	return true
}
