package main

import (
	"database/sql"
	"encoding/json"
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
	Path  []string
}

//PDate to keep track of patient dates
type PDate struct {
	Date string
	Row  int
}

//Basic struct contains the basic information
type Basic struct {
	PTID       string
	MRN        string
	ResearchID string
	Source     Source
}

//Source shows the path from which the object came from
type Source struct {
	Type string   `json:"type"`
	Path []string `json:"path"`
}

//Patient email, phone, address, personal info, cardio, GP
type Patient struct {
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	Lname      string `json:"last_name"`
	Fname      string `json:"first_name"`
	Sex        int64  `json:"sex"`
	DOB        string `json:"dob"`
	SOURCE     Source `json:"source"`
}

//Email for each unique email take date, ptid,mrn,research_id
type Email struct {
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	Email      string `json:"email"`
	Date       string `json:"date"`
	SOURCE     Source `json:"source"`
}

//Phone struct for each unique work or home phone
type Phone struct {
	Type       string `json:"phone_type"`
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	PhoneNum   string `json:"phone_number"`
	Date       string `json:"date"`
	SOURCE     Source `json:"source"`
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
	Date       string `json:"date"`
	SOURCE     Source `json:"source"`
}

//Cardio for each unique cardiologist
type Cardio struct {
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	Cardio     string `json:"cardiologist"`
	Date       string `json:"date"`
	SOURCE     Source `json:"source"`
}

//GP struct for each unique general_practitioner
type GP struct {
	PTID       string `json:"ptid"`
	MRN        string `json:"mrn"`
	ResearchID string `json:"research_id"`
	GP         string `json:"general_practitioner"`
	Date       string `json:"date"`
	SOURCE     Source `json:"source"`
}

