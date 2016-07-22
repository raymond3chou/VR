package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/LynneXie1201/Read_From_Excel/helper"
	"github.com/access/excelHelper"
	"github.com/access/periopchecks"
	"github.com/tealeg/xlsx"
)

//Event type for event object to be printed to JSON
type Event struct {
	PTID       string   `json:"ptid"`
	Type       string   `json:"type"`
	Date       string   `json:"date"`
	ResearchID string   `json:"research_id"`
	MRN        string   `json:"mrn"`
	Surgeries  string   `json:"surgeries"`
	Surgeon    string   `json:"surgeon"`
	PeriOpID   int      `json:"periop_id"`
	Fix        []string `json:"fix"`
}

//MultiField type for multiple key-value pairs
type MultiField struct {
	Name string
	Min  int
	Max  int
}

//Field type for key-value pair
type Field struct {
	Name string
	Min  int
	Max  int
}

//PeriOp is the object for periop data
type PeriOp struct {
	PERIOPID   string            `json:"periop_id"`
	PTID       string            `json:"ptid"`
	AREA       string            `json:"area"`
	TRIAGE     string            `json:"triage"`
	SDA        string            `json:"sda"`
	ADDATE     string            `json:"addate"`
	DATEOR     string            `json:"dateor"`
	DISDATE    string            `json:"disdate"`
	DAYSPOST   string            `json:"dayspost"`
	ICUNUM     string            `json:"icunum"`
	ICU        string            `json:"icu"`
	ICUTIME    string            `json:"icutime"`
	VENT       string            `json:"vent"`
	VENTTIME   string            `json:"venttime"`
	SURG       string            `json:"surg"`
	ASSIST     string            `json:"assist"`
	FDOC       string            `json:"fdoc"`
	CDOC       string            `json:"cdoc"`
	CATRIAL    string            `json:"catrial"`
	ANTRIAL    string            `json:"antrial"`
	TIMING     string            `json:"timing"`
	FROM       string            `json:"from"`
	ACBREDO    string            `json:"acbredo"`
	AVREDO     string            `json:"avredo"`
	MVREDO     string            `json:"mvredo"`
	TVREDO     string            `json:"tvredo"`
	OTHREDO    string            `json:"othredo"`
	DATEORP1   string            `json:"dateorp1"`
	ACBREDOP1  string            `json:"acbredop1"`
	AVREDOP1   string            `json:"avredop1"`
	MVREDOP1   string            `json:"mvredop1"`
	TVREDOP1   string            `json:"tvredop1"`
	OTHREDOP1  string            `json:"othredop1"`
	DATEORP2   string            `json:"dateorp2"`
	ACBREDOP2  string            `json:"acbredop2"`
	AVREDOP2   string            `json:"avredop2"`
	MVREDOP2   string            `json:"mvredop2"`
	TVREDOP2   string            `json:"tvredop2"`
	OTHREDOP2  string            `json:"othredop2"`
	DATEORP3   string            `json:"dateorp3"`
	ACBREDOP3  string            `json:"acbredop3"`
	AVREDOP3   string            `json:"avredop3"`
	MVREDOP3   string            `json:"mvredop3"`
	TVREDOP3   string            `json:"tvredop3"`
	OTHREDOP3  string            `json:"othredop3"`
	DATEORP4   string            `json:"dateorp4"`
	ACBREDOP4  string            `json:"acbredop4"`
	AVREDOP4   string            `json:"avredop4"`
	MVREDOP4   string            `json:"mvredop4"`
	TVREDOP4   string            `json:"tvredop4"`
	OTHREDOP4  string            `json:"othredop4"`
	DATEORP5   string            `json:"dateorp5"`
	ACBREDOP5  string            `json:"acbredop5"`
	AVREDOP5   string            `json:"avredop5"`
	MVREDOP5   string            `json:"mvredop5"`
	TVREDOP5   string            `json:"tvredop5"`
	OTHREDOP5  string            `json:"othredop5"`
	DATEORP6   string            `json:"dateorp6"`
	ACBREDOP6  string            `json:"acbredop6"`
	AVREDOP6   string            `json:"avredop6"`
	MVREDOP6   string            `json:"mvredop6"`
	TVREDOP6   string            `json:"tvredop6"`
	OTHREDOP6  string            `json:"othredop6"`
	PRECARD    string            `json:"precard"`
	PIDATE     string            `json:"pidate"`
	PITHROMB   string            `json:"pithromb"`
	PITDATE    string            `json:"pitdate"`
	CATH       string            `json:"cath"`
	CATHDATE   string            `json:"cathdate"`
	ANGINA     string            `json:"angina"`
	PREOPMI    string            `json:"preopmi"`
	MIDATE     string            `json:"midate"`
	NYHA       string            `json:"nyha"`
	CCS        string            `json:"ccs"`
	LVGRADE    string            `json:"lvgrade"`
	STRESS     string            `json:"stress"`
	DIABETES   string            `json:"diabetes"`
	DOI        string            `json:"doi"`
	HYPER      string            `json:"hyper"`
	CHLSTRL    string            `json:"chlstrl"`
	FHX        string            `json:"fhx"`
	SMOKE      string            `json:"smoke"`
	PACKS      string            `json:"packs"`
	COPD       string            `json:"copd"`
	COPDS      string            `json:"copds"`
	THROMB     string            `json:"thromb"`
	PVD        string            `json:"pvd"`
	RF         string            `json:"rf"`
	NEWRF      string            `json:"newrf"`
	DIAL       string            `json:"dial"`
	MARFAN     string            `json:"marfan"`
	CAROTID    string            `json:"carotid"`
	AAD        string            `json:"aad"`
	EU         string            `json:"eu"`
	RECG       string            `json:"recg"`
	CHF        string            `json:"chf"`
	SHOCK      string            `json:"shock"`
	SYNCOPE    string            `json:"syncope"`
	ASP        string            `json:"asp"`
	AMI        string            `json:"ami"`
	CREAT      string            `json:"creat"`
	STATIN     string            `json:"statin"`
	AVDIS      string            `json:"avdis"`
	MVDIS      string            `json:"mvdis"`
	TVDIS      string            `json:"tvdis"`
	ENDOCARD   string            `json:"endocard"`
	URGENT     string            `json:"urgent"`
	AVSURG     string            `json:"avsurg"`
	D2         string            `json:"d2"`
	AVSIZE     string            `json:"avsize"`
	AVPATH     string            `json:"avpath"`
	AVPATH2    string            `json:"avpath2"`
	AVPATH3    string            `json:"avpath3"`
	ANNULEN    string            `json:"annulen"`
	AVPROS     string            `json:"avpros"`
	MVSURG     string            `json:"mvsurg"`
	MVSIZE     string            `json:"mvsize"`
	MVPATH     string            `json:"mvpath"`
	MVPATH2    string            `json:"mvpath2"`
	MVPATH3    string            `json:"mvpath3"`
	MVANN      string            `json:"mvann"`
	CHORD      string            `json:"chord"`
	GORTEX     string            `json:"gortex"`
	MVP        string            `json:"mvp"`
	MVC        string            `json:"mvc"`
	CHORDAL    string            `json:"chordal"`
	MVPROS     string            `json:"mvpros"`
	TVSURG     string            `json:"tvsurg"`
	TVSIZE     string            `json:"tvsize"`
	TVPATH     string            `json:"tvpath"`
	TVPATH2    string            `json:"tvpath2"`
	TVPATH3    string            `json:"tvpath3"`
	TVPROS     string            `json:"tvpros"`
	PVSURG     string            `json:"pvsurg"`
	PVSIZE     string            `json:"pvsize"`
	PVPROS     string            `json:"pvpros"`
	CI         string            `json:"ci"`
	MPAP       string            `json:"mpap"`
	SYSAVG     string            `json:"sysavg"`
	LVEDP      string            `json:"lvedp"`
	PVR        string            `json:"pvr"`
	MVGRADR    string            `json:"mvgradr"`
	AVAREA     string            `json:"avarea"`
	MVAREA     string            `json:"mvarea"`
	AVXPLTYPE  string            `json:"avxpltype"`
	AVXPLSIZE  string            `json:"avxplsize"`
	AVXPLPATH  string            `json:"avxplpath"`
	AVXPLDATE  string            `json:"avxpldate"`
	MVXPLTYPE  string            `json:"mvxpltype"`
	MVXPLSIZE  string            `json:"mvxplsize"`
	MVXPLPATH  string            `json:"mvxplpath"`
	MVXPLDATE  string            `json:"mvxpldate"`
	TVXPLTYPE  string            `json:"tvxpltype"`
	TVXPLSIZE  string            `json:"tvxplsize"`
	TVXPLPATH  string            `json:"tvxplpath"`
	TVXPLDATE  string            `json:"tvxpldate"`
	ASSOCOP    string            `json:"assocop"`
	LVA        string            `json:"lva"`
	SEPT       string            `json:"sept"`
	SEPTYPE    string            `json:"septype"`
	CHD        string            `json:"chd"`
	AAS        string            `json:"aas"`
	AOPATH     string            `json:"aopath"`
	MAZE       string            `json:"maze"`
	MISC       string            `json:"misc"`
	OTHERTYPE  string            `json:"othertype"`
	ASSDEV     string            `json:"assdev"`
	DEVICETYPE string            `json:"devicetype"`
	DISLAD     string            `json:"dislad"`
	DISCX      string            `json:"discx"`
	DISRCA     string            `json:"disrca"`
	LMAIN      string            `json:"lmain"`
	DISNUM     string            `json:"disnum"`
	LIMA       string            `json:"lima"`
	RIMA       string            `json:"rima"`
	RADIAL     string            `json:"radial"`
	SKEGRAFT   string            `json:"skegraft"`
	GFTLAD     string            `json:"gftlad"`
	GFTCX      string            `json:"gftcx"`
	GFTRCA     string            `json:"gftrca"`
	ENDART     string            `json:"endart"`
	CVPA       string            `json:"cvpa"`
	OTHGFT     string            `json:"othgft"`
	ACBNUM     string            `json:"acbnum"`
	PUMPCASE   string            `json:"pumpcase"`
	MININV     string            `json:"mininv"`
	ORTIME     string            `json:"ortime"`
	PUMP       string            `json:"pump"`
	CLAMP      string            `json:"clamp"`
	CIRARR     string            `json:"cirarr"`
	BSA        string            `json:"bsa"`
	HT         string            `json:"ht"`
	WT         string            `json:"wt"`
	MYOPRO     string            `json:"myopro"`
	TECH       string            `json:"tech"`
	DIRECT     string            `json:"direct"`
	HYPOTHER   string            `json:"hypother"`
	OFFPUMP    string            `json:"offpump"`
	IABP       string            `json:"iabp"`
	REOPNUM    string            `json:"reopnum"`
	REOP       string            `json:"reop"`
	REOP2      string            `json:"reop2"`
	REOP3      string            `json:"reop3"`
	REOP4      string            `json:"reop4"`
	REOP5      string            `json:"reop5"`
	REOPPUMP   string            `json:"reoppump"`
	REOPPUMP2  string            `json:"reoppump2"`
	REOPPUMP3  string            `json:"reoppump3"`
	REOPPUMP4  string            `json:"reoppump4"`
	REOPPUMP5  string            `json:"reoppump5"`
	IECG       string            `json:"iecg"`
	CK         string            `json:"ck"`
	CKMB       string            `json:"ckmb"`
	MI         string            `json:"mi"`
	INO        string            `json:"ino"`
	LOS        string            `json:"los"`
	RENALINO   string            `json:"renalino"`
	POSTRF     string            `json:"postrf"`
	PACE       string            `json:"pace"`
	OCVENDYS   string            `json:"ocvendys"`
	AFIB       string            `json:"afib"`
	OCDVT      string            `json:"ocdvt"`
	OCPULMC    string            `json:"ocpulmc"`
	SEIZURES   string            `json:"seizures"`
	TIA        string            `json:"tia"`
	PREHB      string            `json:"prehb"`
	POSTHB     string            `json:"posthb"`
	RBC        string            `json:"rbc"`
	NOBED      string            `json:"nobed"`
	NONURSE    string            `json:"nonurse"`
	ICUCOMP    string            `json:"icucomp"`
	CHRPTS     string            `json:"chrpts"`
	OTHER      string            `json:"other"`
	OTHERNOTE  string            `json:"othernote"`
	STROKE     string            `json:"stroke"`
	INFARM     string            `json:"infarm"`
	INFLEG     string            `json:"infleg"`
	INFSTERN   string            `json:"infstern"`
	INFSEP     string            `json:"infsep"`
	SURVIVAL   string            `json:"survival"`
	DCTO       string            `json:"dcto"`
	PROC       string            `json:"proc"`
	NOTES      string            `json:"notes"`
	DRUG4      string            `json:"drug4"`
	DRUG5      string            `json:"drug5"`
	DRUG6      string            `json:"drug6"`
	CORMATRIX  string            `json:"cormatrix"`
	CELLSAVER  string            `json:"cellsaver"`
	FIX        []periopcheck.Fix `json:"fix"`
}

