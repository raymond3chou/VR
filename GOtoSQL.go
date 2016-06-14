package main

import (
	"database/sql"
	"fmt"
	_ "odbc/driver"
)

func main() {
	conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq=TestDB.accdb")
	if err != nil {
		fmt.Println("Connecting Error")
		return
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Println("Connection Works?")
	}
}
