package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/LynneXie1201/Read_From_Excel/helper"
	"github.com/raymond3chou/VR/accessHelper"
	"github.com/raymond3chou/VR/excelHelper"
	"github.com/raymond3chou/VR/periopchecks"
	"github.com/tealeg/xlsx"
)

//Operation a struct for the operation event
type Operation struct {
	Type       string            `json:"type"`
	MRN        string            `json:"mrn"`
	ResearchID string            `json:"research_id"`
	PeriOpID   int64             `json:"periop_id"`
	PTID       string            `json:"ptid"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	Surgeries  []string          `json:"surgeries"`
	Children   []string          `json:"children"`
	Parent     int64             `json:"parent"`
	Notes      string            `json:"notes"`
	SOURCE     Source            `json:"source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//MI is event for myocardial infarction
type MI struct {
	Type       string            `json:"type"`
	MRN        string            `json:"mrn"`
	ResearchID string            `json:"research_id"`
	PeriOpID   int64             `json:"periop_id"`
	PTID       string            `json:"ptid"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	Mi         int64             `json:"myocardial_infarction"`
	SOURCE     Source            `json:"source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//Pace is the event for a pacemaker
type Pace struct {
	Type       string            `json:"type"`
	MRN        string            `json:"mrn"`
	ResearchID string            `json:"research_id"`
	PeriOpID   int64             `json:"periop_id"`
	PTID       string            `json:"ptid"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	PACE       int64             `json:"pacemaker"`
	SOURCE     Source            `json:"source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//TIA is the event for a Transient ischemic attack
type TIA struct {
	Type       string            `json:"type"`
	MRN        string            `json:"mrn"`
	ResearchID string            `json:"research_id"`
	PeriOpID   int64             `json:"periop_id"`
	PTID       string            `json:"ptid"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	Outcome    int64             `json:"Outcome"`
	Agents     int64             `json:"Agents"`
	When       int64             `json:"When"`
	SOURCE     Source            `json:"Source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//Source the struct to store source information for each event
type Source struct {
	Type string   `json:"type"`
	Path []string `json:"path"`
}

func assignOperation(rowSlice map[string]string, surg []string, source Source, date string) Operation {
	var eventOp Operation

	eventOp.Type = "operation"
	eventOp.PTID = rowSlice["PTID"]
	eventOp.MRN = ""
	eventOp.ResearchID = ""
	eventOp.PeriOpID = excelHelper.StringToInt(rowSlice["ID"])
	eventOp.Date = date
	eventOp.DateEst = 0
	eventOp.Surgeries = surg
	eventOp.SOURCE = source

	return eventOp
}

func assignMI(rowSlice map[string]string, source Source, date string) MI {
	var mi MI
	mi.Type = "myocardial infarction"
	mi.PTID = rowSlice["PTID"]
	mi.PeriOpID = excelHelper.StringToInt(rowSlice["ID"])
	mi.Date = date
	mi.DateEst = 0
	mi.Mi = 1
	mi.SOURCE = source

	return mi
}

func assignPace(rowSlice map[string]string, source Source, date string) Pace {
	var p Pace
	p.Type = "myocardial infarction"
	p.PTID = rowSlice["PTID"]
	p.PeriOpID = excelHelper.StringToInt(rowSlice["ID"])
	p.Date = date
	p.DateEst = 0
	p.PACE = 1
	p.SOURCE = source

	return p
}

func assignTIA(rowSlice map[string]string, source Source, date string) TIA {
	var e TIA

	e.Type = "TIA"
	e.PTID = rowSlice["PTID"]
	e.PeriOpID = excelHelper.StringToInt(rowSlice["ID"])
	e.Date = date
	e.DateEst = 0
	e.Outcome = 3
	e.Agents = 8
	e.When = 1
	e.SOURCE = source

	return e
}

//mapSurgeries creates a map with ptid_date as the key and surgeries as the value
func mapSurgeries(sheet *xlsx.File) map[string]string {
	surgMap := make(map[string]string)
	rowLegth := sheet.Sheets[0].MaxRow

	for r := 1; r < rowLegth; r++ {
		key := sheet.Sheets[0].Rows[r].Cells[2].Value + "_" + sheet.Sheets[0].Rows[r].Cells[3].Value
		surgMap[key] = sheet.Sheets[0].Rows[r].Cells[5].Value
	}
	return surgMap
}

//findSurgeries goes through the surgMap and finds the match for ptid and date returns the surgery
func findSurgeries(ptID string, date string, surgMap map[string]string) string {
	findKey := ptID + "_" + date
	for key := range surgMap {
		if findKey == key {
			return surgMap[key]
		}
	}

	return ""
}

//parseSurgeries takes surgeries string and creates a string slice
func parseSurgeries(s string, redo []string) []string {
	sSlice := strings.Split(s, "|")
	sSlice = append(sSlice, redo...)
	return sSlice
}

func checkRedo(rowSlice map[string]string) []string {
	var redo []string
	if excelHelper.StringToInt(rowSlice["REOP"]) == 6 {
		redo = append(redo, "redoX1")
	}
	return redo
}

//parseData returns row ri from sheet
func parseData(sheet *xlsx.File, ri int) map[string]string {
	rowSlice := make(map[string]string)
	colLength := sheet.Sheets[0].MaxCol

	for ci := 0; ci < colLength; ci++ {
		rowSlice[sheet.Sheets[0].Rows[0].Cells[ci].Value] = sheet.Sheets[0].Rows[ri].Cells[ci].Value
	}
	return rowSlice
}

func objectGenerator(sheet *xlsx.File, surgerieSheet *xlsx.File, tgh bool, jsonFile *os.File, source Source) {
	rowLength := sheet.Sheets[0].MaxRow
	surgMap := mapSurgeries(surgerieSheet)
	for ri := 1; ri < rowLength; ri++ {
		rowSlice := parseData(sheet, ri)

		if tgh {
			date := helper.CheckDateFormat(ri, "DATEOR", rowSlice["DATEOR"])
			ptID := rowSlice["PTID"]
			surgeries := findSurgeries(ptID, date, surgMap)
			redo := checkRedo(rowSlice)
			surg := parseSurgeries(surgeries, redo)
			eventOP := assignOperation(rowSlice, surg, source, date)
			writeJSON(eventOP, jsonFile)

			if excelHelper.StringToInt(rowSlice["MI"]) == 1 {
				eventMI := assignMI(rowSlice, source, date)
				writeJSON(eventMI, jsonFile)
			}

			if excelHelper.StringToInt(rowSlice["PACE"]) == 1 {
				eventPace := assignPace(rowSlice, source, date)
				writeJSON(eventPace, jsonFile)
			}

			if excelHelper.StringToInt(rowSlice["TIA"]) == 1 {
				eventTIA := assignTIA(rowSlice, source, date)
				writeJSON(eventTIA, jsonFile)
			}
		}
	}
}

func sourceGenerator(path string) Source {
	var s Source
	sPath := strings.TrimPrefix(path, "L:\\CVDMC Students\\Raymond Chou\\")
	s.Type = "periop"
	s.Path = append(s.Path, sPath)
	return s
}

//writeJSON writes the struct into JSON format
func writeJSON(newEvent interface{}, jsonFile *os.File) {
	j, err := json.Marshal(newEvent)
	if err != nil {
		log.Println(err)
	}
	jsonFile.Write(j)
}

func main() {
	jsonPath := ""
	accessHelper.CreateFile(jsonPath)
	jsonFile, _ := accessHelper.ConnectToTxt(jsonPath)

	surgeryPath := ""
	sFile := excelHelper.ConnectToXlsx(surgeryPath)

	tghPath := ""
	tghSource := sourceGenerator(tghPath)

	tghFile := excelHelper.ConnectToXlsx(tghPath)
	objectGenerator(tghFile, sFile, true, jsonFile, tghSource)
}
