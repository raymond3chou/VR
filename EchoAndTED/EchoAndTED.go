package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/access"
	"github.com/tealeg/xlsx"
)

//walkDirectory goes through the current folder and looks for excel files and folders
func walkDirectory(folderNames []string, dir string) {
	for _, folderName := range folderNames {
		xlsxNames := []string{}
		newDir := dir + folderName + "\\"
		xlsxNames, newfolderNames := findExcel(newDir)
		if len(xlsxNames) != 0 {
			checkExcel(xlsxNames, newDir)
		}
		walkDirectory(newfolderNames, newDir)
	}
}

//Finds excel files and folders
func findExcel(dir string) ([]string, []string) {
	var xlsxNames []string
	var folderNames []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for _, f := range files {
		if f.IsDir() {
			folderNames = append(folderNames, f.Name())
		} else {
			if strings.Contains(f.Name(), ".xlsx") {
				xlsxNames = append(xlsxNames, f.Name())
			}
		}
	}
	return xlsxNames, folderNames
}

//checkSheet checks if the excel sheet's columns match the checkCriteria
//if so, save a copy of the file to path
func checkSheet(xlsxPath string, checkCriteria string, path string) bool {
	log.Println(xlsxPath)
	sheet, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		log.Fatalln(err)
	}
	if !checkEmptyHeader(sheet) {

		for i := 0; i < sheet.Sheets[0].MaxCol; i++ {
			cellValue := sheet.Sheets[0].Cell(0, i).Value
			cellValue = strings.ToLower(cellValue)
			if strings.Contains(cellValue, checkCriteria) {
				newName := createXlsxName(xlsxPath)
				sheet.Save(path + newName)
				return true
			}
		}
	}
	return false
}

//createXlsxName removes the directory from the path and replaces \ with _
func createXlsxName(xlsxPath string) string {
	xlsxPathNew := strings.TrimPrefix(xlsxPath, "L:\\CVDMC Students\\valve_registry\\")
	newName := strings.Replace(xlsxPathNew, "\\", "_", -1)
	return newName
}

//checkTED checks if the excel file is a ted file by checking if there is ted in the name
func checkTED(xlsx string, checkCriteria string) bool {
	xlsx = strings.ToLower(xlsx)
	if strings.Contains(xlsx, checkCriteria) {
		return true
	}
	return false
}

//checkExcel goes through the excel files and checks if they're echo files or ted files
func checkExcel(xlsxNames []string, newDir string) {
	for _, xlsx := range xlsxNames {
		xlsxPath := newDir + xlsx
		log.Printf("Checking file: %s", xlsx)
		if checkSheet(xlsxPath, "graftsz", "L:\\CVDMC Students\\Raymond Chou\\ted\\") {
			tedFile, _ := access.ConnectToTxt("C:\\Users\\raymond chou\\Desktop\\tedFiles.txt")
			access.FileWrite(tedFile, xlsxPath+"\n")
			tedFile.Close()
		}
		if checkSheet(xlsxPath, "echo", "L:\\CVDMC Students\\Raymond Chou\\echo\\") {
			echoFile, _ := access.ConnectToTxt("C:\\Users\\raymond chou\\Desktop\\echoFiles.txt")
			access.FileWrite(echoFile, xlsxPath+"\n")
			echoFile.Close()
		}
	}
}

//checkEmptyHeader checks if the header is empty
func checkEmptyHeader(sheet *xlsx.File) bool {
	if sheet.Sheets[0].Cell(0, 0).Value == "" {
		return true
	}
	return false
}

func main() {
	dir := "L:\\CVDMC Students\\valve_registry\\"
	folderNames := []string{""}
	errFile := access.CreateErrorLog(true)
	log.SetOutput(errFile)
	defer errFile.Close()
	walkDirectory(folderNames, dir)
}
