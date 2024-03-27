package tengoScriptExecuter

import (
	"fmt"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"log"
	"strings"
)

func ExecuteScripte(tengoFunctionName string, inputParameterArray []interface{}) (responseValue string) {

	// Add TestCaseUUid randomness
	inputParameterArray = append(inputParameterArray, 0)

	concatenateTengoScriptFiles()

	tengoFunctionName = strings.ReplaceAll(tengoFunctionName, ".", "_")

	var script *tengo.Script

	switch tengoFunctionName {

	case "SubCustody_TodayShiftDay":
		script = tengo.NewScript(myTengoFile2)
	case "SubCustody_RandomFloatValue":
		script = tengo.NewScript(myTengoFile3)
	case "SubCustody_RandomFloatValue_ArrayValue":
		script = tengo.NewScript(myTengoFile3)
	default:
		responseValue = fmt.Sprintf("Unknown 'tengoFunctionName¨ - '%s", tengoFunctionName)

		return responseValue

	}

	// Import time module from stdlib
	script.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))

	script.Add("inputArray", inputParameterArray)

	// Compile the script
	compiled, err := script.Compile()
	if err != nil {
		log.Fatalln("Error compiling script:", err)
	}

	err = compiled.Run()
	if err != nil {
		log.Fatalln("Error compiling script:", err)
	}

	responseVariableInTengoScript := tengoFunctionName + "_out"
	functionResponse := compiled.Get(responseVariableInTengoScript)

	// Check if the function variable is nil
	if functionResponse == nil {
		responseValue = fmt.Sprintf("Function %s does not exist in the script", tengoFunctionName)

		return responseValue
	}

	responseValue = functionResponse.String()

	return responseValue
}
