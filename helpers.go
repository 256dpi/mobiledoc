package mobiledoc

func contains(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}

	return false
}

func toInt(v interface{}) (int, bool) {
	i, ok := v.(int)
	if ok {
		return i, true
	}

	f, ok := v.(float64)
	if ok {
		return int(f), true
	}

	return 0, false
}