//iteratePTID connects to the PatientInfo Database and creates JSON file for each of the event collections
//determines the PTID list from the Contact table
//event collections for each unique PTID are created.
func iteratePTID(path string, fields []string, fldr string) {
	conn := connectToDB(path)
	pList := ptidList(conn)
	query := queryGenerator(fields)

	jsonPath := fldr + "\\email.json"
	accessHelper.CreateFile(jsonPath)
	emailJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer emailJSONFile.Close()

	jsonPath = fldr + "\\phone.json"
	accessHelper.CreateFile(jsonPath)
	phoneJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer phoneJSONFile.Close()

	jsonPath = fldr + "\\address.json"
	accessHelper.CreateFile(jsonPath)
	addressJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer addressJSONFile.Close()

	jsonPath = fldr + "\\GP.json"
	accessHelper.CreateFile(jsonPath)
	GPJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer GPJSONFile.Close()

	jsonPath = fldr + "\\cardio.json"
	accessHelper.CreateFile(jsonPath)
	cardioJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer cardioJSONFile.Close()

	jsonPath = fldr + "\\patient.json"
	accessHelper.CreateFile(jsonPath)
	patientJSONFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer patientJSONFile.Close()

	for _, p := range pList {
		ptidInfoSlice := queryTable(conn, query, p)
		patientInfoSlice := queryTable(conn, "SELECT PTID,DATEOR,LNAME,FNAME,CHART,SEX FROM VR WHERE PTID=?", p)
		vrInfoSlice := queryTable(conn, "SELECT PTID,DATEOR,LNAME,FNAME,CHART,SEX,EMAIL,PHONEHOME,PHONEWORK,GP1,GP2,CARDIO1,CARDIO2,STREET,CITY,PROVINCE,POSTCODE FROM VR WHERE PTID=?", p)
		createJSON(ptidInfoSlice, emailJSONFile, phoneJSONFile, addressJSONFile, GPJSONFile, cardioJSONFile)

		if len(patientInfoSlice) == 0 {
			log.Printf("%s %s returned an empty query", "SELECT PTID,DATEOR,LNAME,FNAME,CHART,SEX FROM VR WHERE PTID=", p)
			continue
		}
		if len(vrInfoSlice) == 0 {
			log.Printf("%s %s returned an empty query", "SELECT PTID,DATEOR,LNAME,FNAME,CHART,SEX, EMAIL, PHONEHOME, PHONEWORK, GP1, GP2, CARDIO1, CARDIO2, STREET, CITY, PROVINCE, POSTCODE FROM VR WHERE PTID=", p)
			continue
		}
		createPatientJSON(patientInfoSlice, patientJSONFile)
		createVRJSON(vrInfoSlice, emailJSONFile, phoneJSONFile, addressJSONFile, GPJSONFile, cardioJSONFile)

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

//ptidList returns a list of unique PTIDs from a connected DB with Contact table
func ptidList(conn *sql.DB) []string {
	var ptidList []string
	rows, err := conn.Query("SELECT PTID FROM [Contact] GROUP BY PTID")
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var r sql.NullString
		err = rows.Scan(&r)
		if err != nil {
			log.Println(err)
		}
		if r.String != "" {
			ptidList = append(ptidList, r.String)
		}
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
	fieldString += " FROM [Contact] WHERE PTID=?"
	return fieldString
}

//queryTable queries a table in the database and maps them into a slice of a slice of OrderedMap
//each slice corresponds to a row and the fields in the row are contained in a slice of OrderedMap presented as {"Field":"Value"}
func queryTable(conn *sql.DB, query string, ptid string) [][]accessHelper.OrderedMap {
	rows, err := conn.Query(query, ptid)
	if err != nil {
		log.Println(err)

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

//compareObjects returns the index of the duplicate with the later(furthest from present) date
func compareObjects(objects []Object) int {
	for i := 0; i < len(objects); i++ {
		for j := i + 1; j < len(objects); j++ {
			if j != i && objects[i].Field == objects[j].Field {
				// if compareDates(objects[i].Date, objects[j].Date) {
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
	if d1 == "" {
		return false
	} else if d2 == "" {
		return true
	}
	date1 := strings.Split(d1, "-")
	date2 := strings.Split(d2, "-")

	for i := 0; i < 3; i++ {
		if excelHelper.StringToInt(date1[i], i, "") > excelHelper.StringToInt(date2[i], i, "") {
			return true
		}
	}
	return false
}

//fixDates finds the date field within the OrderedMap and converts date to YYYY-MM-DD format
func fixDate(oMS [][]accessHelper.OrderedMap, field string) [][]accessHelper.OrderedMap {
	for i := 0; i < len(oMS); i++ {
		for j := 0; j < len(oMS[i]); j++ {
			if oMS[i][j].Colname == field {
				date := strings.Split(oMS[i][j].Value, "T")
				oMS[i][j].Value = date[0]
			}
		}
	}
	return oMS
}

//createPatientJSON creates the documents for the patient event
func createPatientJSON(pMS [][]accessHelper.OrderedMap, patientJSONFile *os.File) {
	pMS = fixDate(pMS, "DATEOR")
	patientSlice := cleanPatient(pMS)
	p := toPatient(patientSlice)
	writeJSON(p, patientJSONFile)
}

//cleanPatient returns the most recent patient row
func cleanPatient(pMS [][]accessHelper.OrderedMap) []accessHelper.OrderedMap {
	var d []PDate
	for i := 0; i < len(pMS); i++ {
		for j := 0; j < len(pMS[i]); j++ {
			if checkDate(pMS[i][j].Colname, "DATEOR") {
				var date PDate
				date.Date = pMS[i][j].Value
				date.Row = i
				d = append(d, date)
			}
		}
	}
	recent := 0
	for i := 1; i < len(d); i++ {
		if !compareDates(d[recent].Date, d[i].Date) {
			recent = i

		}
	}
	return pMS[recent]
}

//toPatient assigns slice to Patient stuct
func toPatient(pS []accessHelper.OrderedMap) Patient {
	var p Patient
	for i := range pS {
		if pS[i].Colname == "PTID" {
			p.PTID = pS[i].Value
		}
		if pS[i].Colname == "CHART" {
			p.MRN = pS[i].Value
		}
		if pS[i].Colname == "LNAME" {
			p.Lname = pS[i].Value
		}
		if pS[i].Colname == "FNAME" {
			p.Fname = pS[i].Value
		}
		if pS[i].Colname == "SEX" {
			p.Sex = excelHelper.StringToInt(pS[i].Value, 0, "")
		}
	}
	p.DOB = toDate(p.PTID)
	return p
}

//toDate takes a PTID and parses the DOB and returns it in the YYYY-MM-DD format
func toDate(dob string) string {
	nDob := []byte(dob)
	n := len(nDob)
	year := "19" + string(nDob[n-2]) + string(nDob[n-1])
	month := string(nDob[n-4]) + string(nDob[n-3])
	day := string(nDob[n-6]) + string(nDob[n-5])
	return year + "-" + month + "-" + day
}

//createVRJSON creates the documents for non-patient events from the VR table
func createVRJSON(oMS [][]accessHelper.OrderedMap, emailJSONFile *os.File, phoneJSONFile *os.File, addressJSONFile *os.File, gpJSONFile *os.File, cardioJSONFile *os.File) {
	oMS = fixDate(oMS, "DATEOR")
	emailObject(oMS, emailJSONFile, "DATEOR")
	phoneObject(oMS, phoneJSONFile, "DATEOR")
	gpObject(oMS, gpJSONFile, "DATEOR")
	cardioObject(oMS, cardioJSONFile, "DATEOR")
	addressObject(oMS, addressJSONFile, "DATEOR")
}

//createJSON creates documents for non-patient events from the contact table
func createJSON(oMS [][]accessHelper.OrderedMap, emailJSONFile *os.File, phoneJSONFile *os.File, addressJSONFile *os.File, gpJSONFile *os.File, cardioJSONFile *os.File) {
	oMS = fixDate(oMS, "FU_D")
	oMS = fixDate(oMS, "LKA_D")

	emailObject(oMS, emailJSONFile, "FU_D")
	phoneObject(oMS, phoneJSONFile, "FU_D")
	gpObject(oMS, gpJSONFile, "FU_D")
	cardioObject(oMS, cardioJSONFile, "FU_D")
	addressObject(oMS, addressJSONFile, "FU_D")
}

//emailObject creates email objects
func emailObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File, date string) {
	emails := cleanObject(oMS, "EMAIL", date)
	basic := getBasic(oMS)
	createEmail(emails, basic, jsonFile, "followupVR")
}

//phoneObject creates phone objects one for HOME and another for WORK
func phoneObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File, date string) {
	basic := getBasic(oMS)
	phonesHome := cleanObject(oMS, "PHONEHOME", date)
	createPhones(phonesHome, basic, jsonFile, true, "followupVR")
	phonesWork := cleanObject(oMS, "PHONEWORK", date)
	createPhones(phonesWork, basic, jsonFile, false, "followupVR")
}

//gpObject creates gp objects one for GP1 and another for GP2
func gpObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File, date string) {
	basic := getBasic(oMS)
	gp1 := cleanObject(oMS, "GP1", date)
	createGP(gp1, basic, jsonFile, "followupVR")
	gp2 := cleanObject(oMS, "GP2", date)
	createGP(gp2, basic, jsonFile, "followupVR")
}

//cardioObject creates cardio objects for Cardio1 and another for Cardio2
func cardioObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File, date string) {
	basic := getBasic(oMS)
	cardio1 := cleanObject(oMS, "CARDIO1", date)
	createCardio(cardio1, basic, jsonFile, "followupVR")
	cardio2 := cleanObject(oMS, "CARDIO2", date)
	createCardio(cardio2, basic, jsonFile, "followupVR")
}

//addressObject creates an address object when any of the address related fields are populated
func addressObject(oMS [][]accessHelper.OrderedMap, jsonFile *os.File, date string) {
	basic := getBasic(oMS)
	address := cleanAddress(oMS, date)
	createAddress(address, basic, jsonFile, "followupVR")
}

//getBasic gets the basic information(PTID,MRN) and stores it in a Basic struct
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

//checkDate checks if the field is a date field
func checkDate(c string, field string) bool {
	if c == field {
		return true
	}
	return false
}

//cleanObject reads the OrderedMap for all the populated cells which correspond to the particular event
//returns a slice of Object that contains the duplicate free field,date,and path
func cleanObject(oMS [][]accessHelper.OrderedMap, field string, dateField string) []Object {
	var fieldValue string
	var date string
	var objects []Object
	var path []string
	validField := false

	for i := 0; i < len(oMS); i++ {
		validField = false
		path = path[:0]
		for j := 0; j < len(oMS[i]); j++ {
			if oMS[i][j].Colname == field {
				if oMS[i][j].Value != "" {
					fieldValue = oMS[i][j].Value
					validField = true
				}
			}
			if checkDate(oMS[i][j].Colname, dateField) {
				date = oMS[i][j].Value
			}

		}
		if validField && date == "" && dateField == "FU_D" {
			for j := 0; j < len(oMS[i]); j++ {
				if oMS[i][j].Colname == "LKA_D" {
					date = oMS[i][j].Value
				}
			}
		}
		if validField {
			path = path[:0]
			for j := 0; j < len(oMS[i]); j++ {
				if oMS[i][j].Colname == "PATH" {
					path = append(path, oMS[i][j].Value)
				}
			}
		}

		if validField {
			o := Object{fieldValue, date, path}
			objects = append(objects, o)
		}
	}

	duplicates := 0
	for true {
		duplicates = compareObjects(objects)
		if duplicates != -1 {
			objects[duplicates].Path = append(objects[duplicates].Path, objects[duplicates].Path...)
			objects = append(objects[:duplicates], objects[duplicates+1:]...)
		} else {
			break
		}
	}
	return objects
}

//cleanAddress reads the OrderedMap for all the populated cells which correspond to the address event
//returns a slice of Object that contains the duplicate free field,date,and path
func cleanAddress(oMS [][]accessHelper.OrderedMap, dateField string) []Object {
	var fieldValue []string
	var date string
	var objects []Object
	var path []string

	for i := 0; i < len(oMS); i++ {
		validField := ""
		fieldValue = []string{}
		path = path[:0]

		for j := 0; j < len(oMS[i]); j++ {
			if oMS[i][j].Colname == "STREET" {
				if oMS[i][j].Value != "" {
					fieldValue = append(fieldValue, oMS[i][j].Value)
					validField += "a"
				} else {
					fieldValue = append(fieldValue, " ")
				}
			}
			if oMS[i][j].Colname == "CITY" {
				if oMS[i][j].Value != "" {
					fieldValue = append(fieldValue, oMS[i][j].Value)
					validField += "b"
				} else {
					fieldValue = append(fieldValue, " ")
				}
			}
			if oMS[i][j].Colname == "PROVINCE" {
				if oMS[i][j].Value != "" {
					fieldValue = append(fieldValue, oMS[i][j].Value)
					validField += "c"
				} else {
					fieldValue = append(fieldValue, " ")
				}
			}
			if oMS[i][j].Colname == "POSTCODE" {
				if oMS[i][j].Value != "" {
					fieldValue = append(fieldValue, oMS[i][j].Value)
					validField += "d"
				} else {
					fieldValue = append(fieldValue, " ")
				}
			}
			if checkDate(oMS[i][j].Colname, dateField) {
				date = oMS[i][j].Value
			}
			if validField != "" {

				path = path[:0]

				for j := 0; j < len(oMS[i]); j++ {
					if oMS[i][j].Colname == "PATH" {
						path = append(path, oMS[i][j].Value)
					}
				}
			}
		}
		if validField != "" {
			var fieldString string
			for _, f := range fieldValue {
				fieldString += f + "_"
			}
			fieldString = strings.TrimSuffix(fieldString, "_")
			o := Object{fieldString, date, path}
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

//createPhones generates and writes phone events to JSON
func createPhones(phones []Object, basic Basic, jsonFile *os.File, home bool, s string) {
	var e Phone
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	e.SOURCE.Type = s
	if home {
		e.Type = "home"
	} else {
		e.Type = "work"
	}
	for i := range phones {
		e.PhoneNum = phones[i].Field
		e.Date = phones[i].Date
		e.SOURCE.Path = phones[i].Path
		writeJSON(e, jsonFile)
	}
}

//createGP generates and writes GP events to JSON
func createGP(gps []Object, basic Basic, jsonFile *os.File, s string) {
	var e GP
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	e.SOURCE = basic.Source
	e.SOURCE.Type = s

	for i := range gps {
		e.GP = gps[i].Field
		e.Date = gps[i].Date
		e.SOURCE.Path = gps[i].Path

		writeJSON(e, jsonFile)
	}
}

//createCardio generates and writes Cardio events to JSON
func createCardio(cardio []Object, basic Basic, jsonFile *os.File, s string) {
	var e Cardio
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	e.SOURCE.Type = s
	for i := range cardio {
		e.Cardio = cardio[i].Field
		e.Date = cardio[i].Date
		e.SOURCE.Path = cardio[i].Path

		writeJSON(e, jsonFile)
	}
}

//createEmail generates and writes emails to JSON
func createEmail(emails []Object, basic Basic, jsonFile *os.File, s string) {
	var e Email
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	e.SOURCE.Type = s

	for i := range emails {
		e.Email = emails[i].Field
		e.Date = emails[i].Date
		e.SOURCE.Path = emails[i].Path
		writeJSON(e, jsonFile)
	}
}

//createAddress generates and writes addresses to JSON
func createAddress(addr []Object, basic Basic, jsonFile *os.File, s string) {
	var e Address
	e.PTID = basic.PTID
	e.MRN = basic.MRN
	e.SOURCE.Type = s

	for i := range addr {
		a := strings.Split(addr[i].Field, "_")
		e.Street = a[0]
		e.City = a[1]
		e.Province = a[2]
		e.PostCode = a[3]
		e.Date = addr[i].Date
		e.SOURCE.Path = addr[i].Path
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

func main() {
	errFile := accessHelper.CreateErrorLog("C:\\Users\\ext_hsc\\Desktop\\jsonError.log")
	log.SetOutput(errFile)
	defer errFile.Close()
	path := "C:\\Users\\ext_hsc\\Documents\\valve_registry_PHI\\ContactInfo.accdb"
	fields := []string{"PTID", "CHART", "EMAIL", "FU_D", "PHONEHOME", "PHONEWORK", "GP1", "GP2", "CARDIO1", "CARDIO2", "STREET", "CITY", "PROVINCE", "POSTCODE", "PATH", "LKA_D"}
	fldr := "C:\\Users\\ext_hsc\\Desktop\\JSON"
	iteratePTID(path, fields, fldr)

}
