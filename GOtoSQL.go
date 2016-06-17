package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/alexbrainman/odbc"
)

func selectAccess(conn *sql.DB, file *os.File) int {
	//Queries the database, and passes the row to be appended to the test file
	inserted := 0
	//count number inserted
	rows, err := conn.Query("SELECT PTID, CHART, LName, FName, SEX, AGE, DOB, STREET, CITY, PROV, PCODE, H_NUM, PNUM, Email FROM ContactInfo")
	if err != nil {
		fmt.Println("Query Failed")
		return 0
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
			return inserted
		}
		s := strings.Split(dob, "T")
		row := "\n" + ptid + "|" + chart + "|" + lname + "|" + fname + "|" + sex + "|" + age + "|" + s[0] + "|" + street + "|" + city + "|" + prov + "|" + pcode + "|" + hnum + "|" + pnum + "|" + email

		inserted += fileWrite(file, row)
	}
	//iterate through each row of the executed Query from Originating DB
	err = rows.Err()
	if err != nil {
		fmt.Println("Select Failed")
		return inserted
	}
	// //flag errors from querying Oringinating DB
	return inserted
}

func fileWrite(file *os.File, row string) int {
	//Writes the queried row into a text file
	_, err := io.WriteString(file, row)
	if err != nil {
		fmt.Println("Could Not Write String")
		return 0
	}
	return 1
}

func findDB(dir string) (filename string) {
	//Go through the current directory and identifies folders, .accdb, and .mdb
	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		if strings.Contains(f.Name(), ".accdb") {
			fmt.Println(f.Name())
		} else if strings.Contains(f.Name(), ".mdb") {
			fmt.Println(f.Name())
		} else if !(strings.Contains(f.Name(), ".")) {
			fmt.Println(f.Name())
		}
	}
}

func main() {

	//Iterate through Files, finds folders, .accdb,and .mdb
	//Check if file is a followup/peri op**
	//if so, parse the necessary fields**

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
	//Opens text file that can be constantly appended to. ONLY NEEDS TO BE CALLED ONCE
	inserted := selectAccess(conn, file)
	fmt.Printf("Total Number of Rows Read= %d\n", inserted)
}
