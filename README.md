# VR
##Programs
### Echo and TED
This program searches through xlsx files and checks whether they are an ECHO file or a TEDCODE file. If either is true then the program saves a copy of the xlsx file to a specified directory.

###Peri
This program reads an xlsx file, parses each row, and checks if they are valid cells. If so, the program iterates through a premade struct and assigns each cell to their corresponding struct field.

###VR GoSQL
This program connects to MS Access, goes through each sheet/table in the database and extract all the needed information then congregates them into a single sheet/table. 

###AR
This program parses data from MS Access and combines them into another MS Access DB

###PHIJsonConvertor
This program reads in data from a MS Access DB and parses them into JSON objects using a set of rules.

###Peri
This program reads data from MS Excel and parses the date into a JSON object for perioperative data.

###PeriOpEvents
This program reads data from MS Excel and for a specific set of columns a JSON object is created.

##Packages
###accessHelper
This package contains helper functions for MS Access such as type conversions and file I/O.

###excelHelper 
This package contains helper functions for MS Excel such as reading columns/rows and comparisons.

###periopchecks
This package constains check functions for Perioperative data code as code validations.
