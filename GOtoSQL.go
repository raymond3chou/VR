package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/miquella"
)

func selectAccess(conn *sql.DB, file *os.File, tablename string) int {
	//Queries the database, and passes the row to be appended to the test file
	inserted := 0
	//count number inserted
	rows, err := conn.Query("SELECT PTID, CHART, LName, FName, SEX, AGE, DOB, STREET, CITY, PROV, PCODE, H_NUM, PNUM, Email FROM " + tablename)
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

func findDB(dir string) ([]string, []string, []string) {
	//Go through the current directory and identifies folders, .accdb, and .mdb
	var mdbnames []string
	var accdbnames []string
	var foldernames []string

	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if strings.Contains(f.Name(), ".accdb") {
			mdbnames = append(mdbnames, f.Name())
		} else if strings.Contains(f.Name(), ".mdb") {
			accdbnames = append(accdbnames, f.Name())
		} else if !(strings.Contains(f.Name(), ".")) {
			foldernames = append(foldernames, f.Name())
		}
	}
	return mdbnames, accdbnames, foldernames
}

func findtable(conn *sql.DB) []string {
	//Iterates through all the tables in the database.
	//Currently only works with .mdb. .accdb does not have permission.
	var tablenames []string
	rows, err := conn.Query("SELECT Name FROM MSysObjects WHERE Type=1 AND Flags=0;")
	if err != nil {
		fmt.Println(err)
		return "Query Failed"
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			fmt.Println("Select Row Failed")
		}
		tablenames = append(tablenames, table)
	}
	return tablenames
}

func cleantable(tablenames []string, match string) {
	var PHItablenames []string
	for _, tablename := range tablenames {
		if strings.Contains(tablename, match) {
			PHItablenames = append(PHItablenames, tablename)
		}
	}
	return PHItablenames
}

func connectandexecute(dir string, mdbnames []string) string {

	for _, mdbname := range mdbnames {
		dbq := dir + mdbname
		conn, err := sql.Open("mgodbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq="+dbq)
		if err != nil {
			fmt.Println("Connecting Error")
			return "Failed"
		}
		//Originating Database connection established
		tablenames := cleantable(findtable(conn), "FU")

		file, err := os.OpenFile("C:\\Users\\raymond chou\\Desktop\\ContactInfo.txt", os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("Could not open text file")
			return "Failed"
		}

		//Opens text file that can be constantly appended to. ONLY NEEDS TO BE CALLED ONCE

		tablenames := findtable(conn)
		for _, tablename := range tablenames {
			inserted := selectAccess(conn, file, tablename)
			fmt.Printf("Total Number of Rows Read= %d\n", inserted)
		}
		conn.Close()
		file.Close()
	}
	return "Files Closed and Function Executed"
}

func dbPresent(dir string, mdbnames []string, accdbnames []string) {
	//Checks if there are mdbnames or accdbnames and then executes the code to connect to the DB
	var result string

	if len(mdbnames) != 0 {
		connectandexecute(dir, mdbnames)
		result += "mdb Executed"
	} else {
		result += "mdb Empty"
	}

	if len(accdbnames) != 0 {
		connectandexecute(dir, accdbnames)
		result += "accdb Executed"
	} else {
		result += "accdb Empty"
	}
	return result
}

func gothroughfolder(foldernames []string, dir string) bool {
	if len(foldernames) != 0 {
		for _, foldername := range foldernames {
			dir = "/" + foldername
			mdbnames, accdbnames, newfoldernames := findDB(dir)
			dbPresent(dir, mdbnames, accdbnames)
			gothroughfolder(newfoldernames, dir)
		}
	} else {
		return false
	}

}
func main() {

	dir := "./"
	mdbnames, accdbnames, foldernames := findDB(dir)
	dbPresent(dir, mdbnames, accdbnames)
	gothroughfolder(foldernames, dir)

}
