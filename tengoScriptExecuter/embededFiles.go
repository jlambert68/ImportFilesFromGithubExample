package tengoScriptExecuter

import (
	_ "embed"
)

// Embed the file
//
//go:embed SubCustody.Today(n).tengo
var myTengoFile []byte
