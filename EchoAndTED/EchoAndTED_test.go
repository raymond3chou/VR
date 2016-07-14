package main

import (
	"log"
	"testing"

	"github.com/access"
)

func TestWalkDirectory(t *testing.T) {
	errFile := access.CreateErrorLog(true)
	log.SetOutput(errFile)
	defer errFile.Close()
	dir := "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\"
	folderNames := []string{"Example"}

	walkDirectory(folderNames, dir)

	tedFile := access.ReadFile("C:\\Users\\raymond chou\\Desktop\\tedFiles.txt")
	if tedFile != "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test Folder\\tedcodes.xlsx\n" {
		t.Errorf("Printed %s instead of C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test Folder\\tedcodes.xlsx\n", tedFile)
	}

	echoFile := access.ReadFile("C:\\Users\\raymond chou\\Desktop\\echoFiles.txt")
	if echoFile != "C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test Folder\\Echo.xlsx\nC:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test Folder\\Echo2.xlsx\n" {
		t.Errorf("Printed %s instead of C:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test Folder\\Echo.xlsx\nC:\\Users\\raymond chou\\Desktop\\WorkingFiles\\VR\\EchoAndTED\\Example\\Test Folder\\Echo2.xlsx\n", echoFile)
	}

}
