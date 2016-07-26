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

//TODO VPROS combine into one array and add null when empty
//TODO somehow parse data into struct
//TODO create test program

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
	PERIOPID  string  `json:"periop_id"`
	PTID      string  `json:"ptid"`
	AREA      int64   `json:"area"`
	TRIAGE    int64   `json:"triage"`
	SDA       int64   `json:"sda"`
	ADDATE    string  `json:"addate"`
	DATEOR    string  `json:"dateor"`
	DISDATE   string  `json:"disdate"`
	DAYSPOST  int64   `json:"dayspost"`
	ICUNUM    int64   `json:"icunum"`
	ICU       int64   `json:"icu"`
	ICUTIME   float64 `json:"icutime"`
	VENT      int64   `json:"vent"`
	VENTTIME  float64 `json:"venttime"`
	SURG      string  `json:"surg"`
	ASSIST    string  `json:"assist"`
	FDOC      string  `json:"fdoc"`
	CDOC      string  `json:"cdoc"`
	CATRIAL   int64   `json:"catrial"`
	ANTRIAL   int64   `json:"antrial"`
	TIMING    int64   `json:"timing"`
	FROM      int64   `json:"from"`
	ACBREDO   int64   `json:"acbredo"`
	AVREDO    int64   `json:"avredo"`
	MVREDO    int64   `json:"mvredo"`
	TVREDO    int64   `json:"tvredo"`
	OTHREDO   int64   `json:"othredo"`
	DATEORP1  string  `json:"dateorp1"`
	ACBREDOP1 string  `json:"acbredop1"`
	AVREDOP1  string  `json:"avredop1"`
	MVREDOP1  string  `json:"mvredop1"`
	TVREDOP1  string  `json:"tvredop1"`
	OTHREDOP1 string  `json:"othredop1"`
	DATEORP2  string  `json:"dateorp2"`
	ACBREDOP2 string  `json:"acbredop2"`
	AVREDOP2  string  `json:"avredop2"`
	MVREDOP2  string  `json:"mvredop2"`
	TVREDOP2  string  `json:"tvredop2"`
	OTHREDOP2 string  `json:"othredop2"`
	DATEORP3  string  `json:"dateorp3"`
	ACBREDOP3 string  `json:"acbredop3"`
	AVREDOP3  string  `json:"avredop3"`
	MVREDOP3  string  `json:"mvredop3"`
	TVREDOP3  string  `json:"tvredop3"`
	OTHREDOP3 string  `json:"othredop3"`
	DATEORP4  string  `json:"dateorp4"`
	ACBREDOP4 string  `json:"acbredop4"`
	AVREDOP4  string  `json:"avredop4"`
	MVREDOP4  string  `json:"mvredop4"`
	TVREDOP4  string  `json:"tvredop4"`
	OTHREDOP4 string  `json:"othredop4"`
	DATEORP5  string  `json:"dateorp5"`
	ACBREDOP5 string  `json:"acbredop5"`
	AVREDOP5  string  `json:"avredop5"`
	MVREDOP5  string  `json:"mvredop5"`
	TVREDOP5  string  `json:"tvredop5"`
	OTHREDOP5 string  `json:"othredop5"`
	DATEORP6  string  `json:"dateorp6"`
	ACBREDOP6 string  `json:"acbredop6"`
	AVREDOP6  string  `json:"avredop6"`
	MVREDOP6  string  `json:"mvredop6"`
	TVREDOP6  string  `json:"tvredop6"`
	OTHREDOP6 string  `json:"othredop6"`
	PRECARD   int64   `json:"precard"`
	PIDATE    string  `json:"pidate"`
	PITHROMB  int64   `json:"pithromb"`
	PITDATE   string  `json:"pitdate"`
	CATH      string  `json:"cath"`
	CATHDATE  string  `json:"cathdate"`
	ANGINA    int64   `json:"angina"`
	PREOPMI   int64   `json:"preopmi"`
	MIDATE    string  `json:"midate"`
	NYHA      int64   `json:"nyha"`
	CCS       string  `json:"ccs"`
	LVGRADE   int64   `json:"lvgrade"`
	STRESS    int64   `json:"stress"`
	DIABETES  int64   `json:"diabetes"`
	DOI       int64   `json:"doi"`
	HYPER     int64   `json:"hyper"`
	CHLSTRL   int64   `json:"chlstrl"`
	FHX       int64   `json:"fhx"`
	SMOKE     int64   `json:"smoke"`
	PACKS     int64   `json:"packs"`
	COPD      int64   `json:"copd"`
	COPDS     int64   `json:"copds"`
	THROMB    int64   `json:"thromb"`
	PVD       int64   `json:"pvd"`
	RF        int64   `json:"rf"`
	NEWRF     int64   `json:"newrf"`
	DIAL      int64   `json:"dial"`
	MARFAN    int64   `json:"marfan"`
	CAROTID   int64   `json:"carotid"`
	AAD       int64   `json:"aad"`
	EU        int64   `json:"eu"`
	RECG      int64   `json:"recg"`
	CHF       int64   `json:"chf"`
	SHOCK     int64   `json:"shock"`
	SYNCOPE   int64   `json:"syncope"`
	ASP       int64   `json:"asp"`
	AMI       int64   `json:"ami"`
	CREAT     int64   `json:"creat"`
	STATIN    int64   `json:"statin"`
	AVDIS     int64   `json:"avdis"`
	MVDIS     int64   `json:"mvdis"`
	TVDIS     int64   `json:"tvdis"`
	ENDOCARD  int64   `json:"endocard"`
	URGENT    int64   `json:"urgent"`
	AVSURG    int64   `json:"avsurg"`
	D2        int64   `json:"d2"`
	AVSIZE    int64   `json:"avsize"`
	AVPATH    []int64 `json:"avpath"`

	ANNULEN int64   `json:"annulen"`
	AVPROS  string  `json:"avpros"`
	MVSURG  int64   `json:"mvsurg"`
	MVSIZE  int64   `json:"mvsize"`
	MVPATH  []int64 `json:"mvpath"`

	MVANN   int64   `json:"mvann"`
	CHORD   int64   `json:"chord"`
	GORTEX  int64   `json:"gortex"`
	MVP     int64   `json:"mvp"`
	MVC     int64   `json:"mvc"`
	CHORDAL int64   `json:"chordal"`
	MVPROS  string  `json:"mvpros"`
	TVSURG  int64   `json:"tvsurg"`
	TVSIZE  int64   `json:"tvsize"`
	TVPATH  []int64 `json:"tvpath"`

	TVPROS     string  `json:"tvpros"`
	PVSURG     int64   `json:"pvsurg"`
	PVSIZE     int64   `json:"pvsize"`
	PVPROS     string  `json:"pvpros"`
	CI         int64   `json:"ci"`
	MPAP       int64   `json:"mpap"`
	SYSAVG     int64   `json:"sysavg"`
	LVEDP      int64   `json:"lvedp"`
	PVR        int64   `json:"pvr"`
	MVGRADR    int64   `json:"mvgradr"`
	AVAREA     int64   `json:"avarea"`
	MVAREA     int64   `json:"mvarea"`
	AVXPLTYPE  string  `json:"avxpltype"`
	AVXPLSIZE  int64   `json:"avxplsize"`
	AVXPLPATH  int64   `json:"avxplpath"`
	AVXPLDATE  string  `json:"avxpldate"`
	MVXPLTYPE  string  `json:"mvxpltype"`
	MVXPLSIZE  int64   `json:"mvxplsize"`
	MVXPLPATH  int64   `json:"mvxplpath"`
	MVXPLDATE  string  `json:"mvxpldate"`
	TVXPLTYPE  string  `json:"tvxpltype"`
	TVXPLSIZE  int64   `json:"tvxplsize"`
	TVXPLPATH  int64   `json:"tvxplpath"`
	TVXPLDATE  string  `json:"tvxpldate"`
	ASSOCOP    int64   `json:"assocop"`
	LVA        int64   `json:"lva"`
	SEPT       int64   `json:"sept"`
	SEPTYPE    int64   `json:"septype"`
	CHD        int64   `json:"chd"`
	AAS        int64   `json:"aas"`
	AOPATH     int64   `json:"aopath"`
	MAZE       int64   `json:"maze"`
	MISC       int64   `json:"misc"`
	OTHERTYPE  int64   `json:"othertype"`
	ASSDEV     int64   `json:"assdev"`
	DEVICETYPE int64   `json:"devicetype"`
	DISLAD     int64   `json:"dislad"`
	DISCX      int64   `json:"discx"`
	DISRCA     int64   `json:"disrca"`
	LMAIN      int64   `json:"lmain"`
	DISNUM     int64   `json:"disnum"`
	LIMA       int64   `json:"lima"`
	RIMA       int64   `json:"rima"`
	RADIAL     int64   `json:"radial"`
	SKEGRAFT   int64   `json:"skegraft"`
	GFTLAD     int64   `json:"gftlad"`
	GFTCX      int64   `json:"gftcx"`
	GFTRCA     int64   `json:"gftrca"`
	ENDART     int64   `json:"endart"`
	CVPA       int64   `json:"cvpa"`
	OTHGFT     int64   `json:"othgft"`
	ACBNUM     int64   `json:"acbnum"`
	PUMPCASE   int64   `json:"pumpcase"`
	MININV     int64   `json:"mininv"`
	ORTIME     int64   `json:"ortime"`
	PUMP       int64   `json:"pump"`
	CLAMP      int64   `json:"clamp"`
	CIRARR     int64   `json:"cirarr"`
	BSA        float64 `json:"bsa"`
	HT         int64   `json:"ht"`
	WT         int64   `json:"wt"`
	MYOPRO     int64   `json:"myopro"`
	TECH       int64   `json:"tech"`
	DIRECT     int64   `json:"direct"`
	HYPOTHER   int64   `json:"hypother"`
	OFFPUMP    int64   `json:"offpump"`
	IABP       int64   `json:"iabp"`
	REOPNUM    int64   `json:"reopnum"`
	REOP       []int64 `json:"reop"`

	REOPPUMP []int64 `json:"reoppump"`

	IECG      int64  `json:"iecg"`
	CK        int64  `json:"ck"`
	CKMB      int64  `json:"ckmb"`
	MI        int64  `json:"mi"`
	INO       int64  `json:"ino"`
	LOS       int64  `json:"los"`
	RENALINO  int64  `json:"renalino"`
	POSTRF    int64  `json:"postrf"`
	PACE      int64  `json:"pace"`
	OCVENDYS  int64  `json:"ocvendys"`
	AFIB      int64  `json:"afib"`
	OCDVT     int64  `json:"ocdvt"`
	OCPULMC   int64  `json:"ocpulmc"`
	SEIZURES  int64  `json:"seizures"`
	TIA       int64  `json:"tia"`
	PREHB     int64  `json:"prehb"`
	POSTHB    int64  `json:"posthb"`
	RBC       int64  `json:"rbc"`
	NOBED     int64  `json:"nobed"`
	NONURSE   int64  `json:"nonurse"`
	ICUCOMP   int64  `json:"icucomp"`
	CHRPTS    int64  `json:"chrpts"`
	OTHER     int64  `json:"other"`
	OTHERNOTE string `json:"othernote"`
	STROKE    int64  `json:"stroke"`
	INFARM    int64  `json:"infarm"`
	INFLEG    int64  `json:"infleg"`
	INFSTERN  int64  `json:"infstern"`
	INFSEP    int64  `json:"infsep"`
	SURVIVAL  int64  `json:"survival"`
	DCTO      int64  `json:"dcto"`
	PROC      int64  `json:"proc"`

	FIX []periopcheck.Fix `json:"fix"`
}

