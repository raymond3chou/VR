package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/alexbrainman/odbc"
	"github.com/raymond3chou/VR/accessHelper"
	"github.com/raymond3chou/VR/excelHelper"
)

//Object struct to sort duplicates
type Object struct {
	Field string
	Date  string
}

//Basic struct contains the basic information
type Basic struct {
	PTID       string
	MRN        string
	ResearchID string
}

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

//Phone struct for each unique work or home phone
type Phone struct {
	Type       string `json:"phone_type"`
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	PhoneNum   string `json:"phone_number"`
	Date       string `json:"FU_Date"`
}

//Address unique address take date,PTID,mrn,research_id
type Address struct {
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostCode   string `json:"postal_code"`
	Date       string `json:"FU_Date"`
}

//Cardio for each unique cardiologist  take date,PTID,mrn,research_id
type Cardio struct {
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	Cardio     string `json:"cardiologist"`
	Date       string `json:"FU_Date"`
}

//GP struct for each unique GP
type GP struct {
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	GP         string `json:"general_practitioner"`
	Date       string `json:"FU_Date"`
}

//for each unique GP take date,PTID,mrn,research_id
//patients:
//should be only one set, match with VR

func iteratePTID(path string, fields []string) {
	conn := connectToDB(path)
	pList := ptidList(conn)
	query := queryGenerator(fields)

	jsonPath := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\JSON\\email.json"
	emailJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer emailJSONFile.Close()
	jsonPath = "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\JSON\\phone.json"
	phoneJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer phoneJSONFile.Close()
	jsonPath = "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\JSON\\address.json"
	addressJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer addressJSONFile.Close()
	jsonPath = "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\JSON\\GP.json"
	GPJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer GPJSONFile.Close()
	jsonPath = "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PHIJsonConvertor\\JSON\\cardio.json"
	cardioJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)

	for _, p := range pList {
		ptidInfoSlice := queryTable(conn, query, p)
		createJSON(ptidInfoSlice, emailJSONFile, phoneJSONFile, addressJSONFile, GPJSONFile, cardioJSONFile)
	}
	conn.Close()
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
	j := -1
	var orderedMapSlice []([]accessHelper.OrderedMap)
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			log.Println(err)
		}

		rowString := accessHelper.ConvertToString(vals)

		orderedMap := accessHelper.ConvertToOrderedMap(valueMap, rowString)
		//appends a struct slice on to a slice of slice struct
		j++
		for i := 0; i < len(orderedMap); i++ {
			orderedMapSlice = append(orderedMapSlice, make([]accessHelper.OrderedMap, len(orderedMap)))
			orderedMapSlice[j][i].Colname = orderedMap[i].Colname
			orderedMapSlice[j][i].Value = orderedMap[i].Value
		}
		//end of append function
	}
	return orderedMapSlice
}

//compareObjects returns the duplicate that should be deleted
func compareObjects(objects []Object) int {
	for i := 0; i < len(objects); i++ {
		for j := i + 1; j < len(objects); j++ {
			if j != i && objects[i].Field == objects[j].Field {
				if compareDates(objects[i].Date, objects[j].Date) {
					return j
				}
				return i
			}
		}
	}
	return -1
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

func createJSON(oMS [][]accessHelper.OrderedMap, emailJSONFile *os.File, phoneJSONFile *os.File, addressJSONFile *os.File, gpJSONFile *os.File, cardioJSONFile *os.File) {
	for i := 0; i < len(oMS); i++ {
		for j := 0; j < len(oMS[i]); j++ {
			if oMS[i][j].Colname == "FU_D" {
				date := strings.Split(oMS[i][j].Value, "T")
				oMS[i][j].Value = date[0]
			}
		}
	}
	emailObject(oMS, emailJSONFile)
	phoneObject(oMS, phoneJSONFile)
	gpObject(oMS, gpJSONFile)
	cardioObject(oMS, cardioJSONFile)
	addressObject(oMS, addressJSONFile)
}

func emailObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File) {
	basic := getBasic(oMS)
	emails := cleanObject(oMS, "EMAIL")
	createEmail(emails, basic, jsonFile)
}

func phoneObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File) {
	basic := getBasic(oMS)
	phonesHome := cleanObject(oMS, "PHONEHOME")
	fmt.Println(phonesHome)
	createPhones(phonesHome, basic, jsonFile, true)
	phonesWork := cleanObject(oMS, "PHONEWORK")
	createPhones(phonesWork, basic, jsonFile, false)
}

func gpObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File) {
	basic := getBasic(oMS)
	gp1 := cleanObject(oMS, "GP1")
	createGP(gp1, basic, jsonFile)
	gp2 := cleanObject(oMS, "GP2")
	createGP(gp2, basic, jsonFile)
}

func cardioObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File) {
	basic := getBasic(oMS)
	cardio1 := cleanObject(oMS, "CARDIO1")
	createCardio(cardio1, basic, jsonFile)
	cardio2 := cleanObject(oMS, "CARDIO2")
	createCardio(cardio2, basic, jsonFile)
}

func addressObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File) {
	basic := getBasic(oMS)
	address := cleanAddress(oMS)
	createAddress(address, basic, jsonFile)
}

