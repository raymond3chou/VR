package event

import "testing"

type Test struct {
	A string
	B int64
}

func TestCutArray(t *testing.T) {
	rowSlice := make(map[string]string)
	rowSlice["A"] = "A"
	rowSlice["B"] = "B"
	rowSlice["C"] = "C"
	reqFields := []string{"B"}
	newRowSlice := CutArray(rowSlice, reqFields)
	for nRS, value := range newRowSlice {
		if nRS != "B" {
			t.Errorf("%s not equal to B", nRS)
		}
		if value != "B" {
			t.Errorf("%s not equal to B", value)
		}
	}
}

// func TestAssignBasic(t *testing.T) {
// 	rowSlice := make(map[string]string)
// 	rowSlice["A"] = "A"
// 	rowSlice["B"] = "2"
// 	rowSlice["C"] = "C"
// 	var event Test
// 	basicStrFields := []string{"A"}
// 	basicIntFields := []string{"B"}
// 	newEvent := AssignBasic(rowSlice, event, basicStrFields, basicIntFields)
// 	newE, err := newEvent.(Test)
// 	if !err {
// 		log.Println("Cant convert")
// 	}
// 	if newE.A != "A" {
// 		t.Errorf("newEvent.A=%s not A", newE.A)
// 	}
// 	if newE.B != 2 {
// 		t.Errorf("newEvent.B=%d not 2", newE.B)
// 	}
// }
