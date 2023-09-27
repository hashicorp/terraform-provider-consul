// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func formatKey(key string) string {
	tokens := strings.Split(key, "_")
	keyToReturn := ""
	for _, token := range tokens {
		if token == "tls" || token == "ttl" {
			keyToReturn += strings.ToUpper(token)
		} else {
			caser := cases.Title(language.English)
			keyToReturn += caser.String(token)
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
