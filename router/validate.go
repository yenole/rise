package router

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Handler func(v string, field reflect.Value) bool

var validas = map[string]Handler{}

func init() {
	validas["eq"] = vEq
	validas["gte"] = vGte
	validas["lte"] = vLte
	validas["len"] = vLen
	validas["min"] = vMin
	validas["max"] = vMax
	validas["contain"] = vContain
	validas["endwith"] = vEndWith
	validas["startwith"] = vStartWith
}

func Validate(val any) error {
	v := reflect.ValueOf(val).Elem()
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get("v")
		if tag == "" {
			continue
		}
		field := v.Field(i)
		omit := strings.HasPrefix(tag, "omit")
		if omit {
			if field.Kind() == reflect.String && field.Len() == 0 {
				continue
			}
			if strings.HasPrefix(tag, "omit,") {
				tag = strings.TrimPrefix(tag, "omit,")
			} else {
				tag = strings.TrimPrefix(tag, "omit")
			}
		}

		tip := v.Type().Field(i).Tag.Get("tip")
		for _, kv := range strings.Split(tag, ",") {
			k, v := kv, ""
			if strings.Contains(kv, "=") {
				sp := strings.Split(kv, "=")
				k, v = sp[0], sp[1]
			}

			if fn, ok := validas[k]; ok {
				if !fn(v, field) {
					return errors.New(tip)
				}
			} else {
				return errors.New(tip)
			}
		}
	}
	return nil
}

func vGte(v string, field reflect.Value) bool {
	dst, _ := strconv.ParseFloat(v, 10)
	src, _ := strconv.ParseFloat(fmt.Sprint(field.Interface()), 10)
	return src >= dst
}

func vLte(v string, field reflect.Value) bool {
	dst, _ := strconv.ParseFloat(v, 10)
	src, _ := strconv.ParseFloat(fmt.Sprint(field.Interface()), 10)
	return src <= dst
}

func vLen(v string, field reflect.Value) bool {
	if field.Kind() == reflect.String {
		size, _ := strconv.Atoi(v)
		return field.Len() == size
	}
	return true
}

func vMin(v string, field reflect.Value) bool {
	if field.Kind() == reflect.String {
		size, _ := strconv.Atoi(v)
		return field.Len() >= size
	}
	return true
}

func vMax(v string, field reflect.Value) bool {
	if field.Kind() == reflect.String {
		size, _ := strconv.Atoi(v)
		return field.Len() <= size
	}
	return true
}

func vEq(v string, field reflect.Value) bool {
	return fmt.Sprint(field.Interface()) == v
}

func vContain(v string, field reflect.Value) bool {
	return field.Kind() == reflect.String && strings.Contains(field.String(), v)
}

func vStartWith(v string, field reflect.Value) bool {
	return field.Kind() == reflect.String && strings.HasPrefix(field.String(), v)
}

func vEndWith(v string, field reflect.Value) bool {
	return field.Kind() == reflect.String && strings.HasSuffix(field.String(), v)
}