//getBasic gets the basic information and stores it in a Basic struct
func getBasic(oMS [][]accessHelper.OrderedMap) Basic {
	var b Basic
	for i := 0; i < len(oMS); i++ {
		for j := 0; j < len(oMS[i]); j++ {
			if oMS[i][j].Colname == "PTID" {
				if oMS[i][j].Value != "" {
					b.PTID = oMS[i][j].Value
				}
			}
			if oMS[i][j].Colname == "CHART" {
				if oMS[i][j].Value != "" {
					b.MRN = oMS[i][j].Value
				}

			}
		}
	}
	return b
}

func checkDate(c string) bool {
	if c == "FU_D" {
		return true
	}
	return false
}

func cleanObject(oMS [][]accessHelper.OrderedMap, field string) []Object {
	var fieldValue string
	var date string
	var objects []Object

	for i := 0; i < len(oMS); i++ {
		validField := false
		for j := 0; j < len(oMS[i]); j++ {
			if oMS[i][j].Colname == field {
				if oMS[i][j].Value != "" {
					fieldValue = oMS[i][j].Value
					validField = true
				}
			}
			if checkDate(oMS[i][j].Colname) {
				date = oMS[i][j].Value
			}
		}
		if validField {
			o := Object{fieldValue, date}
			fmt.Println(o)
			objects = append(objects, o)
		}
	}
	duplicates := 0
	for true {
		duplicates = compareObjects(objects)
		if duplicates != -1 {
			objects = append(objects[:duplicates], objects[duplicates+1:]...)
		} else {
			break
		}
	}
	return objects
}

func cleanAddress(oMS [][]accessHelper.OrderedMap) []Object {
	var fieldValue []string
	var date string
	var objects []Object

	for i := 0; i < len(oMS); i++ {
		validField := 0
		fieldValue = []string{}
		for j := 0; j < len(oMS[i]); j++ {
			if oMS[i][j].Colname == "STREET" {
				if oMS[i][j].Value != "" {
					fieldValue = append(fieldValue, oMS[i][j].Value)
					validField++
				}
				fieldValue = append(fieldValue, "")

			}
			if oMS[i][j].Colname == "CITY" {
				if oMS[i][j].Value != "" {
					fieldValue = append(fieldValue, oMS[i][j].Value)
					validField++
				}
				fieldValue = append(fieldValue, "")

			}
			if oMS[i][j].Colname == "PROVINCE" {
				if oMS[i][j].Value != "" {
					fieldValue = append(fieldValue, oMS[i][j].Value)
					validField++
				}
				fieldValue = append(fieldValue, "")

			}
			if oMS[i][j].Colname == "POSTCODE" {
				if oMS[i][j].Value != "" {
					fieldValue = append(fieldValue, oMS[i][j].Value)
					validField++
				}
				fieldValue = append(fieldValue, "")

			}
			if checkDate(oMS[i][j].Colname) {
				date = oMS[i][j].Value
			}
		}
		if validField > 0 {
			var fieldString string
			for _, f := range fieldValue {
				fieldString += f + "_"
			}
			fieldString = strings.TrimSuffix(fieldString, "_")
			o := Object{fieldString, date}
			fmt.Println(o)
			objects = append(objects, o)
		}
	}
	duplicates := 0
	for true {
		duplicates = compareObjects(objects)
		if duplicates != -1 {
			objects = append(objects[:duplicates], objects[duplicates+1:]...)
		} else {
			break
		}
	}
	return objects
}

func createPhones(phones []Object, basic Basic, jsonFile *os.File, home bool) {
	var e Phone
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	if home {
		e.Type = "home"
	} else {
		e.Type = "work"
	}
	for i := range phones {
		e.PhoneNum = phones[i].Field
		e.Date = phones[i].Date
		writeJSON(e, jsonFile)
	}
}

func createGP(gps []Object, basic Basic, jsonFile *os.File) {
	var e GP
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	for i := range gps {
		e.GP = gps[i].Field
		e.Date = gps[i].Date
		writeJSON(e, jsonFile)
	}
}

func createCardio(cardio []Object, basic Basic, jsonFile *os.File) {
	var e Cardio
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	for i := range cardio {
		e.Cardio = cardio[i].Field
		e.Date = cardio[i].Date
		writeJSON(e, jsonFile)
	}
}

func createEmail(emails []Object, basic Basic, jsonFile *os.File) {
	var e Email
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	for i := range emails {
		e.Email = emails[i].Field
		e.Date = emails[i].Date
		writeJSON(e, jsonFile)
	}
}

func createAddress(addr []Object, basic Basic, jsonFile *os.File) {
	var e Address
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	for i := range addr {
		a := strings.Split(addr[i].Field, "_")
		fmt.Println(a[1])

		e.Street = a[0]
		e.City = a[2]
		e.Province = a[4]
		e.PostCode = a[6]
		e.Date = addr[i].Date
		writeJSON(e, jsonFile)
	}
}

//writeJSON writes the struct into JSON format
func writeJSON(newEvent interface{}, jsonFile *os.File) {
	j, err := json.Marshal(newEvent)
	if err != nil {
		log.Println(err)
	}
	jsonFile.Write(j)
	jsonFile.WriteString("\n")
}
