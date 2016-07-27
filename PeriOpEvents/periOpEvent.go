package main

import (
	"github.com/raymond3chou/VR/excelHelper"
	"github.com/raymond3chou/VR/periopchecks"
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

//Source the struct to store source information for each event
type Source struct {
	Type string   `json:"type"`
	Path []string `json:"path"`
}

func ()  {

}

func main() {
	tghPath := ""
	tghFile := excelHelper.ConnectToXlsx(tghPath)
	tghRowSlice := excelHelper.ParseData(tghFile)

	twhPath := ""
	twhFile := excelHelper.ConnectToXlsx(twhPath)
	twhRowSlice := excelHelper.ParseData(twhFile)
}
