package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	_ "github.com/alexbrainman/odbc"
	"github.com/raymond3chou/VR/accessHelper"
)

func TestCreateErrorLog(t *testing.T) {
	expectedOutput := "Success: C:\\Users\\raymond chou\\Desktop\\TestLog.log"

	//Checks if Testlog is created using path and works
	path := "C:\\Users\\raymond chou\\Desktop\\TestLog.log"
	accessHelper.CreateFile(path)
	errFile, err := accessHelper.ConnectToTxt(path)
	if !err {
		log.Fatal("Did not connect\n")
	}
	log.SetOutput(errFile)
	log.Println("Success: C:\\Users\\raymond chou\\Desktop\\TestLog.log")
	output := accessHelper.ReadFile("C:\\Users\\raymond chou\\Desktop\\TestLog.log")
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Output: %s does not Contain Expected Output: %s", output, expectedOutput)
	}
	errFile.Close()
	error := os.Remove("C:\\Users\\raymond chou\\Desktop\\TestLog.log")
	if error != nil {
		t.Error(error)
	}
}

func TestCheckFollowup(t *testing.T) {
	fileErr := accessHelper.CreateErrorLog("C:\\Users\\raymond chou\\Desktop\\TestError.log")
	log.SetOutput(fileErr)
	defer fileErr.Close()
	actualqueryforT1 := " [PTID], [CHART], [SEX], [STREET], [POSTCODE]"
	tableArray := []string{"Contact Info", "Table1"}
	conn := connectToDB("C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\goSql\\Example", "TestDb.mdb")
	defer conn.Close()
	conn.Ping()
	for _, tablename := range tableArray {
		FU, _, query := checkFollowup(conn, tablename)
		if tablename == "Contact Info" && FU == true {
			t.Errorf("%s is not a Follow up table but function returns that it is", tablename)
		}
		if tablename == "Table1" && query != actualqueryforT1 {
			t.Errorf("Actual Query for %s was %s not %s", tablename, query, actualqueryforT1)
		}
		// if tablename == "Copy of Table1" && query != actualqueryforCopy {
		// 	t.Errorf("Actual Query for %s was %s not %s", tablename, query, actualqueryforCopy)
		// }
	}
}

func TestConvertToString(t *testing.T) {
	vals := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		vals[i] = new(sql.NullString)
	}
	stringval := make([]string, 10)
	row := accessHelper.ConvertToString(vals)
	actualtype := reflect.TypeOf(row)
	expectedtype := reflect.TypeOf(stringval)

	if actualtype != expectedtype {
		t.Errorf("TestConvertToString failed because actual type is %s whereas expected type is %s", actualtype, expectedtype)
	}
}

func TestConvertToText(t *testing.T) {
	dbq := "C:\\Desktop\\Test.mdb"
	maincol := []string{"A", "B", "C", "D", "E", "G"}
	cols := make([]accessHelper.OrderedMap, 4)
	cols[0] = accessHelper.OrderedMap{Colname: "A", Value: "a"}
	cols[1] = accessHelper.OrderedMap{Colname: "B", Value: "b"}
	cols[2] = accessHelper.OrderedMap{Colname: "C", Value: " "}
	cols[3] = accessHelper.OrderedMap{Colname: "D_wdad", Value: "d"}
	actualrow := accessHelper.ConvertToText(maincol, cols, dbq)
	expectedrow := "\na|b| |d||"

	if actualrow != expectedrow {
		t.Errorf("TestConvertToText failed because the converted string is %s instead of %s", actualrow, expectedrow)
	}
}

func TestConvertToOrderedMap(t *testing.T) {
	//can not test if the values matched to the correct map key because maps are random.
	actualcols := make([]accessHelper.OrderedMap, 4)
	actualcols[0] = accessHelper.OrderedMap{Colname: "A", Value: ""}
	actualcols[1] = accessHelper.OrderedMap{Colname: "B", Value: ""}
	actualcols[2] = accessHelper.OrderedMap{Colname: "C", Value: ""}
	actualcols[3] = accessHelper.OrderedMap{Colname: "D", Value: "e"}

	rowstring := []string{"a", "b", "c", "d"}
	rowwithcols := accessHelper.ConvertToOrderedMap(actualcols, rowstring)
	for i := range rowwithcols {
		if actualcols[i].Value == "" {
			t.Errorf("actualcols[%s]='%s' instead of ''", actualcols[i].Colname, actualcols[i].Value)
		}
	}
}

func TestFileWrite(t *testing.T) {
	fileErr := accessHelper.CreateErrorLog("C:\\Users\\raymond chou\\Desktop\\TestError.log")
	log.SetOutput(fileErr)
	defer fileErr.Close()

	path := "C:\\Users\\raymond chou\\Desktop\\TestFile.txt"
	accessHelper.CreateFile(path)
	file, _ := accessHelper.ConnectToTxt(path)

	row := "Hello World 1234..."
	accessHelper.FileWrite(file, row)
	file.Close()
	actualrow := accessHelper.ReadFile(path)
	if row != actualrow {
		t.Errorf("Read %s but Wrote %s", actualrow, row)
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}

func TestSelectAccess(t *testing.T) {
	fileErr := accessHelper.CreateErrorLog("C:\\Users\\raymond chou\\Desktop\\TestError.log")
	log.SetOutput(fileErr)
	defer fileErr.Close()
	//create file
	dbq := "C:\\Desktop\\Test.mdb"
	path := "C:\\Users\\raymond chou\\Desktop\\TestWrite.txt"
	accessHelper.CreateFile(path)
	file, _ := accessHelper.ConnectToTxt(path)

	//connect to database
	conn := connectToDB("./Example", "TestDb.accdb")
	defer conn.Close()
	//array of tablenames
	tableArray := []string{"AV Sparing 2013 FU", "ContactInfo3", "Copy Of Table1", "Table1", "Table2"}

	for _, tablename := range tableArray {
		insertedRows, numberofRows := selectAccess(conn, file, tablename, dbq)
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
	actualOutput := accessHelper.ReadFile(path)
	if expectedOutput != actualOutput {
		t.Errorf("Read %s but should be %s", actualOutput, expectedOutput)
	}
	file.Close()
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}

func TestFindTable(t *testing.T) {
	conn := connectToDB("C:/Users/raymond chou/Desktop/WorkingFiles/VR/Example", "New Microsoft Access Database.accdb")
	tablenames := findTable(conn)
	expectedtablenames := []string{"Table1", "Table5"}
	for i, tablename := range tablenames {
		if tablename != expectedtablenames[i] {
			t.Errorf("Read %s but expected %s", tablename, expectedtablenames[i])
		}
	}
}

func TestWalkDir(t *testing.T) {
	foldername := []string{""}
	walkDir(foldername, "./")
}
