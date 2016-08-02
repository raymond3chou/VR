package excelHelper

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/raymond3chou/VR/accessHelper"
	"github.com/tealeg/xlsx"
)

//ConnectToXlsx connects to the excel file through the path
func ConnectToXlsx(xlsxPath string) *xlsx.File {
	sheet, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		log.Fatalln(err)
	}
	return sheet
}

//IdentifyCols returns all the columns in the sheet in a slice of string
func IdentifyCols(sheet *xlsx.File) []string {
	var colNamesSlice []string
	for i := 0; i < sheet.Sheets[0].MaxCol; i++ {
		cellValue := sheet.Sheets[0].Cell(0, i).Value
		cellValue = strings.ToLower(cellValue)
		colNamesSlice = append(colNamesSlice, cellValue)
	}
	return colNamesSlice
}

//ColCompare identifies the common strings in the two string slices
func ColCompare(tghCols []string, twhCols []string) []string {
	var commonColsSlice []string
	for _, twh := range twhCols {
		for _, tgh := range tghCols {
			if tgh == twh {
				commonColsSlice = append(commonColsSlice, tgh)
			}
		}
	}
	return commonColsSlice
}

//NotPresentinSlice identifies the strings that are not common between the two string slices
func NotPresentinSlice(originalCols []string, commonColsSlice []string) []string {
	var unCommonColsSlice []string
	var present bool
	for _, original := range originalCols {
		present = false
		for _, commonCol := range commonColsSlice {
			if commonCol == original {
				present = true
				break
			}
		}
		if !present {
			unCommonColsSlice = append(unCommonColsSlice, original)
		}
	}
	return unCommonColsSlice
}

//PrintSlice prints the slice
func PrintSlice(slice []string) {
	for _, s := range slice {
		fmt.Printf(" %s |", s)
	}
}

//ReadRow reads a row in excel and returns it as a slice of string
func ReadRow(sheet *xlsx.File) []string {
	return nil
}

//WriteStruct writes the struct type to text so it can be copied into peri.go
func WriteStruct(colNameSlice []string) {
	path := "C:\\Users\\raymond chou\\Desktop\\PeriOp\\struct.txt"
	accessHelper.CreateFile(path)
	file, _ := accessHelper.ConnectToTxt(path)
	for _, c := range colNameSlice {
		var structPrint string
		lowerC := strings.ToLower(c)
		upperC := strings.ToUpper(c)
		structPrint += upperC
		structPrint += " " + "int"
		structPrint += "`json:\"" + lowerC + "\"`\n"
		accessHelper.FileWrite(file, structPrint)
	}
}

//PeriOpLiteral prints to text the literal for the stuct
func PeriOpLiteral(colNameSlice []string) {
	path := "C:\\Users\\raymond chou\\Desktop\\periopliteral.txt"
	accessHelper.CreateFile(path)
	file, _ := accessHelper.ConnectToTxt(path)
	for _, c := range colNameSlice {
		var structPrint string
		lowerC := strings.ToLower(c)
		if strings.Contains(lowerC, "reop") {
			continue
		}
		if strings.Contains(lowerC, "mi") {
			continue
		}
		if strings.Contains(lowerC, "pace") {
			continue
		}
		if strings.Contains(lowerC, "tia") {
			continue
		}
		if strings.Contains(lowerC, "stroke") {
			continue
		}
		if strings.Contains(lowerC, "survival") {
			continue
		}
		upperC := strings.ToUpper(c)
		structPrint += upperC
		structPrint += " " + "string"
		structPrint += "`json:\"" + lowerC + "\"`\n"
		accessHelper.FileWrite(file, structPrint)
	}
}

//StringToInt converts string to int
func StringToInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return i
}

//SliceTrimSpace trims trailing and leading spaces from each string in the slice
func SliceTrimSpace(s []string) []string {
	for i, ss := range s {
		s[i] = strings.TrimSpace(ss)
	}
	return s
}

//StringToFloat converts string to float64
func StringToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println(err)
	}
	return f
}

//DateConvertor converts a float64 to date
func DateConvertor(f float64) string {
	date := xlsx.TimeFromExcelTime(f, false)
	dateSlice := strings.Split(date.String(), " ")
	d := dateSlice[0]
	return d
}