//addToStruct an updated row and parses them into the correct struct field by matching the field names with map keys
//the function also checks whether the field is an int64 or string and converts the mapped value accordingly
func addToStruct(rowSlice map[string]string, newPeriOp PeriOp) PeriOp {
	skipFields := []string{"REDOP", "DATEORP", "VPATH", "REOP", "REOPPUMP", "FIX"}
	//iterate through each field of the Struct
	for i := 0; i < reflect.TypeOf(newPeriOp).NumField(); i++ {
		//sf is the field
		field := reflect.TypeOf(newPeriOp).Field(i)
		constainsSkip := false
		//checks if the field should be skipped
		for _, skip := range skipFields {
			if strings.Contains(field.Name, skip) {
				constainsSkip = true
			}
		}
		//if the field doesnt need to be skipped then check its type and insert value accordingly
		if !constainsSkip {
			if field.Type.Name() == "int64" {
				valueInt := excelHelper.StringToInt(rowSlice[field.Name])
				reflect.ValueOf(&newPeriOp).Elem().Field(i).SetInt(valueInt)

			} else if field.Type.Name() == "string" {
				valueStr := rowSlice[field.Name]
				reflect.ValueOf(&newPeriOp).Elem().Field(i).SetString(valueStr)

			}
		}
	}
	return newPeriOp
}

