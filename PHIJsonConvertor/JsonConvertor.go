package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/alexbrainman/odbc"
	"github.com/raymond3chou/VR/accessHelper"
	"github.com/raymond3chou/VR/excelHelper"
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

func iteratePTID(path string, fields []string) {
	conn := connectToDB(path)
	pList := ptidList(conn)
	query := queryGenerator(fields)
	for _, p := range pList {
		ptidInfoSlice := queryTable(conn, query, p)
		createJSON(ptidInfoSlice)
	}
}

//connectToDB Connects to a specified database in a specified directory
func connectToDB(path string) *sql.DB {
	conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq="+path)
	if err != nil {
		log.Println("Connection to " + path + " Failed")
	}
	return conn
}

//ptidList returns a list of unique PTIDs
func ptidList(conn *sql.DB) []string {
	var ptidList []string
	rows, err := conn.Query("SELECT PTID FROM [AV Sparing 2013 FU] GROUP BY PTID")
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
	fieldString += " FROM [AV Sparing 2013 FU] WHERE PTID=?"
	return fieldString
}

func queryTable(conn *sql.DB, query string, ptid string) [][]accessHelper.OrderedMap {
	rows, err := conn.Query(query, ptid)
	if err != nil {
		log.Fatalf("Query: %s failed to run\n", query)
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

	var orderedMapSlice [][]accessHelper.OrderedMap
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			log.Println(err)
		}
		rowString := accessHelper.ConvertToString(vals)
		orderedMap := accessHelper.ConvertToOrderedMap(valueMap, rowString)
		orderedMapSlice = append(orderedMapSlice, orderedMap)
	}
	return orderedMapSlice
}

func compareObjects(s []string, d []string) map[string]string {
	var sd map[string]string
	if len(s) != len(d) {
		return sd
	}
	for i := 0; i < len(s)-1; i++ {
		dup := false
		for j := i + 1; j < len(s)-1; j++ {
			if s[i] == s[j] {
				dup = true
				if compareDates(d[i], d[j]) {
					sd[s[i]] = d[i]
				} else {
					sd[s[i]] = d[j]
				}
				break
			}
		}
		if !dup {
			sd[s[i]] = d[i]
		}
	}

	return sd
}

//compareDates returns true if the first date is earlier than the second
func compareDates(d1 string, d2 string) bool {
	date1 := strings.Split(d1, "-")
	date2 := strings.Split(d2, "-")

	for i := 0; i < 3; i++ {
		if excelHelper.StringToInt(date1[i]) > excelHelper.StringToInt(date2[i]) {
			return true
		}
	}
	return false
}

func createJSON(oMS [][]accessHelper.OrderedMap) {
	var emails []string
	var date []string

	for i := 0; i < len(oMS); i++ {
		validEmail := false
		for j := 0; j < len(oMS[i]); j++ {
			if oMS[i][j].Colname == "EMAIL" {
				if oMS[i][j].Value != "" {
					emails = append(emails, oMS[i][j].Value)
					validEmail = true
				}
			}
			if oMS[i][j].Colname == "FU_D" && validEmail {
				date = append(date, oMS[i][j].Value)
			}
		}
	}

}
