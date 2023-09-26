// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"reflect"
	"strings"
)

func isSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array
}

func isMap(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}

func isSetSchema(v interface{}) bool {
	return reflect.TypeOf(v).String() == "*schema.Set"
}

func isStruct(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

func formatKey(key string) string {
	tokens := strings.Split(key, "_")
	keyToReturn := ""
	for _, token := range tokens {
		if token == "tls" {
			keyToReturn += strings.ToUpper(token)
		} else {
			keyToReturn += strings.ToTitle(token)
		}
	}
	return keyToReturn
}

func formatKeys(config interface{}, formatFunc func(string) string) (interface{}, error) {
	if isMap(config) {
		formattedMap := make(map[string]interface{})
		for key, value := range config.(map[string]interface{}) {
			formattedKey := formatFunc(key)
			formattedValue, err := formatKeys(value, formatKey)
			if err != nil {
				return nil, err
			}
			if formattedValue != nil {
				formattedMap[formattedKey] = formattedValue
			}
		}
		return formattedMap, nil
	} else if isSlice(config) {
		var newSlice []interface{}
		listValue := config.([]interface{})
		for _, elem := range listValue {
			newElem, err := formatKeys(elem, formatKey)
			if err != nil {
				return nil, err
			}
			newSlice = append(newSlice, newElem)
		}
		return newSlice, nil
	} else if isStruct(config) {
		var modifiedStruct map[string]interface{}
		jsonValue, err := json.Marshal(config)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonValue, &modifiedStruct)
		if err != nil {
			return nil, err
		}
		formattedStructKeys, err := formatKeys(modifiedStruct, formatKey)
		if err != nil {
			return nil, err
		}
		return formattedStructKeys, nil
	} else if isSetSchema(config) {
		valueList := config.(*schema.Set).List()
		if len(valueList) > 0 {
			formattedSetValue, err := formatKeys(valueList[0], formatKey)
			if err != nil {
				return nil, err
			}
			return formattedSetValue, nil
		}
		return nil, nil
	}
	return config, nil
}
