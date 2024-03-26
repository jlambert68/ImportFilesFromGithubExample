package tengoScriptExecuter

import (
	"fmt"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

func ExecuteScripte(dateShift string) (responseValue string) {
	// Initialize a new Tengo script
	/*
			script := tengo.NewScript([]byte(`
				// ShiftDays function shifts the current date by 'n' days and returns the new date in YYYY-MM-DD format
				fmt := import("fmt")
				times := import("times")

				shift_days := func(n) {
		    		today := times.now()               // Get the current date and time
		   			shiftedDate := times.add_date(today, 0, 0, n) // Shift the date by 'n' daysadd_date(t time, years int, months int, days int)
		    		return times.time_format(shiftedDate,"2006-01-02")  // Return the formatted date
				}

				dateShifted := shift_days(input)

				fmt.println("Shifted Date:", dateShifted)
			`))

	*/

	script := tengo.NewScript(myTengoFile)

	// Set 'input' variable in the script
	inputDays := dateShift // This can be positive, negative, or zero
	_ = script.Add("input", inputDays)

	// Import time module from stdlib
	script.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))

	// Compile and run the script
	if compiled, err := script.Run(); err != nil {
		fmt.Println("Error:", err)
		responseValue = err.Error()
	} else {
		// Retrieve the shifted date from the script
		dateShifted := compiled.Get("dateShifted")

		responseValue = dateShifted.String()

		fmt.Println("Shifted Date:", dateShifted)
	}

	return responseValue
}
