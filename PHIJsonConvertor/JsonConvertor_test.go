package main

import (
	"fmt"
	"testing"
)

func TestConnectToDB(t *testing.T) {
	path := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\raymond3chou\\VR\\VR GoSQL\\definition.mdb"
	connectToDB(path)
}

func TestQueryGenerator(t *testing.T) {
	fields := []string{"A", "B d", "C"}
	q := queryGenerator(fields)
	if q != "SELECT [A],[B d],[C] FROM ContactInfo" {
		t.Errorf("The Printed Query was %s", q)
	}
}

func TestQueryTable(t *testing.T) {
	path := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\TestDB.accdb"
	conn := connectToDB(path)
	defer conn.Close()
	query := "SELECT PTID,CHART,Email FROM ContactInfo3"
	queryTable(conn, query)
}

func TestPtidList(t *testing.T) {
	path := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\TestDB.accdb"
	conn := connectToDB(path)
	defer conn.Close()
	list := ptidList(conn)
	fmt.Println(list)
}
