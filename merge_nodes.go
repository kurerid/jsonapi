package jsonapi

import (
	"reflect"
	"time"
)

func mergeNodes[T any](base, source T) T {
	baseVal := reflect.ValueOf(base)
	sourceVal := reflect.ValueOf(source)

	if baseVal.Kind() == reflect.Ptr {
		result := mergePointers(baseVal, sourceVal)
		return result.Interface().(T)
	}

	result := mergeValues(baseVal, sourceVal)
	return result.Interface().(T)
}

func mergePointers(basePtr, sourcePtr reflect.Value) reflect.Value {
	switch {
	case basePtr.IsNil() && sourcePtr.IsNil():
		return basePtr
	case basePtr.IsNil():
		return sourcePtr
	case sourcePtr.IsNil():
		return basePtr
	}

	baseValue := basePtr.Elem()
	sourceValue := sourcePtr.Elem()

	var mergedValue reflect.Value
	if baseValue.Kind() == reflect.Struct {
		mergedValue = mergeValues(baseValue, sourceValue)
	} else {
		mergedValue = mergeSimpleValue(baseValue, sourceValue)
	}

	resultPtr := reflect.New(baseValue.Type())
	resultPtr.Elem().Set(mergedValue)
	return resultPtr
}

func mergeValues(baseVal, sourceVal reflect.Value) reflect.Value {
	result := reflect.New(baseVal.Type()).Elem()

	for i := 0; i < baseVal.NumField(); i++ {
		baseField := baseVal.Field(i)
		sourceField := sourceVal.Field(i)
		resultField := result.Field(i)

		mergeField(baseField, sourceField, resultField)
	}

	return result
}

func mergeField(baseField, sourceField, resultField reflect.Value) {
	switch baseField.Kind() {
	case reflect.Ptr:
		mergedPtr := mergePointers(baseField, sourceField)
		resultField.Set(mergedPtr)

	case reflect.Slice:
		mergeSlice(baseField, sourceField, resultField)

	case reflect.Struct:
		mergeStruct(baseField, sourceField, resultField)

	default:
		mergeBasicType(baseField, sourceField, resultField)
	}
}

func mergeSlice(baseSlice, sourceSlice, resultSlice reflect.Value) {
	if !sourceSlice.IsNil() && sourceSlice.Len() > 0 {
		resultSlice.Set(sourceSlice)
	} else {
		resultSlice.Set(baseSlice)
	}
}

func mergeStruct(baseStruct, sourceStruct, resultStruct reflect.Value) {
	if baseStruct.Type() == sourceStruct.Type() {
		mergedStruct := mergeValues(baseStruct, sourceStruct)
		resultStruct.Set(mergedStruct)
	} else {
		mergeBasicType(baseStruct, sourceStruct, resultStruct)
	}
}

func mergeBasicType(baseVal, sourceVal, resultVal reflect.Value) {
	if !isZeroValue(sourceVal) {
		resultVal.Set(sourceVal)
	} else {
		resultVal.Set(baseVal)
	}
}

func mergeSimpleValue(baseVal, sourceVal reflect.Value) reflect.Value {
	if !isZeroValue(sourceVal) {
		return sourceVal
	}
	return baseVal
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		// Особый случай для time.Time
		if t, ok := v.Interface().(time.Time); ok {
			return t.IsZero()
		}
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Interface:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}
