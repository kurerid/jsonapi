package jsonapi

import (
	"reflect"
	"time"
)

// MergeNodes объединяет две структуры Node, рекурсивно выбирая вариант с большим количеством non-zero значений
func mergeNodes(base, source *Node) *Node {
	result := &Node{}

	// Простые строковые поля
	result.Type = chooseNonEmpty(base.Type, source.Type)
	result.ID = chooseNonEmpty(base.ID, source.ID)
	result.Lid = chooseNonEmpty(base.Lid, source.Lid)
	result.ClientID = chooseNonEmpty(base.ClientID, source.ClientID)

	// Map поля - рекурсивно выбираем тот, где больше non-zero values
	result.Attributes = chooseMapWithMoreData(base.Attributes, source.Attributes)
	result.Relationships = chooseMapWithMoreData(base.Relationships, source.Relationships)

	// Для указателей на Links и Meta
	result.Links = choosePointerMapWithMoreData(base.Links, source.Links)
	result.Meta = choosePointerMapWithMoreData(base.Meta, source.Meta)

	return result
}

// choosePointerMapWithMoreData выбирает указатель на map с бОльшим количеством non-zero значений (рекурсивно)
func choosePointerMapWithMoreData[T ~map[string]interface{}](base, source *T) *T {
	if source == nil && base == nil {
		return nil
	}
	if source == nil {
		return base
	}
	if base == nil {
		return source
	}

	// Рекурсивно считаем non-zero values в каждом
	sourceCount := countNonZeroValuesRecursive(*source)
	baseCount := countNonZeroValuesRecursive(*base)

	if sourceCount > baseCount {
		return source
	} else if baseCount > sourceCount {
		return base
	}

	return source
}

// chooseMapWithMoreData выбирает map с бОльшим количеством non-zero значений (рекурсивно)
func chooseMapWithMoreData(base, source map[string]interface{}) map[string]interface{} {
	if source == nil && base == nil {
		return nil
	}
	if source == nil {
		return base
	}
	if base == nil {
		return source
	}

	// Рекурсивно считаем non-zero values в каждом
	sourceCount := countNonZeroValuesRecursive(source)
	baseCount := countNonZeroValuesRecursive(base)

	if sourceCount > baseCount {
		return source
	} else if baseCount > sourceCount {
		return base
	}

	return source
}

// countNonZeroValuesRecursive рекурсивно подсчитывает количество non-zero значений
func countNonZeroValuesRecursive(data interface{}) int {
	if data == nil {
		return 0
	}

	switch v := data.(type) {
	case map[string]interface{}:
		if len(v) == 0 {
			return 0 // Пустая map = 0 points
		}

		count := 1 // +1 point за сам факт непустой map
		for _, value := range v {
			count += countNonZeroValuesRecursive(value)
		}
		return count

	case []interface{}:
		if len(v) == 0 {
			return 0 // Пустой slice = 0 points
		}

		count := 1 // +1 point за сам факт непустого slice
		for _, item := range v {
			count += countNonZeroValuesRecursive(item)
		}
		return count

	default:
		if !isZeroValueRecursive(v) {
			return 1 // +1 point за non-zero значение
		}
		return 0
	}
}

// isZeroValueRecursive рекурсивно проверяет zero-value
func isZeroValueRecursive(v interface{}) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return val == ""
	case int, int8, int16, int32, int64:
		return val == 0
	case uint, uint8, uint16, uint32, uint64:
		return val == 0
	case float32, float64:
		return val == 0
	case bool:
		return !val
	case time.Time:
		return val.IsZero()
	case map[string]interface{}:
		// Рекурсивно проверяем все значения в map
		for _, value := range val {
			if !isZeroValueRecursive(value) {
				return false
			}
		}
		return true
	case []interface{}:
		// Для пустого slice считаем zero
		return len(val) == 0
	default:
		// Для неизвестных типов используем рефлексию
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Interface:
			return rv.IsNil()
		case reflect.Struct:
			return reflect.DeepEqual(v, reflect.Zero(rv.Type()).Interface())
		default:
			return false
		}
	}
}

// chooseNonEmpty выбирает непустую строку
func chooseNonEmpty(base, source string) string {
	if source != "" {
		return source
	}
	return base
}
