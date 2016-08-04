package main

import "testing"

func TestConnectToDB(t *testing.T) {
	path := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\raymond3chou\\VR\\VR GoSQL\\definition.mdb"
	connectToDB(path)
}

func TestQueryGenerator(t *testing.T) {
	fields := []string{"A", "B d", "C"}
	q := queryGenerator(fields)
	if q != "SELECT [A],[B d],[C] FROM [AV Sparing 2013 FU] WHERE PTID=?" {
		t.Errorf("The Printed Query was %s", q)
	}
}

func TestQueryTable(t *testing.T) {
	path := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\TestDB.accdb"
	conn := connectToDB(path)
	defer conn.Close()
	query := "SELECT [PTID],[CHART] FROM [AV Sparing 2013 FU] WHERE PTID=?"
	ptid := "Test"
	orderedMapSlice := queryTable(conn, query, ptid)
	if orderedMapSlice[0][0].Value != "Test" {
		t.Errorf("%s", orderedMapSlice[0][0].Value)
	}
}

func TestPtidList(t *testing.T) {
	path := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\TestDB.accdb"
	conn := connectToDB(path)
	defer conn.Close()
	list := ptidList(conn)
	if list[0] != "Test" {
		t.Errorf("the first PTID is not Test but is %s", list[0])
	}
}

func TestIteratePTID(t *testing.T) {

	path := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\TestDB.accdb"
	fields := []string{"PTID", "CHART"}
	iteratePTID(path, fields)

}
