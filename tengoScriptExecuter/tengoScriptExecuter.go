package tengoScriptExecuter

import (
	"fmt"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"log"
)

func ExecuteScripte(tengoFunctionName string, inputParameters []string) (responseValue string) {

	concatenateTengoScriptFiles()

	script := tengo.NewScript(myTengoFile2)

	// Import time module from stdlib
	script.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))

	script.Add("shift_days_f883cffd80", 4)

	// Compile the script
	compiled, err := script.Compile()
	if err != nil {
		log.Fatalln("Error compiling script:", err)
	}

	compiled.Set("shift_days_f883cffd80", 5) // Set myVar to 5

	compiled, err = script.Compile()
	if err != nil {
		log.Fatalln("Error compiling script:", err)
	}

	err = compiled.Run()
	if err != nil {
		log.Fatalln("Error compiling script:", err)
	}

	// Get the function from the script
	functionVariable2 := compiled.Get("SubCustody_TodayShiftDay_InputParameters") // This returns *tengo.Variable

	// Check if the function variable is nil
	if functionVariable2 == nil {
		log.Fatalf("Function %s does not exist in the script", functionVariable2)
	}

	ouut := compiled.Get("SubCustody_TodayShiftDay_InputParameters_out")
	fmt.Println(ouut)

	ouut2 := compiled.Get("SubCustody_TodayShiftDay_out")
	fmt.Println(ouut2)

	err = compiled.Set("shift_days_f883cffd80", 1)
	if err != nil {
		log.Fatalln("Error compiling script:", err)
	}

	err = compiled.Run()
	if err != nil {
		log.Fatalln("Error compiling script:", err)
	}

	ouut3 := compiled.Get("SubCustody_TodayShiftDay_out")
	fmt.Println(ouut3)
	/*
		// Get the function from the script for input parameters
		parameterFuntion := "SubCustody_TodayShiftDay_InputParameters" //tengoFunctionName + "_InputParameters"

		functionVariable := compiled.Get(parameterFuntion) // This returns *tengo.Variable

		// Assert the underlying value to be *tengo.UserFunction
		tengoFunctionToCall, ok := functionVariable.Value().(*tengo.UserFunction)
		if !ok {
			//log.Fatalln("Function " + parameterFuntion + " not found or not a user function")
		}

		fmt.Println(tengoFunctionToCall.CanCall())

		// Call the function to get input parameters
		result, err := tengoFunctionToCall.Call()
		if err != nil {
			log.Fatalln("Error calling function:", err)
		}

		responseValue = result.String()

		fmt.Println("result:", result)
	*/
	return responseValue
}