//addMultitoStruct adds the array of codes to the struct's multifield
func addMultitoStruct(newPeriOp PeriOp, insert []int64, fieldName string) PeriOp {
	insertValue := reflect.ValueOf(insert)
	reflect.ValueOf(&newPeriOp).Elem().FieldByName(fieldName).Set(insertValue)
	return newPeriOp
}

//addFixtoStruct adds the Fix array to the Fix field
func addFixtoStruct(newPeriOp PeriOp, fixArray []periopcheck.Fix) PeriOp {
	newPeriOp.FIX = fixArray
	return newPeriOp
}

//writeJSON writes the struct into JSON format
func writeJSON(newEvent Event, jsonFile *os.File) {
	j, err := json.Marshal(newEvent)
	if err != nil {
		log.Println(err)
	}
	jsonFile.Write(j)
}

//fixDate takes in a date and puts it into YYYY-MM-DD format
func fixDate(rowNum int, fieldName string, date string) string {
	return helper.CheckDateFormat(rowNum, fieldName, date)
}

func checkRow(rowSlice map[string]string, rowNum int) (map[string]string, []periopcheck.Fix) {
	var fixArray []periopcheck.Fix

	return rowSlice, fixArray
}

//checkStringCodes checks if the string codes are empty
func checkStringCodes(fixArray []periopcheck.Fix, rowSlice map[string]string, rowNum int) {
	var fix periopcheck.Fix
	nameArray := []string{"SURG", "ASSIST", "FDOC", "CDOC"}
	for _, name := range nameArray {
		if rowSlice[name] == "" {
			fix.Field = "code"
			fix.Msg = "invalid code:empty"
		}
	}
}

