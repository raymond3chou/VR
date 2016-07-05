package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/access"
	_ "github.com/alexbrainman/odbc"
)

var ( //Global Variables to Track files accessed
	numtables    int
	numfiles     int
	numfolders   int
	rowsinserted int
)

//checkFollowup queries the specified table to check if it is a Followup Table(i.e contains FU_D,DIED,DTH_D)
//based on the columns present, it also returns a list of column names matched to a masterlist
//If no column names are present then it returns an empty query
func checkFollowup(conn *sql.DB, tablename string) (bool, []string, string) {
	FU := false
	//Attempts to Run the Query
	rows, err := conn.Query("SELECT * FROM [" + tablename + "]")
	if err != nil {
		issue := "Cant Run SELECT * FROM [" + tablename + "]"
		fmt.Println(err)
		access.ErrorWrite(issue)
	}
	//Returns the Columns from the QUERY
	column, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		access.ErrorWrite(err.Error())
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
				query += " " + columnname + ","
			}
		}
	}
	query = strings.TrimSuffix(query, ",")
	return FU, maincolumns, query
}

//selectAccess
func selectAccess(conn *sql.DB, file *os.File, tablename string) (int, int) {
	//Queries the database, and passes the row to be appended to the test file
	//returns # inserted
	var selectquery string
	var query string

	FU, maincolumns, query := checkFollowup(conn, tablename)
	if FU && query != "" {
		selectquery = "SELECT" + query + " FROM [" + tablename + "]"
	} else if query == "" {
		issue := tablename + " does not contain any columns related to PHI\n"
		access.ErrorWrite(issue)
		return 0, 0
	} else {
		issue := tablename + " is not a Follow Up Table\n"
		access.ErrorWrite(issue)
		return 0, 0
	}
	numberofRows := 0
	inserted := 0
	//count number inserted
	rows, err := conn.Query(selectquery)
	if err != nil {
		issue := "Select query failed to execute " + tablename + "\n"
		fmt.Println("Select query failed to execute " + tablename)
		fmt.Println(err)
		access.ErrorWrite(issue)
		return 0, 0
	}
	defer rows.Close()
	//queried Oringinal DB for ***
	queriedcols, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
	}
	colsOMap := make([]access.OrderedMap, len(queriedcols))
	for i, colname := range queriedcols {
		colsOMap[i].Colname = colname
		colsOMap[i].Value = ""
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
			access.ErrorWrite(issue)
			return inserted, numberofRows
		}
		numberofRows++
		rowstring := access.ConvertToString(vals)
		cols := access.ConvertToOrderedMap(colsOMap, rowstring)
		row := access.ConvertToText(maincolumns, cols)
		inserted += access.FileWrite(file, row)
	}
	//iterate through each row of the executed Query from Originating DB
	err = rows.Err()
	if err != nil {
		fmt.Println("rows.Err failure")
		access.ErrorWrite(err.Error())
		return inserted, numberofRows
	}
	// //flag errors from querying Oringinating DB
	return inserted, numberofRows
}

//findDB goes through the current directory and identifies folders, .accdb, and .mdb
func findDB(dir string) ([]string, []string, []string) {
	var mdbnames []string
	var accdbnames []string
	var foldernames []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		if f.IsDir() {
			foldernames = append(foldernames, f.Name())
		} else {
			if strings.Contains(f.Name(), ".accdb") {
				accdbnames = append(accdbnames, f.Name())
			} else if strings.Contains(f.Name(), ".mdb") {
				mdbnames = append(mdbnames, f.Name())
			}
		}
	}
	return mdbnames, accdbnames, foldernames
}

//findTable iterates through all the tables in the database.
//Currently only works with .mdb. .accdb does not have permission.
func findTable(conn *sql.DB) []string {
	var tablenames []string

	rows, err := conn.Query("SELECT Name FROM MSysObjects WHERE Type=1 AND Flags=0;")
	defer rows.Close()
	if err != nil {
		fmt.Println("Failed to Select Tablenames")
		fmt.Println(err)
		return tablenames
	}

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

//connectToDB Connects to a specified database in a specified directory
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

//connectExecute connects to every database in a specified directory and checks if its is a FU table
//if so, it reads the columns; sorts them; and writes them on to a textfile
func connectExecute(dir string, dbnames []string) string {
	var dbaccessed int
	for _, dbname := range dbnames {
		conn := connectToDB(dir, dbname)
		if conn != nil {
			dbaccessed++
		} else {
			continue
		}
		//Originating Database connection established

		tablenames := findTable(conn)
		tablelength := len(tablenames)
		fmt.Printf("%d Tables in %s\n", tablelength, dbname)

		file, connection := access.ConnectToTxt("C:\\Users\\raymond chou\\Desktop\\ContactInfo.txt")
		if !connection {
			continue
		}
		//Opens text file that can be constantly appended to. ONLY NEEDS TO BE CALLED ONCE
		var tableused int
		for _, tablename := range tablenames {
			inserted, numberofRows := selectAccess(conn, file, tablename)
			fmt.Printf("Total Number of Rows Read and Inserted from %s = %d/%d\n", tablename, inserted, numberofRows)
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

//dbPresent checks if there are mdbnames or accdbnames and then executes the code to connect to the DB
//currently needed because of no permission on .accdb files.
func dbPresent(dir string, mdbnames []string, accdbnames []string) string {
	var result string

	if len(mdbnames) != 0 {
		results := connectExecute(dir, mdbnames)
		fmt.Println(results + "\n*************************\n")
		result += ".mdb Files Executed"
	} else {
		result += ".mdb Files Empty"
	}

	if len(accdbnames) != 0 {
		fmt.Println(connectExecute(dir, accdbnames) + "\n*************************\n")
		result += " .accdb Files Executed"
	} else {
		result += " .accdb Files Empty"
	}
	return result
}

// walkDir goes through the current dir executes on all db files and then moves on to the next dir recursively
func walkDir(foldernames []string, dir string) {
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
			walkDir(newfoldernames, dir)
		}
	}
}

//printDirInfo prints information about the dir
func printDirInfo(mdbnames []string, accdbnames []string, foldernames []string, dir string) {
	fmt.Printf("Number of .mdb files in %s: %d \nNumber of .accdb files in %s: %d \nNumber of Folders in %s: %d\n\n", dir, len(mdbnames), dir, len(accdbnames), dir, len(foldernames))
}

func main() {
	start := time.Now()
	fmt.Println("\n\n------START OF PROGRAM------")

	dir := "./"
	foldernames := []string{""}
	walkDir(foldernames, dir)

	elapsed := time.Since(start)
	fmt.Printf("\n------COMPLETE------\nFolder(s) Accessed: %d\nFile(s) Accessed: %d\nTable(s) Accessed: %d\nRow(s) Inserted: %d\nTime Taken: %s", numfolders, numfiles, numtables, rowsinserted, elapsed)
}
