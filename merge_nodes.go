package jsonapi

func mergeNodes(base, source *Node) *Node {
	result := Node{}

	// Простые поля
	result.Type = firstNonEmpty(source.Type, base.Type)
	result.ID = firstNonEmpty(source.ID, base.ID)
	result.Lid = firstNonEmpty(source.Lid, base.Lid)
	result.ClientID = firstNonEmpty(source.ClientID, base.ClientID)

	// Map - если source не nil и не пустой, берем source
	if source.Attributes != nil &&
		len(source.Attributes) > 0 &&
		len(source.Attributes) > len(base.Attributes) {
		result.Attributes = source.Attributes
	} else {
		result.Attributes = base.Attributes
	}

	if source.Relationships != nil &&
		len(source.Relationships) > 0 &&
		len(source.Relationships) > len(base.Relationships) {
		result.Relationships = source.Relationships
	} else {
		result.Relationships = base.Relationships
	}

	// Указатели на структуры
	if source.Links != nil {
		result.Links = source.Links
	} else {
		result.Links = base.Links
	}

	if source.Meta != nil {
		result.Meta = source.Meta
	} else {
		result.Meta = base.Meta
	}

	return &result
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