//checkBinaryCodes checks the binary codes
func checkBinaryCodes(fixArray []periopcheck.Fix, rowSlice map[string]string, rowNum int) {
	binaryCodeArray := []string{"AREA", "TRIANGE", "SDA", "CATRIAL", "ANTRIAL", "PITHROMB", "DIABETES", "HYPER", "CHLSTRL", "FHX", "COPD", "COPDS", "THROMB", "NEWRF", "DIAL", "MARFAN", "CHF", "SHOCK", "SYNCOPE", "ASP", "AMI", "STATIN", "GORTEX", "SEPT", "MAZE", "DISLAD", "DISCX", "DISRCA", "LMAIN", "SKEGRAFT", "MININV", "MI", "INO", "RENALINO", "LOS", "PACE", "OCVENDYS", "AFIB", "OCDVT", "OCPULMC", "SEIZURES", "TIA", "INFARM", "INFSEP", "SURVIVAL"}
	var fix periopcheck.Fix
	// For Binary Codes
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
}

//checkNonNegCodes checks if the code is negative
func checkNonNegCodes(fixArray []periopcheck.Fix, rowSlice map[string]string, rowNum int) {
	var fix periopcheck.Fix
	nonNegativeArray := []string{"DAYSPOST", "ICUNUM", "ICU", "VENT", "CREAT", "AVSIZE", "MVSIZE", "TVSIZE", "PVSIZE", "CI", "MPAP", "SYSAVG", "LVEDP", "PVR", "MVGRADR", "AVAREA", "MVAREA", "AVXPLSIZE", "MVXPLSIZE", "TVXPLSIZE", "DISNUM", "GFTLAD", "GFTCX", "GFTRCA", "ACBNUM", "ORTIME", "PUMP", "CLAMP", "CIRARR", "HT", "WT", "REOPNUM", "CK", "CKMB", "PREHB", "POSTHB", "PACKCELLS"}
	//For non negative codes
	for _, n := range nonNegativeArray {
		if !periopcheck.CheckNonNegative(rowSlice[n]) {
			fix = periopcheck.OutBoundsErrorHandler(rowNum, n, rowSlice)
			fixArray = append(fixArray, fix)
		}
	}
}

//checkNonMultiCode checks all non-multi codes
func checkNonMultiCode(fixArray []periopcheck.Fix, rowSlice map[string]string, rowNum int) {
	var fieldArray []Field
	var fix periopcheck.Fix
	//DEMOGRAPHICS:
	f := Field{"AREA", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"SEX", 0, 2}
	fieldArray = append(fieldArray, f)
	//GENERAL PATIENT DATA :
	f = Field{"TIMING", 1, 4}
	fieldArray = append(fieldArray, f)
	f = Field{"FROM", 1, 5}
	fieldArray = append(fieldArray, f)
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
	//MITRAL VALVE SURGERY:
	f = Field{"MVSURG", 0, 2}
	fieldArray = append(fieldArray, f)
	f = Field{"AREA", 0, 2}
	fieldArray = append(fieldArray, f)
	//MultiField MVPATH
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
	//TRICUSPID VALVE SURGERY:
	f = Field{"TVSURG", 0, 3}
	fieldArray = append(fieldArray, f)
	//PULMONARY VALVE SURGERY:
	f = Field{"PVSURG", 0, 3}
	fieldArray = append(fieldArray, f)
	//EXPLANTED VALVES:
	//AORTIC:
	f = Field{"AVXPLTYPE", 0, 6}
	fieldArray = append(fieldArray, f)
	//MITRAL:
	f = Field{"MVXPLTYPE", 0, 6}
	fieldArray = append(fieldArray, f)
	//TRICUSPID:
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
	//For all other non-multi field codes
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
}

