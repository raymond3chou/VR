package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	_ "github.com/miquella"
)

var (
	numtables    int
	numfiles     int
	numfolders   int
	rowsinserted int
)

func selectAccess(conn *sql.DB, file *os.File, tablename string) (int, int) {
	//Queries the database, and passes the row to be appended to the test file
	//returns # inserted
	NumberofRows := 0
	inserted := 0
	//count number inserted
	rows, err := conn.Query("SELECT PTID, CHART, LNAME, FNAME, SEX, AGE, STREET, CITY, PROVINCE, POSTCODE, PHONEHOME,PHONEWORK,PHONECELL, EMAIL FROM " + tablename)
	if err != nil {
		fmt.Println("Select query failed to execute " + tablename)
		fmt.Println(err)
		return 0, 0
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
			street string
			city   string
			prov   string
			pcode  string
			hnum   string
			wnum   string
			cnum   string
			email  sql.NullString //accounts for NULL entry
		)
		err = rows.Scan(&ptid, &chart, &lname, &fname, &sex, &age, &street, &city, &prov, &pcode, &hnum, &wnum, &cnum, &email)
		if err != nil {
			fmt.Println("Read row failed. Not enough parameters?")
			fmt.Println(err)
			return inserted, NumberofRows
		}
		NumberofRows++

		// s := strings.Split(dob, "T")
		row := "\n" + ptid + "|" + chart + "|" + lname + "|" + fname + "|" + sex + "|" + age + "|" + street + "|" + city + "|" + prov + "|" + pcode + "|" + hnum + "|" + wnum + "|" + cnum + "|" + email.String

		inserted += fileWrite(file, row)
	}
	//iterate through each row of the executed Query from Originating DB
	err = rows.Err()
	if err != nil {
		fmt.Println("rows.Err failure")
		return inserted, NumberofRows
	}
	// //flag errors from querying Oringinating DB
	return inserted, NumberofRows
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

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".accdb") {
			accdbnames = append(accdbnames, f.Name())
		} else if strings.Contains(f.Name(), ".mdb") {
			mdbnames = append(mdbnames, f.Name())
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
		fmt.Println("Failed to Select Tablenames")
		fmt.Println(err)
		return tablenames
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			fmt.Println("Failed to pass tablenames")
		}
		tablenames = append(tablenames, table)
	}
	return tablenames
}

func matchTable(tablenames []string, match string) []string {
	//finds tables with names that matches the match string and then returns a slice of the tables
	var PHItablenames []string
	for _, tablename := range tablenames {
		if strings.Contains(tablename, match) {
			PHItablenames = append(PHItablenames, tablename)
		}
	}
	return PHItablenames
}

func connectandexecute(dir string, dbnames []string) string {
	//Connects to Database and File.
	//Calls matchTable Function and then iterates through the tables using SelectAccess
	var dbaccessed int
	for _, dbname := range dbnames {
		dbq := dir + "/" + dbname
		fmt.Println("Connecting to " + dbq)
		conn, err := sql.Open("mgodbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq="+dbq)
		if err != nil {
			return "Connecting Error"
		}
		fmt.Println("Connected to " + dbq)
		dbaccessed++
		//Originating Database connection established

		tablenames := findtable(conn)
		tablelength := len(tablenames)
		fmt.Printf("%d Tables in %s\n", tablelength, dbname)

		if len(tablenames) != 0 {
			tablenames = matchTable(tablenames, "Info")
			fmt.Printf("%d/%d Tables Match the Criteria\n", len(tablenames), tablelength)
		} else {
			return "No Table Names"
		}
		file, err := os.OpenFile("C:\\Users\\raymond chou\\Desktop\\ContactInfo.txt", os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return "Could not open text file"
		}
		//Opens text file that can be constantly appended to. ONLY NEEDS TO BE CALLED ONCE
		var tableused int
		for _, tablename := range tablenames {
			inserted, NumberofRows := selectAccess(conn, file, tablename)
			fmt.Printf("Total Number of Rows Read and Inserted from %s = %d/%d\n", tablename, inserted, NumberofRows)
			tableused++
			rowsinserted += inserted
		}
		conn.Close()
		file.Close()
		fmt.Printf("%s Closed and %d Table(s) Extracted\n", dbname, tableused)
		numtables += tableused
	}
	result := strconv.Itoa(dbaccessed) + " Files Accessed and Closed"
	numfiles += dbaccessed
	return result
}

func dbPresent(dir string, mdbnames []string, accdbnames []string) string {
	//Checks if there are mdbnames or accdbnames and then executes the code to connect to the DB
	var result string

	if len(mdbnames) != 0 {
		results := connectandexecute(dir, mdbnames)
		fmt.Println(results)
		result += ".mdb Files Executed"
	} else {
		result += ".mdb Files Empty"
	}
	//
	// if len(accdbnames) != 0 {
	// 	fmt.Println(connectandexecute(dir, accdbnames))
	// 	result += "accdb Executed"
	// } else {
	// 	result += "accdb Empty"
	// }
	return result
}

func gothroughfolder(foldernames []string, dir string) bool {
	if len(foldernames) != 0 {
		for _, foldername := range foldernames {
			dir = dir + foldername
			fmt.Printf("\nEntering Folder Directory: %s\n", dir)
			numfolders++
			mdbnames, accdbnames, newfoldernames := findDB(dir)
			printDirInfo(mdbnames, accdbnames, foldernames, dir)
			result := dbPresent(dir, mdbnames, accdbnames)
			fmt.Printf("\n%s \n", result)
			dir += "/"
			gothroughfolder(newfoldernames, dir)
		}
	}
	return true
}

func printDirInfo(mdbnames []string, accdbnames []string, foldernames []string, dir string) {
	//prints out the info for the current folder
	fmt.Printf("Number of .mdb files in %s: %d \nNumber of .accdb files in %s: %d \nNumber of Folders in %s: %d\n\n", dir, len(mdbnames), dir, len(accdbnames), dir, len(foldernames))
}

func main() {

	dir := "./"
	mdbnames, accdbnames, foldernames := findDB(dir)
	printDirInfo(mdbnames, accdbnames, foldernames, dir)
	result := dbPresent(dir, mdbnames, accdbnames)
	fmt.Printf("\n%s \n", result)
	status := gothroughfolder(foldernames, dir)
	if status == true {
		fmt.Printf("Complete:\nFolders Accessed: %d\nFiles Accessed: %d\nTables Accessed: %d\nRows Inserted: %d", numfolders, numfiles, numtables, rowsinserted)
	}

}
