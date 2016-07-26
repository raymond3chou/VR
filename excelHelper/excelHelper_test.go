package excelHelper

import "testing"

func TestWriteStruct(t *testing.T) {
	colNameSlice := []string{"DAYSPOST", "ICUNUM", "ICU", "VENT", "CREAT", "AVSIZE", "MVSIZE", "TVSIZE", "PVSIZE", "CI", "MPAP", "SYSAVG", "LVEDP", "PVR", "MVGRADR", "AVAREA", "MVAREA", "AVXPLSIZE", "MVXPLSIZE", "TVXPLSIZE", "DISNUM", "GFTLAD", "GFTCX", "GFTRCA", "ACBNUM", "ORTIME", "PUMP", "CLAMP", "CIRARR", "HT", "WT", "REOPNUM", "CK", "CKMB", "PREHB", "POSTHB", "PACKCELLS"}
	WriteStruct(colNameSlice)
}
