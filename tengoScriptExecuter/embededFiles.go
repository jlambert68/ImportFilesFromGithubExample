package tengoScriptExecuter

import (
	_ "embed"
)

// Embed files

//go:embed importTengoLibraries.tengo
var myTengoFile1 []byte

//go:embed SubCustody.Today(n).tengo
var myTengoFile2 []byte

//go:embed "SubCustody.RandomFloatValue(integerSize, decimalSize).tengo"
var myTengoFile3 []byte

var myTengoFile []byte

// Create on file based on all other script files
func concatenateTengoScriptFiles() {

	var tempScript string
	tempScript = string(myTengoFile1)
	tempScript = tempScript + string(myTengoFile2)
	tempScript = tempScript + string(myTengoFile3)

	myTengoFile = []byte(tempScript)
}
