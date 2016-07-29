package event

import (
	"reflect"

	"github.com/raymond3chou/VR/excelHelper"
)

//CutArray takes in a map of string and a slice of required fields
//returns a map with only key's matching the required fields
func CutArray(rowSlice map[string]string, reqFields map[string]string) map[string]string {
	newRowSlice := make(map[string]string)
	for rf, fieldName := range reqFields {
		for key := range rowSlice {
			if rf == key {
				newRowSlice[fieldName] = rowSlice[key]
				break
			}
		}
	}
	return newRowSlice
}

//AssignBasic assigns the basic values to their respective fields(e.g PTID,MRN)
func AssignBasic(rowSlice map[string]string, event interface{}, basicStrFields []string, basicIntFields []string) interface{} {
	for _, bStr := range basicStrFields {
		valueStr := rowSlice[bStr]
		reflect.ValueOf(&event).Elem().FieldByName(bStr).SetString(valueStr)
	}
	for _, bInt := range basicIntFields {
		valueInt := excelHelper.StringToInt(rowSlice[bInt])
		reflect.ValueOf(&event).Elem().FieldByName(bInt).SetInt(valueInt)
	}
	return event
}
