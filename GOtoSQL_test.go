package main

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestMatchTable(t *testing.T) {
	tablenames := []string{"contact info", "contactInfo", "contactin19fo", "infoinfo", "inf"}
	match := "info"
	actualphitablenames := matchTable(tablenames, match)
	expectedPHItablenames := []string{"contact info", "infoinfo"}
	nummatched := 0
	for _, expected := range expectedPHItablenames {
		for _, actual := range actualphitablenames {
			if expected == actual {
				nummatched++
			}
		}
	}
	if nummatched != len(expectedPHItablenames) {
		t.Errorf("TestMatchTable failed because %d matched instead of %d", nummatched, len(expectedPHItablenames))
	} else {
		t.Log("MatchTable Works")
	}
}
func TestConvertToString(t *testing.T) {
	vals := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		vals[i] = new(sql.NullString)
	}
	stringval := make([]string, 10)
	row := convertToString(vals)
	actualtype := reflect.TypeOf(row)
	expectedtype := reflect.TypeOf(stringval)

	if actualtype != expectedtype {
		t.Errorf("TestConvertToString failed because actual type is %s whereas expected type is %s", actualtype, expectedtype)
	} else {
		t.Log("convertToString Works")
	}
}
func TestConvertToText(t *testing.T) {
	rowstring := make([]string, 10)
	for i := range rowstring {
		rowstring[i] = "Hello1&3"
	}
	actualrow := convertToText(rowstring)
	expectedrow := "\nHello1&3|Hello1&3|Hello1&3|Hello1&3|Hello1&3|Hello1&3|Hello1&3|Hello1&3|Hello1&3|Hello1&3"
	if actualrow == expectedrow {
		t.Error("convertToText Works")
	} else {
		t.Errorf("TestConvertToText failed because the converted string is %s instead of %s", actualrow, expectedrow)
	}
}
func TestFileWrite(t *testing.T) {

}
