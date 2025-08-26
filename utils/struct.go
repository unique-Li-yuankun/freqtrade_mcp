package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func StructJsonParams(s any) []string {
	var params []string
	ts := reflect.TypeOf(s).Elem()
	vs := reflect.ValueOf(s).Elem()
	for i := 0; i < ts.NumField(); i++ {
		field := ts.Field(i)
		fieldV := vs.Field(i)
		key, ok := field.Tag.Lookup("json")
		if !ok {
			continue
		}

		var param string
		switch field.Type.Kind() {
		case reflect.Slice:
			var value []string
			for j := 0; j < fieldV.Len(); j++ {
				value = append(value, fmt.Sprintf("%v", fieldV.Index(j).Interface()))
			}
			param = fmt.Sprintf("--%s %v", key, strings.Join(value, " "))
		default:
			param = fmt.Sprintf("--%s %v", key, fieldV.Interface())
		}
		params = append(params, param)
	}
	return params
}
