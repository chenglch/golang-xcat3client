package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func KeyValueArrayToMapArray(values []string) ([]interface{}, error) {
	// transform the string array like
	// -i mac=42:87:0a:05:00:00,primary=True,name=eth0 -i mac=42:87:0a:05:00:00,name=eth1`
	var m_array []interface{}
	for _, value := range values {
		m, err := KeyValueToMap(value, ",")
		if err != nil {
			return nil, err
		}
		m_array = append(m_array, m)
	}
	return m_array, nil
}

func KeyValueToMap(value string, sep string) (map[string]interface{}, error) {
	// transform the string like
	// bmc_address=11.0.0.0,bmc_password=password,bmc_username=admin
	m := make(map[string]interface{})
	temps := strings.Split(value, sep)
	for _, temp := range temps {
		item := strings.Split(temp, "=")
		if len(temp) < 2 {
			return nil, fmt.Errorf("The format of %s is not correct.", value)
		}
		b, err := strconv.ParseBool(item[1])
		if err == nil {
			m[item[0]] = b
		} else {
			m[item[0]] = item[1]
		}
	}
	return m, nil
}

func KeyValueArrayToMap(values []string, sep string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for _, value := range values {
		item := strings.Split(value, sep)
		if len(item) < 2 {
			return nil, fmt.Errorf("The format of %s is not correct.", item)
		}
		b, err := strconv.ParseBool(item[1])
		if err == nil {
			m[item[0]] = b
		} else {
			m[item[0]] = item[1]
		}
	}
	return m, nil
}

func ToNodeArray(value string) (names []string, err error) {
	items := strings.Split(value, ",")
	for _, item := range items {
		if strings.Contains(item, "[") && strings.Contains(item, "-") && strings.Contains(item, "]") {

			parts := strings.Split(item, "[")
			prefix := parts[0]

			num_parts := strings.Split(parts[1], "-")
			left, err := strconv.Atoi(num_parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid node format %s.", item)
			}
			right, err := strconv.Atoi(strings.Split(num_parts[1], "]")[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid node format %s.", item)
			}
			for left <= right {
				name := fmt.Sprintf("%s%d", prefix, left)
				left += 1
				names = append(names, name)
			}
		} else {
			item = strings.Replace(item, "[", "", -1)
			item = strings.Replace(item, "]", "", -1)
			names = append(names, item)
		}
	}

	names = RmDuplicate(names)
	return names, nil
}

func RmDuplicate(strs []string) (ret []string) {
	sort.Strings(strs)
	for i := 0; i < len(strs); i++ {
		if i > 0 && strs[i-1] == strs[i] {
			continue
		}
		ret = append(ret, strs[i])
	}
	return ret
}

func MergeMap(a map[string]interface{}, b map[string]interface{}) {
	for k, v := range b {
		a[k] = v
	}
}

func MergeSlice(a []interface{}, b []interface{}) {
	for _, value := range b {
		a = append(a, value)
	}
}

func InterfaceToMap(in interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	value := reflect.ValueOf(in)
	if value.Kind() == reflect.Map {
		for _, key := range value.MapKeys() {
			v := value.MapIndex(key)
			key, ok := key.Interface().(string)
			if ok != true {
				return nil
			}
			m[key] = v.Interface()
		}
	}
	return m
}

func InterfaceToSlice(in interface{}) []interface{} {
	s := make([]interface{}, 0)
	if reflect.TypeOf(in).Kind() == reflect.Slice {
		value := reflect.ValueOf(in)
		for i := 0; i < value.Len(); i++ {
			s = append(s, value.Index(i).Interface())
		}
	}
	return s
}

func PrintJson(in interface{}) {
	var out bytes.Buffer
	json.Indent(&out, in.([]byte), "", "\t")
	out.WriteTo(os.Stdout)
	fmt.Printf("\n")
}

func Contains(collection interface{}, obj interface{}) (bool, error) {
	collectionValue := reflect.ValueOf(collection)
	switch reflect.TypeOf(collection).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < collectionValue.Len(); i++ {
			if collectionValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if collectionValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in")
}

func GetBytes(key interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	return buf, err
}
