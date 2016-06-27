package main

import "testing"

func TestMatchTable(t *testing.T) {
	tablenames := []string{"contact info", "contactInfo", "contactin19fo", "infoinfo", "inf"}
	match := "info"
	actualphitablenames := matchTable(tablenames, match)
	expectedPHItablenames := []string{"contact info", "infoinfo"}
	nummatched := 0
	for _, expected := range expectedPHItablenames {
		for _, actual := range actualphitablenames {
			if expected == actual {
				nummatched++
			}
		}
	}
	if nummatched != len(expectedPHItablenames) {
		t.Errorf("TestmatchTable failed because %d matched instead of %d", nummatched, len(expectedPHItablenames))
	} else {
		t.Log("MatchTable Works")
	}
}
