package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/access"
)

//Info blah
type Info struct {
	Name  int     `json:"Name"`
	Value []Value `json:"Value_"`
}

//Value blah
type Value struct {
	V string `json:"field"`
}

func writeJSON() {
	path := "C:\\Users\\raymond chou\\Desktop\\PeriOp\\test.txt"
	accessHelper.CreateFile(path)
	jsonFile, _ := accessHelper.ConnectToTxt(path)
	v := Value{"Jim"}
	v2 := Value{"Time:3"}
	var vArray []Value
	vArray = append(vArray, v)
	vArray = append(vArray, v2)

	i := Info{2, vArray}

	j, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
	}
	jsonFile.Write(j)

}

func printStructFields() {
	i := reAttributes()
	fmt.Printf("Name: %d\n", i.Name)

}

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

func reAttributes() Info {
	var i Info
	for n := 0; n < reflect.TypeOf(i).NumField(); n++ {
		field := reflect.TypeOf(i).Field(n).Name
		if strings.Contains(field, "Name") {
			reflect.ValueOf(&i).Elem().FieldByName(field).SetInt(9)
		}
	}
	return i
}
