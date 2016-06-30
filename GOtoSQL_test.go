package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

//create mock files to test database and files.
//print errors
//check if the function works and if it doesnt
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
	}
}

func TestConvertToText(t *testing.T) {

	maincol := []string{"A", "B", "C", "D", "E"}
	cols := map[string]string{"A": "a", "B": "b", "C": "c", "D_wdad": "d"}
	actualrow := convertToText(maincol, cols)
	expectedrow := "\na|b|c|d"

	if actualrow == expectedrow {

	} else {
		t.Errorf("TestConvertToText failed because the converted string is %s instead of %s", actualrow, expectedrow)
	}
}
func TestConvertToMap(t *testing.T) {
	actualcols := map[string]string{"A": "", "B": "", "C": "", "D": "e"}
	rowstring := []string{"a", "b", "", "d"}
	expectedcols := map[string]string{"A": "a", "B": "b", "C": "", "D": "d"}

	rowwithcols := convertToMap(actualcols, rowstring)
	for col := range rowwithcols {
		if actualcols[col] != expectedcols[col] {
			t.Errorf("expectedcols[%s]='%s' instead of actualcols[%s]='%s'", col, actualcols[col], col, expectedcols[col])
		}
	}
}
func readFile(filename string) string {
	fileoutput, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	return string(fileoutput)
}

func TestFileWrite(t *testing.T) {
	path := "C:\\Users\\raymond chou\\Desktop\\TestFile.txt"

	var f, err = os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Close()

	file, _ := connectToTxt(path)
	row := "Hello World 1234..."
	fileWrite(file, row)
	file.Close()
	actualrow := readFile(path)
	if row != actualrow {
		t.Errorf("Read %s but Wrote %s", actualrow, row)
	}
	err = os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCheckFollowup(t *testing.T) {
	actualqueryforT1 := " PTID CHART SEX STREET POSTCODE"
	actualqueryforCopy := " CHART SEX STREET PTID POSTCODE"
	tableArray := []string{"ContactInfo3", "Table1", "Copy Of Table1"}
	conn := connectToDB("./", "TestDb.accdb")
	defer conn.Close()
	for _, tablename := range tableArray {
		FU, _, query := checkFollowup(conn, tablename)
		if tablename == "ContactInfo3" && FU == true {
			t.Errorf("%s is not a Follow up table but function returns that it is", tablename)
		}
		if tablename == "Table1" && query != actualqueryforT1 {
			t.Errorf("Actual Query for %s was %s not %s", tablename, query, actualqueryforT1)
		}
		if tablename == "Copy of Table1" && query != actualqueryforCopy {
			t.Errorf("Actual Query for %s was %s not %s", tablename, query, actualqueryforCopy)
		}
	}
}
