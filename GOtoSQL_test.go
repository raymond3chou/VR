package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

//create mock files to test database and files.
//print errors
//check if the function works and if it doesnt
func TestCheckFollowup(t *testing.T) {
	actualqueryforT1 := " PTID, CHART, SEX, STREET, POSTCODE"
	actualqueryforCopy := " CHART, SEX, STREET, PTID, POSTCODE"
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

	maincol := []string{"A", "B", "C", "D", "E", "G"}
	cols := make([]orderedMap, 4)
	cols[0] = orderedMap{"A", "a"}
	cols[1] = orderedMap{"B", "b"}
	cols[2] = orderedMap{"C", " "}
	cols[3] = orderedMap{"D_wdad", "d"}
	actualrow := convertToText(maincol, cols)
	expectedrow := "\na|b| |d||"

	if actualrow != expectedrow {
		t.Errorf("TestConvertToText failed because the converted string is %s instead of %s", actualrow, expectedrow)
	}
}

func TestConvertToOrderedMap(t *testing.T) {
	//can not test if the values matched to the correct map key because maps are random.
	actualcols := make([]orderedMap, 4)
	actualcols[0] = orderedMap{"A", ""}
	actualcols[1] = orderedMap{"B", ""}
	actualcols[2] = orderedMap{"C", ""}
	actualcols[3] = orderedMap{"D", "e"}

	rowstring := []string{"a", "b", "c", "d"}
	rowwithcols := convertToOrderedMap(actualcols, rowstring)
	for i := range rowwithcols {
		if actualcols[i].value == "" {
			t.Errorf("actualcols[%s]='%s' instead of ''", actualcols[i].colname, actualcols[i].value)
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

func createFile(path string) *os.File {
	var f, err = os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	file, _ := connectToTxt(path)
	return file
}

func TestFileWrite(t *testing.T) {
	path := "C:\\Users\\raymond chou\\Desktop\\TestFile.txt"
	file := createFile(path)

	row := "Hello World 1234..."
	fileWrite(file, row)
	file.Close()
	actualrow := readFile(path)
	if row != actualrow {
		t.Errorf("Read %s but Wrote %s", actualrow, row)
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}

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

func TestSelectAccess(t *testing.T) {
	//create file
	path := "C:\\Users\\raymond chou\\Desktop\\TestWrite.txt"
	file := createFile(path)
	//connect to database
	conn := connectToDB("./", "TestDb.accdb")
	defer conn.Close()
	//array of tablenames
	tableArray := []string{"AV Sparing 2013 FU", "ContactInfo3", "Copy Of Table1", "Table1", "Table2"}

	for _, tablename := range tableArray {
		insertedRows, numberofRows := selectAccess(conn, file, tablename)
		if tablename == "ContactInfo3" && insertedRows != 0 && numberofRows != 0 {
			t.Errorf("%s is not a follow up table but is read", tablename)
			t.Fail()
		}
		if tablename == "AV Sparing 2013 FU" && insertedRows != 2 && numberofRows != 2 {
			t.Errorf("%s has insertedRows:%d X | %d and numberofRows:%d X |%d", tablename, insertedRows, 2, numberofRows, 2)
			t.Fail()
		}
		if tablename == "Copy Of Table1" && insertedRows != 1 && numberofRows != 1 {
			t.Errorf("%s has insertedRows:%d X | %d and numberofRows:%d X |%d", tablename, insertedRows, 2, numberofRows, 2)
			t.Fail()
		}
		if tablename == "Table1" && insertedRows != 1 && numberofRows != 1 {
			t.Errorf("%s has insertedRows:%d X | %d and numberofRows:%d X |%d", tablename, insertedRows, 2, numberofRows, 2)
			t.Fail()
		}
		if tablename == "Table2" && insertedRows != 0 && numberofRows != 0 {
			t.Errorf("%s is not a follow up table but is read", tablename)
			t.Fail()
		}
	}
	expectedOutput := "\nTest|1234567|Jim|Lim|1||Main St|Toronto|BC|M7K A6D|123456798||6478032654|@gmail.com|\nTest2|1234567|Jim2|Lim2|1||Main St2|Toronto2|BC|M7K A6D||123456798|6478032654|@gmail.com|\nTestRow|135435|||1||Test St|||m27d9dq|||||\nTestRow|111111|||1||Test St|||m27d9dq|||||"
	actualOutput := readFile(path)
	if expectedOutput != actualOutput {
		t.Errorf("Read %s but should be %s", actualOutput, expectedOutput)
	}
	file.Close()
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
	//checks if it is a followup
	//queries the table
	//identifies the column names
	//scans the rows
	//writes the result in to a textfile * create a text file

}
