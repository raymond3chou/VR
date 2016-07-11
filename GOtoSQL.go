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

//Global Variables to Track files accessed
var (
	numtables    int
	numfiles     int
	numfolders   int
	rowsinserted int
)

//checkFollowup queries the specified table to check if it is a Followup Table(i.e contains FU_D,DIED,DTH_D)
//based on the columns present, it also returns a list of column names matched to a masterlist
//If no column names are present then it returns an empty query
func checkFollowup(conn *sql.DB, tablename string) (bool, []string, string) {
	followup := false
	maincolumns := []string{"PTID", "CHART", "LNAME", "FNAME", "SEX", "AGE", "STREET", "CITY", "PROVINCE", "POSTCODE", "PHONEHOME", "PHONEWORK", "PHONECELL", "EMAIL", "DOB", "CARDIO1", "CARDIO2", "GP1", "GP2"}
	var query string
	//Attempts to Run the Query
	rows, err := conn.Query("SELECT * FROM [" + tablename + "]")
	if err != nil {
		issue := "Cant Run SELECT * FROM [" + tablename + "]"
		fmt.Println(err)
		log.Panic(issue)
	} else {
		//Returns the Columns from the QUERY
		columnArray, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
			log.Panic(err.Error())
		} else {

			for _, columnname := range columnArray {
				for _, maincolname := range maincolumns {
					if columnname == "FU_D" {
						followup = true
					}
					if columnname == "DIED" {
						followup = true
					}
					if columnname == "DTH_D" {
						followup = true
					}
					if strings.Contains(columnname, maincolname) {
						query += " [" + columnname + "],"
					}
				}
			}
			query = strings.TrimSuffix(query, ",")
		}
	}
	return followup, maincolumns, query
}

//selectAccess
func selectAccess(conn *sql.DB, file *os.File, tablename string) (int, int) {
	//Queries the database, and passes the row to be appended to the test file
	//returns # inserted
	var selectquery string
	var query string

	followup, maincolumns, query := checkFollowup(conn, tablename)
	if followup && query != "" {
		selectquery = "SELECT" + query + " FROM [" + tablename + "]"
	} else if query == "" {
		issue := tablename + " does not contain any columns related to PHI\n"
		log.Print(issue)
		return 0, 0
	} else {
		issue := tablename + " is not a Follow Up Table\n"
		log.Print(issue)
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
		log.Panic(issue)
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
			log.Panic(issue)
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
		log.Panic(err.Error())
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
		log.Println(err)
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
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			var table string
			err = rows.Scan(&table)
			if err != nil {
				log.Println("Failed to pass tablenames")
			}
			tablenames = append(tablenames, table)
		}
	}
	return tablenames

}

//connectToDB Connects to a specified database in a specified directory
func connectToDB(dir string, dbname string) *sql.DB {

	dbq := dir + "/" + dbname
	fmt.Println("Connecting to " + dbq)
	conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq="+dbq)
	if err != nil {
		log.Println("Connection to " + dbq + " Failed")
	}
	fmt.Println("Connected to " + dbq)
	return conn
}

//connectExecute connects to every database in a specified directory and checks if its is a FU table
//if so, it reads the columns; sorts them; and writes them on to a textfile
func connectExecute(dir string, dbnames []string) string {
	var dbaccessed int
	for _, dbname := range dbnames {
		isConnected := connectToDB(dir, dbname)
		if isConnected != nil {
			dbaccessed++
		} else {
			continue
		}
		//Originating Database isConnectedection established

		tablenames := findTable(isConnected)
		tablelength := len(tablenames)
		fmt.Printf("%d Tables in %s\n", tablelength, dbname)

		file, connection := access.ConnectToTxt("C:\\Users\\raymond chou\\Desktop\\ContactInfo.txt")
		if !connection {
			continue
		}
		//Opens text file that can be constantly appended to. ONLY NEEDS TO BE CALLED ONCE
		var tableused int
		for _, tablename := range tablenames {
			inserted, numberofRows := selectAccess(isConnected, file, tablename)
			fmt.Printf("Total Number of Rows Read and Inserted from %s = %d/%d\n", tablename, inserted, numberofRows)
			tableused++
			rowsinserted += inserted
		}
		isConnected.Close()
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
	for _, foldername := range foldernames {
		newDir := dir + foldername + "/"
		fmt.Printf("\nEntering Folder Directory: %s\n", newDir)
		numfolders++
		mdbnames, accdbnames, newfoldernames := findDB(newDir)
		printDirInfo(mdbnames, accdbnames, foldernames, dir)
		result := dbPresent(newDir, mdbnames, accdbnames)
		fmt.Printf("\n%s \n", result)
		fmt.Println(newDir)
		walkDir(newfoldernames, newDir)
	}

}

//printDirInfo prints information about the dir
func printDirInfo(mdbnames []string, accdbnames []string, foldernames []string, dir string) {
	fmt.Printf("Number of .mdb files in %s: %d \nNumber of .accdb files in %s: %d \nNumber of Folders in %s: %d\n\n", dir, len(mdbnames), dir, len(accdbnames), dir, len(foldernames))
}

func main() {
	errFile := access.CreateErrorLog(true)
	log.SetOutput(errFile)
	defer errFile.Close()
	dir := "./"
	start := time.Now()
	fmt.Println("\n\n------START OF PROGRAM------")

	foldernames := []string{""}
	walkDir(foldernames, dir)

	elapsed := time.Since(start)
	fmt.Printf("\n------COMPLETE------\nFolder(s) Accessed: %d\nFile(s) Accessed: %d\nTable(s) Accessed: %d\nRow(s) Inserted: %d\nTime Taken: %s", numfolders, numfiles, numtables, rowsinserted, elapsed)
}
