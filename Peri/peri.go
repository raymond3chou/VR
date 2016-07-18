package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/access/excelHelper"
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

//PeriOp is the object for periop data
type PeriOp struct {
	PERIOPID   string `json:"periop_id"`
	PTID       string `json:"ptid"`
	AREA       string `json:"area"`
	TRIAGE     string `json:"triage"`
	SDA        string `json:"sda"`
	ADDATE     string `json:"addate"`
	DATEOR     string `json:"dateor"`
	DISDATE    string `json:"disdate"`
	DAYSPOST   string `json:"dayspost"`
	ICUNUM     string `json:"icunum"`
	ICU        string `json:"icu"`
	ICUTIME    string `json:"icutime"`
	VENT       string `json:"vent"`
	VENTTIME   string `json:"venttime"`
	SURG       string `json:"surg"`
	ASSIST     string `json:"assist"`
	FDOC       string `json:"fdoc"`
	CDOC       string `json:"cdoc"`
	CATRIAL    string `json:"catrial"`
	ANTRIAL    string `json:"antrial"`
	FROM       string `json:"from"`
	ACBREDO    string `json:"acbredo"`
	AVREDO     string `json:"avredo"`
	MVREDO     string `json:"mvredo"`
	TVREDO     string `json:"tvredo"`
	OTHREDO    string `json:"othredo"`
	DATEORP1   string `json:"dateorp1"`
	ACBREDOP1  string `json:"acbredop1"`
	AVREDOP1   string `json:"avredop1"`
	MVREDOP1   string `json:"mvredop1"`
	TVREDOP1   string `json:"tvredop1"`
	OTHREDOP1  string `json:"othredop1"`
	DATEORP2   string `json:"dateorp2"`
	ACBREDOP2  string `json:"acbredop2"`
	AVREDOP2   string `json:"avredop2"`
	MVREDOP2   string `json:"mvredop2"`
	TVREDOP2   string `json:"tvredop2"`
	OTHREDOP2  string `json:"othredop2"`
	DATEORP3   string `json:"dateorp3"`
	ACBREDOP3  string `json:"acbredop3"`
	AVREDOP3   string `json:"avredop3"`
	MVREDOP3   string `json:"mvredop3"`
	TVREDOP3   string `json:"tvredop3"`
	OTHREDOP3  string `json:"othredop3"`
	DATEORP4   string `json:"dateorp4"`
	ACBREDOP4  string `json:"acbredop4"`
	AVREDOP4   string `json:"avredop4"`
	MVREDOP4   string `json:"mvredop4"`
	TVREDOP4   string `json:"tvredop4"`
	OTHREDOP4  string `json:"othredop4"`
	DATEORP5   string `json:"dateorp5"`
	ACBREDOP5  string `json:"acbredop5"`
	AVREDOP5   string `json:"avredop5"`
	MVREDOP5   string `json:"mvredop5"`
	TVREDOP5   string `json:"tvredop5"`
	OTHREDOP5  string `json:"othredop5"`
	DATEORP6   string `json:"dateorp6"`
	ACBREDOP6  string `json:"acbredop6"`
	AVREDOP6   string `json:"avredop6"`
	MVREDOP6   string `json:"mvredop6"`
	TVREDOP6   string `json:"tvredop6"`
	OTHREDOP6  string `json:"othredop6"`
	PRECARD    string `json:"precard"`
	PIDATE     string `json:"pidate"`
	PITHROMB   string `json:"pithromb"`
	PITDATE    string `json:"pitdate"`
	CATH       string `json:"cath"`
	CATHDATE   string `json:"cathdate"`
	ANGINA     string `json:"angina"`
	NYHA       string `json:"nyha"`
	CCS        string `json:"ccs"`
	LVGRADE    string `json:"lvgrade"`
	STRESS     string `json:"stress"`
	DIABETES   string `json:"diabetes"`
	DOI        string `json:"doi"`
	HYPER      string `json:"hyper"`
	CHLSTRL    string `json:"chlstrl"`
	FHX        string `json:"fhx"`
	SMOKE      string `json:"smoke"`
	PACKS      string `json:"packs"`
	COPD       string `json:"copd"`
	COPDS      string `json:"copds"`
	THROMB     string `json:"thromb"`
	PVD        string `json:"pvd"`
	RF         string `json:"rf"`
	NEWRF      string `json:"newrf"`
	DIAL       string `json:"dial"`
	MARFAN     string `json:"marfan"`
	CAROTID    string `json:"carotid"`
	AAD        string `json:"aad"`
	EU         string `json:"eu"`
	RECG       string `json:"recg"`
	CHF        string `json:"chf"`
	SHOCK      string `json:"shock"`
	SYNCOPE    string `json:"syncope"`
	ASP        string `json:"asp"`
	CREAT      string `json:"creat"`
	STATIN     string `json:"statin"`
	AVDIS      string `json:"avdis"`
	MVDIS      string `json:"mvdis"`
	TVDIS      string `json:"tvdis"`
	ENDOCARD   string `json:"endocard"`
	URGENT     string `json:"urgent"`
	AVSURG     string `json:"avsurg"`
	D2         string `json:"d2"`
	AVSIZE     string `json:"avsize"`
	AVPATH     string `json:"avpath"`
	AVPATH2    string `json:"avpath2"`
	AVPATH3    string `json:"avpath3"`
	ANNULEN    string `json:"annulen"`
	AVPROS     string `json:"avpros"`
	MVSURG     string `json:"mvsurg"`
	MVSIZE     string `json:"mvsize"`
	MVPATH     string `json:"mvpath"`
	MVPATH2    string `json:"mvpath2"`
	MVPATH3    string `json:"mvpath3"`
	MVANN      string `json:"mvann"`
	CHORD      string `json:"chord"`
	GORTEX     string `json:"gortex"`
	MVP        string `json:"mvp"`
	MVC        string `json:"mvc"`
	CHORDAL    string `json:"chordal"`
	MVPROS     string `json:"mvpros"`
	TVSURG     string `json:"tvsurg"`
	TVSIZE     string `json:"tvsize"`
	TVPATH     string `json:"tvpath"`
	TVPATH2    string `json:"tvpath2"`
	TVPATH3    string `json:"tvpath3"`
	TVPROS     string `json:"tvpros"`
	PVSURG     string `json:"pvsurg"`
	PVSIZE     string `json:"pvsize"`
	PVPROS     string `json:"pvpros"`
	CI         string `json:"ci"`
	MPAP       string `json:"mpap"`
	SYSAVG     string `json:"sysavg"`
	LVEDP      string `json:"lvedp"`
	PVR        string `json:"pvr"`
	MVGRADR    string `json:"mvgradr"`
	AVAREA     string `json:"avarea"`
	MVAREA     string `json:"mvarea"`
	AVXPLTYPE  string `json:"avxpltype"`
	AVXPLSIZE  string `json:"avxplsize"`
	AVXPLPATH  string `json:"avxplpath"`
	AVXPLDATE  string `json:"avxpldate"`
	MVXPLTYPE  string `json:"mvxpltype"`
	MVXPLSIZE  string `json:"mvxplsize"`
	MVXPLPATH  string `json:"mvxplpath"`
	MVXPLDATE  string `json:"mvxpldate"`
	TVXPLTYPE  string `json:"tvxpltype"`
	TVXPLSIZE  string `json:"tvxplsize"`
	TVXPLPATH  string `json:"tvxplpath"`
	TVXPLDATE  string `json:"tvxpldate"`
	ASSOCOP    string `json:"assocop"`
	LVA        string `json:"lva"`
	SEPT       string `json:"sept"`
	SEPTYPE    string `json:"septype"`
	CHD        string `json:"chd"`
	AAS        string `json:"aas"`
	AOPATH     string `json:"aopath"`
	MAZE       string `json:"maze"`
	OTHERTYPE  string `json:"othertype"`
	ASSDEV     string `json:"assdev"`
	DEVICETYPE string `json:"devicetype"`
	DISLAD     string `json:"dislad"`
	DISCX      string `json:"discx"`
	DISRCA     string `json:"disrca"`
	LMAIN      string `json:"lmain"`
	DISNUM     string `json:"disnum"`
	LIMA       string `json:"lima"`
	RIMA       string `json:"rima"`
	RADIAL     string `json:"radial"`
	SKEGRAFT   string `json:"skegraft"`
	GFTLAD     string `json:"gftlad"`
	GFTCX      string `json:"gftcx"`
	GFTRCA     string `json:"gftrca"`
	ENDART     string `json:"endart"`
	CVPA       string `json:"cvpa"`
	OTHGFT     string `json:"othgft"`
	ACBNUM     string `json:"acbnum"`
	PUMPCASE   string `json:"pumpcase"`
	ORTIME     string `json:"ortime"`
	PUMP       string `json:"pump"`
	CLAMP      string `json:"clamp"`
	CIRARR     string `json:"cirarr"`
	BSA        string `json:"bsa"`
	HT         string `json:"ht"`
	WT         string `json:"wt"`
	MYOPRO     string `json:"myopro"`
	TECH       string `json:"tech"`
	DIRECT     string `json:"direct"`
	HYPOTHER   string `json:"hypother"`
	OFFPUMP    string `json:"offpump"`
	IABP       string `json:"iabp"`
	IECG       string `json:"iecg"`
	CK         string `json:"ck"`
	CKMB       string `json:"ckmb"`
	INO        string `json:"ino"`
	LOS        string `json:"los"`
	RENALINO   string `json:"renalino"`
	POSTRF     string `json:"postrf"`
	OCVENDYS   string `json:"ocvendys"`
	AFIB       string `json:"afib"`
	OCDVT      string `json:"ocdvt"`
	OCPULMC    string `json:"ocpulmc"`
	SEIZURES   string `json:"seizures"`
	PREHB      string `json:"prehb"`
	POSTHB     string `json:"posthb"`
	RBC        string `json:"rbc"`
	NOBED      string `json:"nobed"`
	NONURSE    string `json:"nonurse"`
	ICUCOMP    string `json:"icucomp"`
	CHRPTS     string `json:"chrpts"`
	OTHER      string `json:"other"`
	OTHERNOTE  string `json:"othernote"`
	INFARM     string `json:"infarm"`
	INFLEG     string `json:"infleg"`
	INFSTERN   string `json:"infstern"`
	INFSEP     string `json:"infsep"`
	DCTO       string `json:"dcto"`
	PROC       string `json:"proc"`
	NOTES      string `json:"notes"`
	DRUG4      string `json:"drug4"`
	DRUG5      string `json:"drug5"`
	DRUG6      string `json:"drug6"`
	CORMATRIX  string `json:"cormatrix"`
	CELLSAVER  string `json:"cellsaver"`
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

func main() {

	dirTGH := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TGH perioperative.xlsx"
	// dirTWH := "L:\\CVDMC Students\\Raymond Chou\\perioperative\\TWH perioperative.xlsx"
	tghFile := excelHelper.ConnectToXlsx(dirTGH)
	// tghCols := excelHelper.IdentifyCols(tghFile)
	// twhFile := excelHelper.ConnectToXlsx(dirTWH)
	// twhCols := excelHelper.IdentifyCols(twhFile)
	parseData(tghFile)
}
