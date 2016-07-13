package main

import "testing"

func TestWalkDirectory(t *testing.T) {
	dir := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\"
	folderNames := []string{"Example"}
	walkDirectory(folderNames, dir)
	// if EchoPaths[0] != "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test.xlsx" {
	// 	t.Errorf("The Path was %s and not %s", EchoPaths[0], "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test.xlsx")
	// }
	if EchoPaths[0] != "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test Folder\\Echo.xlsx" {
		t.Errorf("The Path was %s and not %s", EchoPaths[0], "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test Folder\\Echo.xlsx")
	}
}
