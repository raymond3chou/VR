package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

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
	PTID       string            `json:"patient_id"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	Surgeon    string            `json:"surgeon"`
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
	PTID       string            `json:"patient_id"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	SOURCE     Source            `json:"source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//Reop return to or
type Reop struct {
	Type       string            `json:"type"`
	MRN        string            `json:"mrn"`
	ResearchID string            `json:"research_id"`
	PeriOpID   int64             `json:"periop_id"`
	PTID       string            `json:"patient_id"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	REOP       int64             `json:"code"`
	SOURCE     Source            `json:"source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//Pace is the event for a pacemaker
type Pace struct {
	Type       string            `json:"type"`
	MRN        string            `json:"mrn"`
	ResearchID string            `json:"research_id"`
	PeriOpID   int64             `json:"periop_id"`
	PTID       string            `json:"patient_id"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	SOURCE     Source            `json:"source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//Tia is the event for a Transient ischemic attack
type Tia struct {
	Type       string            `json:"type"`
	MRN        string            `json:"mrn"`
	ResearchID string            `json:"research_id"`
	PeriOpID   int64             `json:"periop_id"`
	PTID       string            `json:"patient_id"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	Outcome    int64             `json:"outcome"`
	Agents     int64             `json:"anti_agents"`
	SOURCE     Source            `json:"source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//Stroke is the event for a stroke
type Stroke struct {
	Type       string            `json:"type"`
	MRN        string            `json:"mrn"`
	ResearchID string            `json:"research_id"`
	PeriOpID   int64             `json:"periop_id"`
	PTID       string            `json:"patient_id"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	Outcome    int64             `json:"outcome"`
	Agents     int64             `json:"anti_agents"`
	When       int64             `json:"when"`
	SOURCE     Source            `json:"source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//Survival is the event for when survial = 0
type Survival struct {
	Type       string            `json:"type"`
	MRN        string            `json:"mrn"`
	ResearchID string            `json:"research_id"`
	PeriOpID   int64             `json:"periop_id"`
	PTID       string            `json:"patient_id"`
	Date       string            `json:"date"`
	DateEst    int64             `json:"date_est"`
	Reason     string            `json:"reason"`
	PrmDeath   int64             `json:"primary_cause"`
	Operative  int64             `json:"operative"`
	SOURCE     Source            `json:"source"`
	FIX        []periopcheck.Fix `json:"fix"`
}

//Source the struct to store source information for each event
type Source struct {
	Type string   `json:"type"`
	Path []string `json:"path"`
}

func assignOperation(rowSlice map[string]string, surg []string, source Source, date string) Operation {
	var eventOp Operation
	var f []periopcheck.Fix
	eventOp.Type = "operation"
	eventOp.PTID = rowSlice["PTID"]
	eventOp.MRN = ""
	eventOp.ResearchID = ""
	eventOp.PeriOpID = excelHelper.StringToInt(rowSlice["ID"], 0, "ID")
	eventOp.Date = date
	eventOp.DateEst = 0
	eventOp.Surgeon = rowSlice["SURG"]
	if surg[0] != "" {
		eventOp.Surgeries = surg
	}

	eventOp.SOURCE = source
	eventOp.FIX = f
	return eventOp
}

func assignMI(rowSlice map[string]string, source Source, date string) MI {
	var mi MI
	var f []periopcheck.Fix

	mi.Type = "myocardial_infarction"
	mi.PTID = rowSlice["PTID"]
	mi.PeriOpID = excelHelper.StringToInt(rowSlice["ID"], 0, "ID")
	mi.Date = date
	mi.DateEst = 0
	mi.SOURCE = source
	mi.FIX = f
	return mi
}

func assignPace(rowSlice map[string]string, source Source, date string) Pace {
	var p Pace
	var f []periopcheck.Fix

	p.Type = "perm_pacemaker"
	p.PTID = rowSlice["PTID"]
	p.PeriOpID = excelHelper.StringToInt(rowSlice["ID"], 0, "ID")
	p.Date = date
	p.DateEst = 0
	p.SOURCE = source
	p.FIX = f
	return p
}

func assignTIA(rowSlice map[string]string, source Source, date string) Tia {
	var e Tia
	var f []periopcheck.Fix

	e.Type = "tia"
	e.PTID = rowSlice["PTID"]
	e.PeriOpID = excelHelper.StringToInt(rowSlice["ID"], 0, "ID")
	e.Date = date
	e.DateEst = 0
	e.SOURCE = source
	e.Outcome = 3
	e.Agents = 2
	e.FIX = f
	return e
}

func assignREOP(rowSlice map[string]string, source Source, date string, reop string) Reop {
	var e Reop
	var f []periopcheck.Fix

	e.Type = "return_to_or"
	e.PTID = rowSlice["PTID"]
	e.PeriOpID = excelHelper.StringToInt(rowSlice["ID"], 0, "ID")
	e.Date = date
	e.DateEst = 0
	e.SOURCE = source
	e.REOP = excelHelper.StringToInt(rowSlice[reop], 0, "REOP")
	e.FIX = f
	return e
}

func assignStroke(rowSlice map[string]string, source Source, date string) Stroke {
	var p Stroke
	var f []periopcheck.Fix

	p.Type = "stroke"
	p.PTID = rowSlice["PTID"]
	p.PeriOpID = excelHelper.StringToInt(rowSlice["ID"], 0, "ID")
	p.Date = date
	p.DateEst = 0
	p.SOURCE = source
	p.Outcome = 3
	p.Agents = 8
	p.When = 1
	p.FIX = f
	return p
}

func assignSurvival(rowSlice map[string]string, source Source, date string) Survival {
	var p Survival
	var f []periopcheck.Fix

	p.Type = "death"
	p.PTID = rowSlice["PTID"]
	p.PeriOpID = excelHelper.StringToInt(rowSlice["ID"], 0, "ID")
	p.Date = date
	p.DateEst = 0
	p.Reason = rowSlice["NOTES"]
	p.PrmDeath = -9
	p.Operative = 1
	p.SOURCE = source
	p.FIX = f
	return p
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
	sSlice = excelHelper.SliceTrimSpace(sSlice)
	sSlice = append(sSlice, redo...)
	return sSlice
}

func checkRedo(rowSlice map[string]string, row int) []string {
	var redo []string
	count := 0
	for k := range rowSlice {
		if strings.Contains(k, "REOP") && !strings.Contains(k, "PUMP") {
			if excelHelper.StringToInt(rowSlice[k], row, k) == 6 {
				count++
			}
		}
	}
	if count != 0 {
		for i := 1; i <= count; i++ {
			t := strconv.Itoa(i)

			redoInsert := "redoX" + t
			redo = append(redo, redoInsert)
		}
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
			dStr := rowSlice["DATEOR"]
			ptID := rowSlice["PTID"]
			surgeries := findSurgeries(ptID, dStr, surgMap)
			redo := checkRedo(rowSlice, ri)
			surg := parseSurgeries(surgeries, redo)
			d := excelHelper.StringToFloat(dStr)
			date := excelHelper.DateConvertor(d)
			eventOP := assignOperation(rowSlice, surg, source, date)
			writeJSON(eventOP, jsonFile)

			if excelHelper.StringToInt(rowSlice["MI"], ri, "MI") == 1 {
				eventMI := assignMI(rowSlice, source, date)
				writeJSON(eventMI, jsonFile)
			}

			if excelHelper.StringToInt(rowSlice["PACE"], ri, "PACE") == 1 {
				eventPace := assignPace(rowSlice, source, date)
				writeJSON(eventPace, jsonFile)
			}
			for k := range rowSlice {
				if strings.Contains(k, "REOP") && !strings.Contains(k, "PUMP") && !strings.Contains(k, "NUM") {
					reop := excelHelper.StringToInt(rowSlice[k], ri, k)
					if reop != 6 && reop != -9 && reop != 0 && reop != 9 {
						eventREOP := assignREOP(rowSlice, source, date, k)
						writeJSON(eventREOP, jsonFile)
					}
				}
			}

			if excelHelper.StringToInt(rowSlice["TIA"], ri, "TIA") == 1 {
				eventTIA := assignTIA(rowSlice, source, date)
				writeJSON(eventTIA, jsonFile)
			}

			if excelHelper.StringToInt(rowSlice["STROKE"], ri, "STROKE") == 1 {
				eventStroke := assignStroke(rowSlice, source, date)
				writeJSON(eventStroke, jsonFile)
			}

			if excelHelper.StringToInt(rowSlice["SURVIVAL"], ri, "SURVIVAL") == 0 {
				eventSurvival := assignSurvival(rowSlice, source, date)
				writeJSON(eventSurvival, jsonFile)
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
	jsonFile.WriteString("\n")
}

func main() {
	jsonPath := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\periOpEvents.json"
	accessHelper.CreateFile(jsonPath)
	jsonFile, _ := accessHelper.ConnectToTxt(jsonPath)
	defer jsonFile.Close()

	surgeryPath := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\surgeries.xlsx"
	sFile := excelHelper.ConnectToXlsx(surgeryPath)

	tghPath := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TGH perioperative.xlsx"
	tghSource := sourceGenerator(tghPath)

	tghFile := excelHelper.ConnectToXlsx(tghPath)
	objectGenerator(tghFile, sFile, true, jsonFile, tghSource)

	surgeryPath = "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TWH surgeries.xlsx"
	sFile = excelHelper.ConnectToXlsx(surgeryPath)

	twhPath := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TWH perioperative.xlsx"
	twhSource := sourceGenerator(twhPath)

	twhFile := excelHelper.ConnectToXlsx(twhPath)
	objectGenerator(twhFile, sFile, true, jsonFile, twhSource)
}