func writeJSON(newEvent Event, jsonFile *os.File) {
	j, err := json.Marshal(newEvent)
	if err != nil {
		log.Println(err)
	}
	jsonFile.Write(j)
}

func checkDate() {
	// helper.CheckDateFormat()
}

func checkRow(rowSlice map[string]string, rowNum int) {
	var fixArray []periopcheck.Fix
	var fix periopcheck.Fix

	// for loop for all binary codes
	binaryCodeArray := []string{"AREA", "TRIANGE", "SDA", "CATRIAL", "ANTRIAL", "PITHROMB", "DIABETES", "HYPER", "CHLSTRL", "FHX", "COPD", "COPDS", "THROMB", "NEWRF", "DIAL", "MARFAN", "CHF", "SHOCK", "SYNCOPE", "ASP", "AMI", "STATIN", "GORTEX", "SEPT", "MAZE", "DISLAD", "DISCX", "DISRCA", "LMAIN", "SKEGRAFT", "MININV", "MI", "INO", "RENALINO", "LOS", "PACE", "OCVENDYS", "AFIB", "OCDVT", "OCPULMC", "SEIZURES", "TIA", "INFARM", "INFSEP", "SURVIVAL"}
	nonNegativeArray := []string{"DAYSPOST", "ICUNUM", "ICU", "VENT", "CREAT", "AVSIZE", "MVSIZE", "TVSIZE", "PVSIZE", "CI", "MPAP", "SYSAVG", "LVEDP", "PVR", "MVGRADR", "AVAREA", "MVAREA", "AVXPLSIZE", "MVXPLSIZE", "TVXPLSIZE", "DISNUM", "GFTLAD", "GFTCX", "GFTRCA", "ACBNUM", "ORTIME", "PUMP", "CLAMP", "CIRARR", "HT", "WT", "REOPNUM", "CK", "CKMB", "PREHB", "POSTHB", "PACKCELLS"}
	// nameArray := []string{"SURG", "ASSIST", "FDOC", "CDOC"}
	var multiFieldArray []MultiField
	var fieldArray []Field
	for _, b := range binaryCodeArray {
		switch periopcheck.CheckValidNumber(0, 2, rowSlice[b]) {
		case 0:
			continue
		case 1:
			fix = periopcheck.CantReadErrorHandler(rowNum, b, rowSlice)
			fixArray = append(fixArray, fix)
		case 2:
			rowSlice[b] = "-9"
		case 3:
			fix = periopcheck.OutBoundsErrorHandler(rowNum, b, rowSlice)
			fixArray = append(fixArray, fix)
		}
	}

	for _, n := range nonNegativeArray {
		if !periopcheck.CheckNonNegative(rowSlice[n]) {
			fix = periopcheck.OutBoundsErrorHandler(rowNum, n, rowSlice)
			fixArray = append(fixArray, fix)
		}
	}
	//DEMOGRAPHICS:
	f := Field{"AREA", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"SEX", 0, 2}
	fieldArray = append(fieldArray, f)

	//DATES & DOCTORS:

	if !periopcheck.CheckNonNegativeFloat(rowSlice["ICUTIME"]) {
		periopcheck.OutBoundsErrorHandler(rowNum, "ICUTIME", rowSlice)
	}
	if !periopcheck.CheckNonNegativeFloat(rowSlice["VENTIME"]) {
		periopcheck.OutBoundsErrorHandler(rowNum, "VENTIME", rowSlice)
	}
	//GENERAL PATIENT DATA :
	f = Field{"TIMING", 1, 4}
	fieldArray = append(fieldArray, f)
	f = Field{"FROM", 1, 5}
	fieldArray = append(fieldArray, f)

	//****************************************************************
	//MultiField ACBREDO
	multiField := MultiField{"ACBREDOP", 0, 1}
	multiFieldArray = append(multiFieldArray, multiField)
	//MultiField AVREDO
	multiField = MultiField{"AVREDOP", 0, 3}
	multiFieldArray = append(multiFieldArray, multiField)
	//MultiField MVREDO
	multiField = MultiField{"MVREDOP", 0, 3}
	multiFieldArray = append(multiFieldArray, multiField)
	//MultiField TVREDO
	multiField = MultiField{"TVREDOP", 0, 3}
	multiFieldArray = append(multiFieldArray, multiField)
	//MultiField OTHEREDO
	multiField = MultiField{"OTHEREDO", 0, 1}
	multiFieldArray = append(multiFieldArray, multiField)
	//****************************************************************
	//PREVIOUS (NON-SURGICAL) INTERVENTION:
	f = Field{"PRECARD", 1, 2}
	fieldArray = append(fieldArray, f)

	//CLINICAL PRESENTATION:
	f = Field{"ANGINA", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"PREOPMI", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"NYHA", 1, 4}
	fieldArray = append(fieldArray, f)
	if !periopcheck.CheckCCS(rowSlice["CCS"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "CCS", rowSlice)
		fixArray = append(fixArray, fix)
	}
	f = Field{"LVGRADE", 1, 4}
	fieldArray = append(fieldArray, f)
	f = Field{"STRESS", 0, 2}
	fieldArray = append(fieldArray, f)

	//C.A.D. RISKS:
	f = Field{"DOI", 0, 4}
	fieldArray = append(fieldArray, f)
	f = Field{"SMOKE", 0, 2}
	fieldArray = append(fieldArray, f)

	//ASSOCIATED DISEASES:
	if !periopcheck.CheckPVD(rowSlice["PVD"], rowSlice["CORATID"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "PVD", rowSlice)
		fixArray = append(fixArray, fix)
	}
	f = Field{"RF", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"CAROTID", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"ADD", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"RECG", 0, 2}
	fieldArray = append(fieldArray, f)

	//VALVE PATIENT DATA:
	f = Field{"AVDIS", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"ENDOCARD", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"URGENT", 0, 7}
	fieldArray = append(fieldArray, f)
	f = Field{"MVDIS", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"TVDIS", 0, 3}
	fieldArray = append(fieldArray, f)

	//AORTIC VALVE SURGERY:
	f = Field{"AVSURG", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"D2", 0, 4}
	fieldArray = append(fieldArray, f)
	f = Field{"ANNULEN", 0, 3}
	fieldArray = append(fieldArray, f)

	//MultiField AVPATH
	multiField = MultiField{"AVPATH", 0, 8}
	multiFieldArray = append(multiFieldArray, multiField)

	if !periopcheck.CheckVPROS(rowSlice["AVPROS"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "AVPROS", rowSlice)
		fixArray = append(fixArray, fix)
	}
	//MITRAL VALVE SURGERY:
	f = Field{"MVSURG", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"AREA", 0, 2}
	fieldArray = append(fieldArray, f)

	//MultiField MVPATH
	multiField = MultiField{"MVPATH", 0, 7}
	multiFieldArray = append(multiFieldArray, multiField)

	f = Field{"MVANN", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"CHORD", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"MVP", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"MVC", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"CHORDAL", 0, 3}
	fieldArray = append(fieldArray, f)

	if !periopcheck.CheckVPROS(rowSlice["MVPROS"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "MVPROS", rowSlice)
		fixArray = append(fixArray, fix)
	}
	//TRICUSPID VALVE SURGERY:
	f = Field{"TVSURG", 0, 3}
	fieldArray = append(fieldArray, f)

	//MultiField TVPATH
	multiField = MultiField{"TVPATH", 0, 5}
	multiFieldArray = append(multiFieldArray, multiField)

	if !periopcheck.CheckVPROS(rowSlice["TVPROS"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "TVPROS", rowSlice)
		fixArray = append(fixArray, fix)
	}
	//PULMONARY VALVE SURGERY:
	f = Field{"PVSURG", 0, 3}
	fieldArray = append(fieldArray, f)

	if !periopcheck.CheckVPROS(rowSlice["PVPROS"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "PVPROS", rowSlice)
		fixArray = append(fixArray, fix)
	}
	//EXPLANTED VALVES:
	//AORTIC:
	if !periopcheck.CheckVPROS(rowSlice["AVXPLTYPE"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "AVXPLTYPE", rowSlice)
		fixArray = append(fixArray, fix)
	}
	f = Field{"AVXPLTYPE", 0, 6}
	fieldArray = append(fieldArray, f)

	//MITRAL:
	if !periopcheck.CheckVPROS(rowSlice["MVXPLTYPE"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "MVXPLTYPE", rowSlice)
		fixArray = append(fixArray, fix)
	}
	f = Field{"MVXPLTYPE", 0, 6}
	fieldArray = append(fieldArray, f)

	//TRICUSPID:
	if !periopcheck.CheckVPROS(rowSlice["TVXPLTYPE"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "TVXPLTYPE", rowSlice)
		fixArray = append(fixArray, fix)
	}
	f = Field{"TVXPLPATH", 0, 6}
	fieldArray = append(fieldArray, f)

	//OTHER PROCEDURES/OPERATIONS:
	f = Field{"LVA", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"SEPTYPE", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"CHD", 0, 4}
	fieldArray = append(fieldArray, f)
	f = Field{"AAS", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"AOPATH", 0, 4}
	fieldArray = append(fieldArray, f)
	f = Field{"MISC", 1, 11}
	fieldArray = append(fieldArray, f)
	f = Field{"OTHERTYPE", 0, 5}
	fieldArray = append(fieldArray, f)
	f = Field{"DEVICETYPE", 0, 5}
	fieldArray = append(fieldArray, f)

	//ACB PATIENT DATA:
	//ALL BINARY OR NON negative

	//ARTERIES BYPASSED:
	f = Field{"LIMA", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"RIMA", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"RADIAL", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"ENDART", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"CVPA", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"OTHGFT", 0, 3}
	fieldArray = append(fieldArray, f)

	//GENERAL OPERATIVE DATA:
	f = Field{"PUMPCASE", 0, 2}
	fieldArray = append(fieldArray, f)
	if !periopcheck.CheckNonNegativeFloat(rowSlice["BSA"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "BSA", rowSlice)
		fixArray = append(fixArray, fix)
	}
	f = Field{"MYOPRO", 0, 4}
	fieldArray = append(fieldArray, f)
	f = Field{"TECH", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"DIRECT", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"HYPOTHER", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"OFFPUMP", 0, 2}
	fieldArray = append(fieldArray, f)

	//COMPLICATIONS:
	f = Field{"IABP", 0, 4}
	fieldArray = append(fieldArray, f)
	//MultiField REOP
	multiField = MultiField{"REOP", 0, 7}
	multiFieldArray = append(multiFieldArray, multiField)
	//MultiField REOPPUMP
	multiField = MultiField{"REOPPUMP", 0, 2}
	multiFieldArray = append(multiFieldArray, multiField)

	f = Field{"IECG", 0, 4}
	fieldArray = append(fieldArray, f)
	f = Field{"POSTRF", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"STROKE", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"INFLEG", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"INFSTERN", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"DCTO", 0, 3}
	fieldArray = append(fieldArray, f)
	f = Field{"PROC", 1, 3}
	fieldArray = append(fieldArray, f)
	//END
	//check all f in fieldArray
	for i := range fieldArray {
		switch periopcheck.CheckValidNumber(fieldArray[i].Min, fieldArray[i].Max, fieldArray[i].Name) {
		case 0:
			continue
		case 1:
			fix = periopcheck.CantReadErrorHandler(rowNum, fieldArray[i].Name, rowSlice)
			fixArray = append(fixArray, fix)

		case 2:
			rowSlice[fieldArray[i].Name] = "-9"
		case 3:
			fix = periopcheck.OutBoundsErrorHandler(rowNum, fieldArray[i].Name, rowSlice)
			fixArray = append(fixArray, fix)

		}
	}

	for _, f := range multiFieldArray {
		combineMultiFields(f.Name, rowSlice, rowNum, f.Min, f.Max)
	}
}

