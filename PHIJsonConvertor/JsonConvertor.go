package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/alexbrainman/odbc"
	"github.com/raymond3chou/VR/accessHelper"
)

//query for PTID
//collections : email, phone, address, personal info, cardio, GP

//Email for each unique email take date, ptid,mrn,research_id
type Email struct {
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	Email      string `json:"email"`
	Date       string `json:"FU_Date"`
}

//phone:
//for each unique phone take date,PTID,mrn,research_id
//address:
//for each unique address take date,PTID,mrn,research_id
//cardio:
//for each unique cardio take date,PTID,mrn,research_id
//GP:
//for each unique GP take date,PTID,mrn,research_id
//patients:
//should be only one set, match with VR

//connectToDB Connects to a specified database in a specified directory
func connectToDB(path string) *sql.DB {
	log.Println("Connecting to " + path)
	conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq="+path)
	if err != nil {
		log.Println("Connection to " + path + " Failed")
	} else {
		log.Println("Connected to " + path)
	}
	return conn
}

func ptidList(conn *sql.DB) []string {
	var ptidList []string
	rows, err := conn.Query("SELECT PTID FROM ContactInfo3 GROUP BY PTID")
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var r string
		err = rows.Scan(&r)
		if err != nil {
			log.Println(err)
		}
		ptidList = append(ptidList, r)
	}
	return ptidList
}

//queryGenerator generates a SQL quert based on a list of fields
func queryGenerator(fields []string) string {
	var fieldString string
	fieldString += "SELECT "
	for _, f := range fields {
		fieldString += "[" + f + "],"
	}
	fieldString = strings.TrimSuffix(fieldString, ",")
	fieldString += " FROM ContactInfo"
	return fieldString
}

func queryTable(conn *sql.DB, query string) {
	rows, err := conn.Query(query)
	if err != nil {
		log.Printf("Query: %s failed to run\n", query)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println(err)
	}
	valueMap := make([]accessHelper.OrderedMap, len(cols))
	for i, colname := range cols {
		valueMap[i].Colname = colname
		valueMap[i].Value = ""
	}

	vals := make([]interface{}, len(cols))
	for i := range cols {
		vals[i] = new(sql.NullString)
	}

	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			log.Println(err)
		}
		rowString := accessHelper.ConvertToString(vals)
		orderedMap := accessHelper.ConvertToOrderedMap(valueMap, rowString)
		createJSON(orderedMap)
	}

}

func createJSON(orderedMap []accessHelper.OrderedMap) {

}
