package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/tealeg/xlsx"
)

//EchoPaths keeps track of all the Echo file paths
var EchoPaths []string

//Go Through All Databases
func walkDirectory(folderNames []string, dir string) {
	for _, folderName := range folderNames {
		newDir := dir + folderName + "\\"
		xlsxNames, newfolderNames := findExcel(newDir)
		for _, xlsx := range xlsxNames {
			xlsxPath := newDir + xlsx
			if checkExcel(xlsxPath, "Echo") {
				EchoPaths = append(EchoPaths, newDir+xlsx)
			}
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

func checkExcel(xlsxPath string, checkCriteria string) bool {
	sheet, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range sheet.Sheets[0].Rows[0].Cells {
		fmt.Println(v.Value)
		if v.Value == checkCriteria {
			return true
		}
	}
	return false
}
