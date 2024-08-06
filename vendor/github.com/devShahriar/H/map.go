package H

import (
	"reflect"
	"strconv"
)

func PopulateStructFromMap(obj interface{}, data map[string]string) {
	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		key := field.Tag.Get("key")

		if value, ok := data[key]; ok {
			fieldValue := objValue.Field(i)
			setValue(fieldValue, value)
		}
	}
}

func setValue(fieldValue reflect.Value, value string) {

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int:
		if intValue, err := strconv.Atoi(value); err == nil {
			fieldValue.SetInt(int64(intValue))
		}
	}
}
