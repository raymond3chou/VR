package main

import (
	"testing"

	"github.com/raymond3chou/VR/accessHelper"
	"github.com/raymond3chou/VR/excelHelper"
)

func TestSourceGenerator(t *testing.T) {
	path := "L:\\CVDMC Students\\Raymond Chou\\perioperative"
	s := sourceGenerator(path)
	if s.Type != "periop" {
		t.Errorf("s.Type not equal to periop, is %s instead\n", s.Type)
	}
	if s.Path[0] != "perioperative" {
		t.Errorf("s.Path[0] not equal to periop, is %s instead\n", s.Path[0])
	}
}

func TestParseData(t *testing.T) {
	path := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TGH perioperative.xlsx"
	file := excelHelper.ConnectToXlsx(path)
	rowSlice := parseData(file, 1)
	if rowSlice["PTID"] != "AGUP062964" {
		t.Errorf("PTID is %s not AGUP062964", rowSlice["PTID"])
	}
	if rowSlice["AVSURG"] != "0" {
		t.Errorf("AVSURG is %s not 0\n", rowSlice["AVSURG"])
	}

	rowSlice = parseData(file, 14)
	if rowSlice["STRESS"] != "" {
		t.Errorf("STRESS is %s not ``\n", rowSlice["STRESS"])
	}

	path = "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TWH perioperative.xlsx"
	file = excelHelper.ConnectToXlsx(path)
	rowSlice = parseData(file, 1)

	if rowSlice["PTID"] != "DONJ062136" {
		t.Errorf("PTID is %s not DONJ062136", rowSlice["PTID"])
	}

}

func TestCheckRedo(t *testing.T) {
	path := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TGH perioperative.xlsx"
	file := excelHelper.ConnectToXlsx(path)
	rowSlice := parseData(file, 1)
	redo := checkRedo(rowSlice, 0)
	if len(redo) != 0 {
		t.Error("redo not empty")
	}
	rowSlice = parseData(file, 1087)
	redo2 := checkRedo(rowSlice, 0)
	if len(redo2) != 1 {
		t.Error("redo is empty")
	}
}

func TestParseSurgeries(t *testing.T) {
	path := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\surgeries.xlsx"
	file := excelHelper.ConnectToXlsx(path)
	surgMap := mapSurgeries(file)
	surgString := findSurgeries("ANDR040831", "37368", surgMap)
	if surgString != "spAAS | TVrepair" {
		t.Errorf("the surg string is %s not spAAS | TVrepair\n", surgString)
	}
	redo := []string{"redoX1"}
	surgSlice := parseSurgeries(surgString, redo)
	expectedSurgSlice := []string{"spAAS", "TVrepair", "redoX1"}
	if !accessHelper.CompareSlice(surgSlice, expectedSurgSlice) {
		t.Errorf("The two slices do not match\n")
	}
}

// func TestRedo(t *testing.T) {
// 	rowSlice := make(map[string]string)
// 	rowSlice["REOP"] = "6"
// 	rowSlice["Jim"] = "32"
// 	rowSlice["REOP5"] = "6"
// 	rowSlice["REOP2"] = "2"
// 	redo := checkRedo(rowSlice)
// }

func TestObjectGenerator(t *testing.T) {
	tghPath := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\tgh.xlsx"
	tghFile := excelHelper.ConnectToXlsx(tghPath)
	surgPath := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\surgeries.xlsx"
	surgFile := excelHelper.ConnectToXlsx(surgPath)
	jsonPath := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\src\\github.com\\raymond3chou\\VR\\PeriOpEvents\\MockData\\mock.json"
	jsonFile, _ := accessHelper.ConnectToTxt(jsonPath)
	tghSource := sourceGenerator(tghPath)
	objectGenerator(tghFile, surgFile, true, jsonFile, tghSource)
}
