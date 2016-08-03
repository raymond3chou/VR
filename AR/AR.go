package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/alexbrainman/odbc"
	"github.com/raymond3chou/VR/accessHelper"
)

func iterateDB(listFiles []string, file *os.File, dir string) {
	for _, l := range listFiles {
		newPath := dir + "\\" + l
		sqlDB := connectToDB(newPath)
		defer sqlDB.Close()
		tableList := findTable(sqlDB)
		for _, t := range tableList {
			if strings.Contains(strings.ToLower(t), "echo") {
				if !strings.Contains(t, "COUNTRY") {
					selectAR(sqlDB, file, t)
				}
			}
		}
	}
}

//findTable iterates through all the tables in the database.
//Currently only works with .mdb. .accdb does not have permission.
func findTable(conn *sql.DB) []string {
	var tablenames []string

	rows, err := conn.Query("SELECT Name FROM MSysObjects WHERE Type=1 AND Flags=0;")
	if err != nil {
		log.Println(err)
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
func connectToDB(path string) *sql.DB {
	log.Println("Connecting to " + path)
	fmt.Println("Connecting to " + path)

	conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq="+path)
	if err != nil {
		fmt.Println("Connection to " + path + " Failed")
		log.Println("Connection to " + path + " Failed")
	} else {
		log.Println("Connected to " + path)
		fmt.Println("Connected to " + path)

	}
	return conn
}

func selectAR(conn *sql.DB, file *os.File, tableName string) {
	query := "SELECT PTID, ECHODATE,AR FROM [" + tableName + "]"
	rows, err := conn.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ptID sql.NullString
		var echoDate sql.NullString
		var ar sql.NullString

		err = rows.Scan(&ptID, &echoDate, &ar)
		if err != nil {
			log.Println(err)
		}
		date := strings.Split(echoDate.String, "T")
		row := ptID.String + "|" + date[0] + "|" + ar.String + "\n"
		accessHelper.FileWrite(file, row)
	}
}

func readDir(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}

	var listDB []string

	for _, f := range files {
		if strings.Contains(f.Name(), ".mdb") || strings.Contains(f.Name(), ".accdb") || strings.Contains(f.Name(), ".MDB") {
			listDB = append(listDB, f.Name())
		}
	}
	return listDB
}

func main() {
	errFile := accessHelper.CreateErrorLog(true)
	log.SetOutput(errFile)
	defer errFile.Close()

	dir := "F:\\DavidOp"
	listFiles := readDir(dir)
	echotxt := dir + "\\echo.txt"
	accessHelper.CreateFile(echotxt)
	wFile, _ := accessHelper.ConnectToTxt(echotxt)

	iterateDB(listFiles, wFile, dir)
	defer wFile.Close()
}
