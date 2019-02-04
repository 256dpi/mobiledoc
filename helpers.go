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
	i, ok := v.(int)
	if ok {
		return i, true
	}

	// check int64 (bson)
	ii, ok := v.(int64)
	if ok {
		return int(ii), true
	}

	// check float64 (json)
	f, ok := v.(float64)
	if ok {
		return int(f), true
	}

	// check MarkerType (tests)
	mt, ok := v.(MarkerType)
	if ok {
		return int(mt), true
	}

	// check SectionType (tests)
	st, ok := v.(SectionType)
	if ok {
		return int(st), true
	}

	return 0, false
}
