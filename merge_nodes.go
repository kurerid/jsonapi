package jsonapi

import (
	"time"
)

// MergeNodes объединяет две структуры Node, выбирая вариант с большим количеством non-zero значений
func mergeNodes(base, source *Node) *Node {
	result := &Node{}

	// Простые строковые поля (приоритет source)
	result.Type = chooseNonEmpty(base.Type, source.Type)
	result.ID = chooseNonEmpty(base.ID, source.ID)
	result.Lid = chooseNonEmpty(base.Lid, source.Lid)
	result.ClientID = chooseNonEmpty(base.ClientID, source.ClientID)

	// Map поля - выбираем тот, где больше non-zero values
	result.Attributes = chooseMapWithMoreData(base.Attributes, source.Attributes)
	result.Relationships = chooseMapWithMoreData(base.Relationships, source.Relationships)

	// Для указателей на Links и Meta
	result.Links = choosePointerMapWithMoreData(base.Links, source.Links)
	result.Meta = choosePointerMapWithMoreData(base.Meta, source.Meta)

	return result
}

// choosePointerMapWithMoreData выбирает указатель на map с бОльшим количеством non-zero значений
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

	// Считаем non-zero values в каждом
	sourceCount := countNonZeroValuesInMap(*source)
	baseCount := countNonZeroValuesInMap(*base)

	// Выбираем тот, где больше non-zero values
	if sourceCount > baseCount {
		return source
	} else if baseCount > sourceCount {
		return base
	}

	// Если количество одинаковое - приоритет у source
	return source
}

// chooseMapWithMoreData выбирает map с бОльшим количеством non-zero значений
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

	// Считаем non-zero values в каждом
	sourceCount := countNonZeroValuesInMap(source)
	baseCount := countNonZeroValuesInMap(base)

	// Выбираем тот, где больше non-zero values
	if sourceCount > baseCount {
		return source
	} else if baseCount > sourceCount {
		return base
	}

	// Если количество одинаковое - приоритет у source
	return source
}

// countNonZeroValuesInMap подсчитывает количество non-zero значений в map
func countNonZeroValuesInMap(m map[string]interface{}) int {
	if m == nil {
		return 0
	}

	count := 0
	for _, value := range m {
		if !isZeroValue(value) {
			count++
		}
	}
	return count
}

// chooseNonEmpty выбирает непустую строку (приоритет source)
func chooseNonEmpty(base, source string) string {
	if source != "" {
		return source
	}
	return base
}

// isZeroValue проверяет, является ли значение zero-value
func isZeroValue(v interface{}) bool {
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
	default:
		// Для сложных типов считаем non-zero
		return false
	}
}
