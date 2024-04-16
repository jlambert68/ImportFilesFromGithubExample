package luaScriptEngine

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"hash/crc32"
)

// InitiateLuaScriptEngine
// Initiate the Lua Script Engine
func InitiateLuaScriptEngine(luaScriptFiles [][]byte) (err error) {

	// Load Fenix Lua Script files
	var fenixLuaScripts [][]byte
	fenixLuaScripts = loadFenixLuaScripts()

	// Concatenate Fenix Lua scripts with Domain supported Lua scripts
	for _, fenixLuaScript := range fenixLuaScripts {
		luaScriptFiles = append(luaScriptFiles, fenixLuaScript)
	}

	// Save all script into one byte array
	luaScriptFilesAsByteArray = luaScriptFiles

	// Initiate Gopher-Lua Script Engine state variables
	luaState = lua.NewState()

	return err
}

// CloseDownLuaScriptEngine
// Close down the Lua Script Engine in a correct way
func CloseDownLuaScriptEngine() {
	luaState.Close()
}

// LuaScriptEngineExecute
// Execute a specific Lua function
func LuaScriptEngineExecute(inputParameterArray []interface{}, testCaseExecutionUuid string) (responseValue string) {

	var err error
	var luaFunctionToCall string
	var addExtraEntropyValue uint64
	var useEntropyFromTestCaseExecutionUuid bool
	var entropyToUse uint64

	luaFunctionToCall = inputParameterArray[1].(string)
	useEntropyFromTestCaseExecutionUuid = inputParameterArray[4].(bool)
	addExtraEntropyValue = inputParameterArray[5].(uint64)

	// Decide how much entropy to use
	if useEntropyFromTestCaseExecutionUuid == true {

		// Create entropy form TestCaseExecutionUuid by converting string into 32-bit hash
		var entropyBasedOnTestCaseExecutionUuid uint32
		entropyBasedOnTestCaseExecutionUuid = crc32.ChecksumIEEE([]byte(testCaseExecutionUuid))

		// Create final entropy
		entropyToUse = uint64(entropyBasedOnTestCaseExecutionUuid) + addExtraEntropyValue

	}

	// Convert Array to the one used for conversion into Lua Table
	var entropyArray []interface{}
	entropyArray = append(entropyArray, useEntropyFromTestCaseExecutionUuid)
	entropyArray = append(entropyArray, entropyToUse)

	var goArrayToBeConvertedIntoLuaTable []interface{}
	goArrayToBeConvertedIntoLuaTable = inputParameterArray[1:4]
	goArrayToBeConvertedIntoLuaTable = append(goArrayToBeConvertedIntoLuaTable, entropyArray)

	// Create the Lua Input Table
	var luaInputTable *lua.LTable
	luaInputTable = convertToLuaTableRecursively(luaState, inputParameterArray)
	/*
		// Add the 'entropyBasedOnTestCaseExecutionUuid' to the entropy position, which is position 4
		originalValue := luaState.GetField(luaInputTable, "4")
		fmt.Println("Original value at position 4:", originalValue)


		// Add the Go value to the Lua value at position 4
		lv, ok := originalValue.(lua.LNumber)
		if ok == true {
			newValue := lua.LNumber(entropyBasedOnTestCaseExecutionUuid) + lv
			luaState.SetField(luaInputTable, "4", newValue) // Set the new value at position 4
		} else {
			// Original value is not a number
			err = errors.New(fmt.Sprintf("original value, '%s' in position 4 is not a number", originalValue))

			return err.Error()
		}



		// Access and modify the number in the nested table
		nestedTable, ok := originalValue.(*lua.LTable)
		if ok == true {
			nestedValue := luaState.GetField(nestedTable, "1")
			lv, ok := nestedValue.(lua.LNumber)
			if ok == true {
				newValue := lua.LNumber(entropyBasedOnTestCaseExecutionUuid) + lv
				luaState.SetField(nestedTable, "1", newValue) // Update the number in the nested table
			} else {
				// Original value is not a number
				err = errors.New(fmt.Sprintf("original value, '%s' in position 1 in table is not a number", nestedValue))

				return err.Error()
			}
		} else {
			// Original value is not a lua table
			err = errors.New(fmt.Sprintf("original value, '%s' in position 4 is not a table", originalValue))
			return err.Error()
		}
	*/
	// Call function with input parameters and a response
	err = luaState.CallByParam(lua.P{
		Fn:      luaState.GetGlobal(luaFunctionToCall),
		NRet:    1,
		Protect: true,
	},
		luaInputTable)

	if err != nil {
		panic(err)
	}
	returnValue := luaState.Get(-1) // returned value
	luaState.Pop(1)                 // remove received value

	responseValue = returnValue.String()

	return responseValue
}

// 'convertToLuaTableRecursively' recursively converts a Go slice of `[]interface{}` to a Lua table.
func convertToLuaTableRecursively(tempLuaState *lua.LState, goSlice []interface{}) *lua.LTable {
	var luaTable *lua.LTable
	luaTable = tempLuaState.NewTable()
	for _, item := range goSlice {
		switch v := item.(type) {
		case []interface{}:
			// If the item is a slice, recursively convert it to a Lua table
			nestedTable := convertToLuaTableRecursively(tempLuaState, v)
			tempLuaState.RawSet(luaTable, lua.LNumber(luaTable.Len()+1), nestedTable)
		default:
			// Otherwise, add the item directly to the Lua table
			tempLuaState.RawSet(luaTable, lua.LNumber(luaTable.Len()+1), lua.LString(fmt.Sprint(v)))
		}
	}
	return luaTable
}