//checkMultiCodes checks multi code and creates an array filled with fields with multi codes
func checkMultiCodes(fixArray []periopcheck.Fix, rowSlice map[string]string, rowNum int) {
	var multiFieldArray []MultiField
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
	//MultiField AVPATH
	multiField = MultiField{"AVPATH", 0, 8}
	multiFieldArray = append(multiFieldArray, multiField)
	//MultiField MVPATH
	multiField = MultiField{"MVPATH", 0, 7}
	multiFieldArray = append(multiFieldArray, multiField)
	//MultiField TVPATH
	multiField = MultiField{"TVPATH", 0, 5}
	multiFieldArray = append(multiFieldArray, multiField)
	//MultiField REOP
	multiField = MultiField{"REOP", 0, 7}
	multiFieldArray = append(multiFieldArray, multiField)
	//MultiField REOPPUMP
	multiField = MultiField{"REOPPUMP", 0, 2}
	multiFieldArray = append(multiFieldArray, multiField)
	//For multi field codes
	for _, f := range multiFieldArray {
		combineMultiFields(f.Name, rowSlice, rowNum, f.Min, f.Max)
	}
}

//checkNonNegFloats checks non-negative floats
func checkNonNegFloats(fixArray []periopcheck.Fix, rowSlice map[string]string, rowNum int) {
	var fix periopcheck.Fix
	//DATES & DOCTORS:
	if !periopcheck.CheckNonNegativeFloat(rowSlice["ICUTIME"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "ICUTIME", rowSlice)
		fixArray = append(fixArray, fix)
	}
	if !periopcheck.CheckNonNegativeFloat(rowSlice["VENTIME"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "VENTIME", rowSlice)
		fixArray = append(fixArray, fix)
	}
	//GENERAL OPERATIVE DATA:
	if !periopcheck.CheckNonNegativeFloat(rowSlice["BSA"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "BSA", rowSlice)
		fixArray = append(fixArray, fix)
	}
}

//checkVPROSCodes checks if all the VPROS codes are valid
func checkVPROSCodes(fixArray []periopcheck.Fix, rowSlice map[string]string, rowNum int) {
	var fix periopcheck.Fix
	//ASSOCIATED DISEASES:
	if !periopcheck.CheckVPROS(rowSlice["AVPROS"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "AVPROS", rowSlice)
		fixArray = append(fixArray, fix)
	}

	if !periopcheck.CheckVPROS(rowSlice["MVPROS"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "MVPROS", rowSlice)
		fixArray = append(fixArray, fix)
	}

	if !periopcheck.CheckVPROS(rowSlice["TVPROS"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "TVPROS", rowSlice)
		fixArray = append(fixArray, fix)
	}

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

	//MITRAL:
	if !periopcheck.CheckVPROS(rowSlice["MVXPLTYPE"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "MVXPLTYPE", rowSlice)
		fixArray = append(fixArray, fix)
	}

	//TRICUSPID:
	if !periopcheck.CheckVPROS(rowSlice["TVXPLTYPE"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "TVXPLTYPE", rowSlice)
		fixArray = append(fixArray, fix)
	}
}

//checkCCSCodes checks if CCS code is valid
func checkCCSCodes(fixArray []periopcheck.Fix, rowSlice map[string]string, rowNum int) {
	var fix periopcheck.Fix
	//CLINICAL PRESENTATION:
	if !periopcheck.CheckCCS(rowSlice["CCS"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "CCS", rowSlice)
		fixArray = append(fixArray, fix)
	}
}

//checkPVDCode checks if PVD code is valid
func checkPVDCode(fixArray []periopcheck.Fix, rowSlice map[string]string, rowNum int) {
	var fix periopcheck.Fix
	//ASSOCIATED DISEASES:
	if !periopcheck.CheckPVD(rowSlice["PVD"], rowSlice["CORATID"]) {
		fix = periopcheck.OutBoundsErrorHandler(rowNum, "PVD", rowSlice)
		fixArray = append(fixArray, fix)
	}
}

//combineMultiFields looks for field names containing a multifield and appends them to a slice of int
func combineMultiFields(fieldName string, rowSlice map[string]string, rowNum int, min int, max int) []int64 {
	var multiField []int64
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
func uniformDates(fixArray []periopcheck.Fix, rowSlice map[string]string, row int) map[string]string {
	for value := range rowSlice {
		value = strings.ToLower(value)
		if strings.Contains(value, "date") {
			if rowSlice[value] == "" {
				rowSlice[value] = "null"

			}
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
