package service

import (
	"reflect"
	"strconv"
)

func ConvertFloat64(s string) (float float64) {
	float, _ = strconv.ParseFloat(s, 64)
	return
}

func InSlice(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
