package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/alexbrainman/odbc"
)

var ( //Global Variables to Track files accessed
	numtables    int
	numfiles     int
	numfolders   int
	rowsinserted int
)

//Sort the Columns
//fill the empty or not present columns with empty cells
//Error handling ie writing to a new .txt
func errorWrite(issue string) {
	file, conn := connectToTxt("C:\\Users\\raymond chou\\Desktop\\ErrorLog.txt")
	if conn {
		fmt.Println("Error File Opened")
	} else {
		fmt.Println("Unable to Open Error File")
		return
	}
	fileWrite(file, issue+"\n")
	file.Close()
}

func checkFollowup(conn *sql.DB, tablename string) (bool, []string, string) {
	FU := false
	//Attempts to Run the Query
	rows, err := conn.Query("SELECT * FROM [" + tablename + "]")
	if err != nil {
		issue := "Cant Run SELECT * FROM [" + tablename + "]"
		fmt.Println(err)
		errorWrite(issue)
	}
	//Returns the Columns from the QUERY
	column, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		errorWrite(err.Error())
	}

	maincolumns := []string{"PTID", "CHART", "LNAME", "FNAME", "SEX", "AGE", "STREET", "CITY", "PROVINCE", "POSTCODE", "PHONEHOME", "PHONEWORK", "PHONECELL", "EMAIL", "DOB"}
	var query string
	for _, columnname := range column {
		for _, maincolname := range maincolumns {
			if columnname == "FU_D" {
				FU = true
			}
			if columnname == "DIED" {
				FU = true
			}
			if columnname == "DTH_D" {
				FU = true
			}
			if strings.Contains(columnname, maincolname) {
				query += " " + columnname
			}
		}
	}
	return FU, maincolumns, query
}

func convertToString(vals []interface{}) []string {
	row := make([]string, len(vals))
	for i, val := range vals {
		value := val.(*sql.NullString)
		row[i] = value.String
	}
	return row
}

func convertToText(maincolumns []string, cols map[string]string) string {
	//takes in the queried row divided in an array of strings based off of the column
	//maincolumns contains the master columns and a flag for which ever one was used
	//the function arranges based on
	var row string
	found := false
	row = "\n"
	for _, mastercol := range maincolumns {
		for colname := range cols {
			if strings.Contains(colname, mastercol) {
				row += cols[colname] + "|"
				found = true
			}
		}
		if !found {
			row += " |"
		}
	}
	row = strings.TrimSuffix(row, "|")
	return row
}

func convertToMap(cols map[string]string, rowstring []string) map[string]string {
	endindex := len(rowstring)
	i := 0
	for key := range cols {
		if i < endindex {
			cols[key] = rowstring[i]
			i++
		} else {
			break
		}

	}

	return cols
}

func selectAccess(conn *sql.DB, file *os.File, tablename string) (int, int) {
	//Queries the database, and passes the row to be appended to the test file
	//returns # inserted
	var selectquery string
	var query string

	FU, maincolumns, query := checkFollowup(conn, tablename)
	if FU {
		selectquery = "SELECT" + query + " FROM [" + tablename + "]"
	} else {
		issue := tablename + " is not a Follow Up Table"
		errorWrite(issue)
		return 0, 0
	}
	NumberofRows := 0
	inserted := 0
	//count number inserted
	rows, err := conn.Query(selectquery)
	if err != nil {
		issue := "Select query failed to execute " + tablename
		fmt.Println("Select query failed to execute " + tablename)
		fmt.Println(err)
		errorWrite(issue)
		return 0, 0
	}
	defer rows.Close()
	//queried Oringinal DB for ***
	queriedcols, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
	}
	colsmap := make(map[string]string)
	for _, colname := range queriedcols {
		colsmap[colname] = ""
	}

	vals := make([]interface{}, len(queriedcols))
	for i := range queriedcols {
		vals[i] = new(sql.NullString)
	}

	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			issue := "Read row failed. Not enough parameters in: " + tablename
			fmt.Println("Read row failed. Not enough parameters?")
			fmt.Println(err)
			errorWrite(issue)
			return inserted, NumberofRows
		}
		NumberofRows++
		rowstring := convertToString(vals)
		cols := convertToMap(colsmap, rowstring)
		row := convertToText(maincolumns, cols)
		inserted += fileWrite(file, row)
	}
	//iterate through each row of the executed Query from Originating DB
	err = rows.Err()
	if err != nil {
		fmt.Println("rows.Err failure")
		errorWrite(err.Error())
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
	var phitablenames []string
	for _, tablename := range tablenames {
		if strings.Contains(tablename, match) {
			phitablenames = append(phitablenames, tablename)
		}
	}
	return phitablenames
}

func connectToDB(dir string, dbname string) *sql.DB {

	dbq := dir + "/" + dbname
	fmt.Println("Connecting to " + dbq)
	conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq="+dbq)
	if err != nil {
		log.Fatal("Connection to " + dbq + " Failed")
	}
	fmt.Println("Connected to " + dbq)
	return conn
}

func connectToTxt(filedir string) (*os.File, bool) {

	file, err := os.OpenFile(filedir, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Unable to Open Text File: %s", filedir)
		fmt.Print(err)
		errorWrite(err.Error())
		return file, false
	}
	return file, true
}

func connectandexecute(dir string, dbnames []string) string {
	//Connects to Database and File.
	//Calls matchTable Function and then iterates through the tables using SelectAccess
	var dbaccessed int
	for _, dbname := range dbnames {
		conn := connectToDB(dir, dbname)
		if conn != nil {
			dbaccessed++
		} else {
			continue
		}
		//Originating Database connection established

		tablenames := findtable(conn)
		tablelength := len(tablenames)
		fmt.Printf("%d Tables in %s\n", tablelength, dbname)

		if len(tablenames) != 0 {
			tablenames = matchTable(tablenames, "FU")
			fmt.Printf("%d/%d Tables Match the Criteria\n", len(tablenames), tablelength)
		} else {
			return "No Table Names"
		}

		file, connection := connectToTxt("C:\\Users\\raymond chou\\Desktop\\ContactInfo.txt")
		if !connection {
			continue
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
		fmt.Printf("%s Closed and %d Table(s) Extracted\n ---------------------------\n", dbname, tableused)
		numtables += tableused
	}
	result := "File(s) Accessed and Closed: " + strconv.Itoa(dbaccessed)
	numfiles += dbaccessed
	return result
}

func dbPresent(dir string, mdbnames []string, accdbnames []string) string {
	//Checks if there are mdbnames or accdbnames and then executes the code to connect to the DB
	var result string

	if len(mdbnames) != 0 {
		results := connectandexecute(dir, mdbnames)
		fmt.Println(results + "\n*************************\n")
		result += ".mdb Files Executed"
	} else {
		result += ".mdb Files Empty"
	}

	if len(accdbnames) != 0 {
		fmt.Println(connectandexecute(dir, accdbnames) + "\n*************************\n")
		result += " .accdb Files Executed"
	} else {
		result += " .accdb Files Empty"
	}
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
	start := time.Now()
	fmt.Println("\n\n------START OF PROGRAM------")

	dir := "./"
	foldernames := []string{""}
	status := gothroughfolder(foldernames, dir)

	if status == true {
		elapsed := time.Since(start)
		fmt.Printf("\n------COMPLETE------\nFolder(s) Accessed: %d\nFile(s) Accessed: %d\nTable(s) Accessed: %d\nRow(s) Inserted: %d\nTime Taken: %s", numfolders, numfiles, numtables, rowsinserted, elapsed)
	}
}
