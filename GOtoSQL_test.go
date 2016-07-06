package main

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/access"
	_ "github.com/alexbrainman/odbc"
)

func TestCheckFollowup(t *testing.T) {
	actualqueryforT1 := " PTID, CHART, SEX, STREET, POSTCODE"
	actualqueryforCopy := " CHART, SEX, STREET, PTID, POSTCODE"
	tableArray := []string{"ContactInfo3", "Table1", "Copy Of Table1"}
	conn := connectToDB("./Example", "TestDb.accdb")
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
	row := access.ConvertToString(vals)
	actualtype := reflect.TypeOf(row)
	expectedtype := reflect.TypeOf(stringval)

	if actualtype != expectedtype {
		t.Errorf("TestConvertToString failed because actual type is %s whereas expected type is %s", actualtype, expectedtype)
	}
}

func TestConvertToText(t *testing.T) {

	maincol := []string{"A", "B", "C", "D", "E", "G"}
	cols := make([]access.OrderedMap, 4)
	cols[0] = access.OrderedMap{Colname: "A", Value: "a"}
	cols[1] = access.OrderedMap{Colname: "B", Value: "b"}
	cols[2] = access.OrderedMap{Colname: "C", Value: " "}
	cols[3] = access.OrderedMap{Colname: "D_wdad", Value: "d"}
	actualrow := access.ConvertToText(maincol, cols)
	expectedrow := "\na|b| |d||"

	if actualrow != expectedrow {
		t.Errorf("TestConvertToText failed because the converted string is %s instead of %s", actualrow, expectedrow)
	}
}

func TestConvertToOrderedMap(t *testing.T) {
	//can not test if the values matched to the correct map key because maps are random.
	actualcols := make([]access.OrderedMap, 4)
	actualcols[0] = access.OrderedMap{Colname: "A", Value: ""}
	actualcols[1] = access.OrderedMap{Colname: "B", Value: ""}
	actualcols[2] = access.OrderedMap{Colname: "C", Value: ""}
	actualcols[3] = access.OrderedMap{Colname: "D", Value: "e"}

	rowstring := []string{"a", "b", "c", "d"}
	rowwithcols := access.ConvertToOrderedMap(actualcols, rowstring)
	for i := range rowwithcols {
		if actualcols[i].Value == "" {
			t.Errorf("actualcols[%s]='%s' instead of ''", actualcols[i].Colname, actualcols[i].Value)
		}
	}
}

func TestFileWrite(t *testing.T) {
	access.CreateErrorLog(true)
	path := "C:\\Users\\raymond chou\\Desktop\\TestFile.txt"
	access.CreateFile(path)
	file, _ := access.ConnectToTxt(path)

	row := "Hello World 1234..."
	access.FileWrite(file, row)
	file.Close()
	actualrow := access.ReadFile(path)
	if row != actualrow {
		t.Errorf("Read %s but Wrote %s", actualrow, row)
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}

func TestSelectAccess(t *testing.T) {
	access.CreateErrorLog(true)

	//create file
	path := "C:\\Users\\raymond chou\\Desktop\\TestWrite.txt"
	access.CreateFile(path)
	file, _ := access.ConnectToTxt(path)

	//connect to database
	conn := connectToDB("./Example", "TestDb.accdb")
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
	expectedOutput := "\nTest|1234567|Jim|Lim|1||Main St|Toronto|BC|M7K A6D|123456798||6478032654|@gmail.com||C1|C2|G1|G2\nTest2|1234567|Jim2|Lim2|1||Main St2|Toronto2|BC|M7K A6D||123456798|6478032654|@gmail.com||C1||G1|\nTestRow|135435|||1||Test St|||m27d9dq|||||||||\nTestRow|111111|||1||Test St|||m27d9dq|||||||||"
	actualOutput := access.ReadFile(path)
	if expectedOutput != actualOutput {
		t.Errorf("Read %s but should be %s", actualOutput, expectedOutput)
	}
	file.Close()
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}
