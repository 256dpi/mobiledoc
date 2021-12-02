package mobiledoc

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

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
	// convert numbers
	switch i := v.(type) {
	case int:
		return i, true
	case int32: // bson
		return int(i), true
	case int64: // bson
		return int(i), true
	case float64: // json
		return int(i), true
	case MarkerType:
		return int(i), true
	case SectionType:
		return int(i), true
	default:
		r := reflect.ValueOf(v)
		switch r.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int(r.Int()), true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int(r.Uint()), true
		case reflect.Float32, reflect.Float64:
			return int(r.Float()), true
		default:
			return 0, false
		}
	}
}

func toMap(v interface{}) (Map, bool) {
	// convert map
	switch m := v.(type) {
	case Map:
		return m, true
	case bson.M:
		return m, true
	default:
		return nil, false
	}
}

func toList(v interface{}) (List, bool) {
	// convert list
	switch l := v.(type) {
	case List:
		return l, true
	case bson.A:
		return l, true
	default:
		return nil, false
	}
}
