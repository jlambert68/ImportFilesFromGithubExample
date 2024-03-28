package tengoScriptExecuter

import (
	_ "embed"
)

// Embed files

//go:embed importTengoLibraries.tengo
var myTengoFile1 []byte

//go:embed SubCustody.Today(n).tengo
var myTengoFile2 []byte

//go:embed "SubCustody_RandomPositiveFloatValue(integerSize, decimalSize).tengo"
var myTengoFile3 []byte

var myTengoFile [][]byte

// Create on file based on all other script files
func concatenateTengoScriptFiles() {

	myTengoFile = append(myTengoFile, myTengoFile2)
	myTengoFile = append(myTengoFile, myTengoFile3)
}