//combineMultiFields looks for field names containing a multifield and appends them to a slice of int
func combineMultiFields(fieldName string, rowSlice map[string]string, rowNum int, min int, max int) []int {
	var multiField []int
	var fixArray []periopcheck.Fix
	var fix periopcheck.Fix
	for field := range rowSlice {
		if strings.Contains(field, fieldName) {
			switch periopcheck.CheckValidNumber(min, max, rowSlice[field]) {
			case 0:
				continue
			case 1:
				fix = periopcheck.CantReadErrorHandler(rowNum, field, rowSlice)
				fixArray = append(fixArray, fix)
			case 2:
				rowSlice[field] = "-9"
			case 3:
				fix = periopcheck.OutBoundsErrorHandler(rowNum, field, rowSlice)
				fixArray = append(fixArray, fix)
			}
			value := excelHelper.StringToInt(rowSlice[field])
			multiField = append(multiField, value)
		}
	}
	return multiField
}

//uniformDates converts all dates to YYYY-MM-DD
func uniformDates(rowSlice map[string]string, row int) map[string]string {
	for value := range rowSlice {
		value = strings.ToLower(value)
		if strings.Contains(value, "date") {
			date := helper.CheckDateFormat(row, value, rowSlice[value])
			rowSlice[value] = date
		}
	}
	return rowSlice
}

