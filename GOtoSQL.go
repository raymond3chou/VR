package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	_ "odbc/driver"
)

func selectAccess(conn *sql.DB, file *os.File) bool {

	rows, err := conn.Query("SELECT * from ContactInfo")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	//queried Oringinal DB for ***
	for rows.Next() {
		var (
			ptid   string
			chart  string
			lname  string
			fname  string
			sex    string
			age    string
			dob    string
			street string
			city   string
			prov   string
			pcode  string
			hnum   string
			pnum   string
			email  string
		)
		err = rows.Scan(&ptid, &chart, &lname, &fname, &sex, &age, &dob, &street, &city, &prov, &pcode, &hnum, &pnum, &email)
		if err != nil {
			fmt.Println("Select Row Failed")
			return false
		}
		s := strings.Split(dob, "T")
		row := "\n" + ptid + "|" + chart + "|" + lname + "|" + fname + "|" + sex + "|" + age + "|" + s[0] + "|" + street + "|" + city + "|" + prov + "|" + pcode + "|" + hnum + "|" + pnum + "|" + email
		fmt.Println("Here")
		fmt.Println("Here")

		fmt.Println(row)
		fileWrite(file, row)
	}
	//iterate through each row of the executed Query from Originating DB
	err = rows.Err()
	if err != nil {
		fmt.Println("Select Failed")
		return false
	}
	// //flag errors from querying Oringinating DB
	return true
}

func fileWrite(file *os.File, row string) {

	_, err := io.WriteString(file, row)
	if err != nil {
		fmt.Println("Could Not Write String")
	}
}

func main() {

	conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq=.\\TestDB.accdb")
	if err != nil {
		fmt.Println("Connecting Error")
		return
	}
	defer conn.Close()
	//Originating Database connection established
	file, err := os.OpenFile("C:\\Users\\raymond chou\\Desktop\\ContactInfo.txt", os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Could not open text file")
	}
	defer file.Close()
	selectAccess(conn, file)

}

// rec, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq=.\\TestDBRec.accdb")
// if err != nil {
// 	fmt.Println("Connecting Error")
// 	return
// }
// defer rec.Close()
// // Receiving Database connection established
//**************************************************************
// INSERT INTO MS ACCESS USING sql
// CURRENTLY DOES NOT WORK FOR SOME REASON
//
// func insertAccess(conn *sql.DB) string {
// 	rslt, err := conn.Exec("INSERT INTO info VALUES(?)", "dolly")
// 	if err != nil {
// 		log.Fatalln("Could Not Insert")
// 	}
// 	// preparing inset statement for receiving DB
// 	fmt.Println(rslt.LastInsertId())
// 	return "Insert Complete"
// }
//**************************************************************
