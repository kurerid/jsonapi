package jsonapi

import (
	"time"
)

// MergeNodes объединяет две структуры Node, предпочитая non-zero значения из source
func mergeNodes(base, source *Node) *Node {
	result := &Node{}

	// Простые строковые поля
	result.Type = chooseNonEmpty(base.Type, source.Type)
	result.ID = chooseNonEmpty(base.ID, source.ID)
	result.Lid = chooseNonEmpty(base.Lid, source.Lid)
	result.ClientID = chooseNonEmpty(base.ClientID, source.ClientID)

	// Map поля - выбираем ту, где есть non-zero values
	result.Attributes = chooseNonEmptyMap(base.Attributes, source.Attributes)
	result.Relationships = chooseNonEmptyMap(base.Relationships, source.Relationships)

	// Для указателей на Links и Meta - особая логика
	result.Links = chooseNonEmptyLinks(base.Links, source.Links)
	result.Meta = chooseNonEmptyMeta(base.Meta, source.Meta)

	return result
}

// chooseNonEmptyLinks выбирает Links с non-zero values
func chooseNonEmptyLinks(base, source *Links) *Links {
	if source == nil && base == nil {
		return nil
	}
	if source == nil {
		return base
	}
	if base == nil {
		return source
	}

	// Оба не nil - выбираем тот, где есть данные
	if hasNonZeroValuesInLinks(*source) {
		return source
	}
	if hasNonZeroValuesInLinks(*base) {
		return base
	}

	// Оба пустые - возвращаем source (или base)
	return source
}

// chooseNonEmptyMeta выбирает Meta с non-zero values
func chooseNonEmptyMeta(base, source *Meta) *Meta {
	if source == nil && base == nil {
		return nil
	}
	if source == nil {
		return base
	}
	if base == nil {
		return source
	}

	// Оба не nil - выбираем тот, где есть данные
	if hasNonZeroValuesInMeta(*source) {
		return source
	}
	if hasNonZeroValuesInMeta(*base) {
		return base
	}

	// Оба пустые - возвращаем source (или base)
	return source
}

// hasNonZeroValuesInLinks проверяет, содержит ли Links хотя бы одно non-zero значение
func hasNonZeroValuesInLinks(links Links) bool {
	if links == nil {
		return false
	}
	for _, value := range links {
		if !isZeroValue(value) {
			return true
		}
	}
	return false
}

// hasNonZeroValuesInMeta проверяет, содержит ли Meta хотя бы одно non-zero значение
func hasNonZeroValuesInMeta(meta Meta) bool {
	if meta == nil {
		return false
	}
	for _, value := range meta {
		if !isZeroValue(value) {
			return true
		}
	}
	return false
}

// chooseNonEmptyMap выбирает map с non-zero values
func chooseNonEmptyMap(base, source map[string]interface{}) map[string]interface{} {
	if hasNonZeroValuesInMap(source) {
		return source
	}
	if hasNonZeroValuesInMap(base) {
		return base
	}
	return source
}

// hasNonZeroValuesInMap проверяет, содержит ли map хотя бы одно non-zero значение
func hasNonZeroValuesInMap(m map[string]interface{}) bool {
	if m == nil {
		return false
	}
	for _, value := range m {
		if !isZeroValue(value) {
			return true
		}
	}
	return false
}

// chooseNonEmpty выбирает непустую строку
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
		return false
	}
}