//parseData reads the sheet and inserts the cell values into a map with the column name as the key
func parseData(sheet *xlsx.File) {
	rowSlice := make(map[string]string)
	colLength := sheet.Sheets[0].MaxCol
	rowLength := sheet.Sheets[0].MaxRow
	for ri := 1; ri < rowLength; ri++ {
		for ci := 0; ci < colLength; ci++ {
			rowSlice[sheet.Sheets[0].Rows[0].Cells[ci].Value] = sheet.Sheets[0].Rows[ri].Cells[ci].Value
		}
	}
}

//attributes returns the name of the fields of a particular struct
func attributes(m interface{}) []string {
	typ := reflect.TypeOf(m)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// create an attribute data structure as a map of types keyed by a string.
	var attrs []string
	// Only structs are supported so return an empty result if the passed object
	// isn't a struct
	if typ.Kind() != reflect.Struct {
		fmt.Printf("%v type can't have attributes inspected\n", typ.Kind())
		return attrs
	}

	// loop through the struct's fields and set the map
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		if !p.Anonymous {
			attrs = append(attrs, p.Name)
		}
	}

	return attrs
}

func main() {

	dirTGH := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TGH perioperative.xlsx"
	// dirTWH := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TWH perioperative.xlsx"
	tghFile := excelHelper.ConnectToXlsx(dirTGH)
	tghCols := excelHelper.IdentifyCols(tghFile)
	// twhFile := excelHelper.ConnectToXlsx(dirTWH)
	// twhCols := excelHelper.IdentifyCols(twhFile)
	excelHelper.WriteStruct(tghCols)
}
