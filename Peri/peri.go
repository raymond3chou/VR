package main

import (
	"log"
	"strings"

	"github.com/tealeg/xlsx"
)

//connectToXlsx connects to the excel file through the path
func connectToXlsx(xlsxPath string) *xlsx.File {
	sheet, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		log.Fatalln(err)
	}
	return sheet
}

//identifyCols returns all the columns in the sheet in a slice of string
func identifyCols(sheet *xlsx.File) []string {
	var colNamesSlice []string
	for i := 0; i < sheet.Sheets[0].MaxCol; i++ {
		cellValue := sheet.Sheets[0].Cell(0, i).Value
		cellValue = strings.ToLower(cellValue)
		colNamesSlice = append(colNamesSlice, cellValue)
	}
	return colNamesSlice
}

func colCompare() {

}
