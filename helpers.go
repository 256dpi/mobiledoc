package mobiledoc

func contains(list []string, str string) bool {
	// check existence
	for _, item := range list {
		if item == str {
			return true
		}
	}

	return false
}

func toInt(v interface{}) (int, bool) {
	// check int
	if i, ok := v.(int); ok {
		return i, true
	}

	// check int64 (bson)
	if i, ok := v.(int64); ok {
		return int(i), true
	}

	// check float64 (json)
	if f, ok := v.(float64); ok {
		return int(f), true
	}

	// check MarkerType (tests)
	if mt, ok := v.(MarkerType); ok {
		return int(mt), true
	}

	// check SectionType (tests)
	if st, ok := v.(SectionType); ok {
		return int(st), true
	}

	return 0, false
}
