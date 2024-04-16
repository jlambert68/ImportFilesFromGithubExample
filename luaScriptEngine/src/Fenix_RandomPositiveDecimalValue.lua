---------------------------------------------------------------------------------------
-- Module for replacing a Fenix Inception 'Placeholder'
--
-- Version 1.0

-- Placeholder usage
-- 'Fenix.RandomPositiveDecimalValue(IntegerSize, FractionSize)'
-- 'Fenix.RandomPositiveDecimalValue(IntegerSize, FractionSize, IntergerSpace, FractionSpace)'
-- 'Fenix.RandomPositiveDecimalValue[ArraysIndex](IntegerSize, FractionSize)'
-- 'Fenix.RandomPositiveDecimalValue[ArraysIndex](IntegerSize, FractionSize)(UseTestCaseExecutionUuidEntropi)'
-- 'Fenix.RandomPositiveDecimalValue[ArraysIndex](IntegerSize, FractionSize)(UseTestCaseExecutionUuidEntropi, ExtraEntropiNumber)'
-- 'Fenix.RandomPositiveDecimalValue['Integer']('Integer', 'Integer')('Boolean', 'Integer')'
--
-- Resesponse is a lua table;  {value, success, errorMessage} with the following types {'strings', 'boolean', 'string'}
--
-- Usage exmaples
-- 'Fenix.RandomPositiveDecimalValue(2, 3)' is same as 'Fenix.RandomPositiveDecimalValue[1](2, 3)' which is the same as 'Fenix.RandomPositiveDecimalValue[1](2, 3)(true)'
-- They all could produce i.e. '35.693' and with same input the placeholder will allways have the same output.
-- UseTestCaseExecutionUuidEntropi(true/false) is based on the 'TestCaseExecutionUuid' and within one TestCaseExecution values with the same input will have the same output
--
-- 'Fenix.RandomPositiveDecimalValue[2](2, 3)(true)' clould have the output of '27.568'
-- 'Fenix.RandomPositiveDecimalValue[2](2, 3)(false)' will allways produce the same output, independently of 'TestCaseExecutionUuid'
-- 'Fenix.RandomPositiveDecimalValue[2](2, 3)(true, 1)' will add extra entropi to seed, by adding 1 to the value based on 'TestCaseExecutionUuid'.

-- 'IntergerSpace' and 'FractionSpace' define the spaces for 'Intergerpart' and 'FractionPart'. Zeros will be added before 'Intergerpart' and after 'FractionPart'
-- If 'IntergerSpace' is less than 'IntegerSize' then it will be ignored
-- If 'FractionSpace' is less than 'FractionSize' then it will be ignored
-- 'Fenix.RandomPositiveDecimalValue(1, 2, 3, 4)' will produce "004.5700"
-- 'Fenix.RandomPositiveDecimalValue(3, 2, 2, 1)' will produce "344.54"
-- 'Fenix.RandomPositiveDecimalValue(0, 0, 2, 2)' will produce "00.00"
-- 
-- Examples of different parameter values 'IntegerSize' and 'FractionSize' and a possible output
-- 'Fenix.RandomPositiveDecimalValue(1, 2)' = "5.48"
-- 'Fenix.RandomPositiveDecimalValue(2, 2)' = "28.87"
-- 'Fenix.RandomPositiveDecimalValue(0, 1)' = "0.37"
-- 'Fenix.RandomPositiveDecimalValue(0, 0)' = "0"
-- 'Fenix.RandomPositiveDecimalValue(0, 0, 2, 3)' = "00.000"
-- 'Fenix.RandomPositiveDecimalValue(0, 0)' = "0.0"
-- 'Fenix.RandomPositiveDecimalValue(3, 0)' = "293"
-- 'Fenix.RandomPositiveDecimalValue(3, 0, 3, 2)' = "293.00"



local function tableToString(tbl, sep)

    -- More then one parameter in table
    if #tbl > 1 then
        sep = sep or ", "
        local result = "["
        for _, v in ipairs(tbl) do
            result = result .. tostring(v) .. sep
        end

        result = result:sub(1, -#sep - 1)

        result = result .. "]"

        return result
    end

    -- Only one parameter in table
    local result = tostring(tbl)

    return result
end

-- ***********************************************************************************
-- round
--
-- Function to round a decimal to a specific number of decimal places

local function round(x, places)

    local shift = 10 ^ places

    return math.floor(x * shift + 0.5) / shift
end

-- ***********************************************************************************



-- ***********************************************************************************
-- formatDecimal
--
-- Function to format a decimal with a specific number of decimals
local function formatDecimal(number, numberOfDecimals)

    -- Convert the number to a string
    local str = tostring(number)

    -- Find the position of the decimal point
    local dotIndex = nil
    for i = 1, #str do
        if str:sub(i, i) == '.' then
            dotIndex = i
            break
        end
    end


    -- Add a decimal point if it doesn't exist
    if numberOfDecimals > 0 and dotIndex == nil then
        dotIndex = #str + 1
        str = str .. "."
    end

        -- Handle case when there should be decimals 
    if numberOfDecimals > 0 then

        -- Calculate the number of decimal places currently in the string
        local currentDecimals = dotIndex and #str - dotIndex or 0

        -- Add zeros to reach the desired number of decimal places
        while currentDecimals < numberOfDecimals do
            str = str .. "0"
            currentDecimals = currentDecimals + 1
        end

        return str
    end

    -- No decimals so remove decimal point and any following zero, i.e. 13.0
    if dotIndex ~= nil then
    local integerWithOutDecimals = string.sub(str, 1, dotIndex -1)

    return integerWithOutDecimals

    else
        return str

    end

end

-- ***********************************************************************************

-- ***********************************************************************************
-- padValueWithZeros
--
-- Function to pad the integer part and the fraction part with correct number of zeros

local function padValueWithZeros(valueAsString, integerSpace, fractionSpace)

    -- Split the value into integer and fractional parts
    local integerPart, fractionPart = valueAsString:match("^(%d+)%.(%d+)$")
    local noFractions = false
    
    -- Check if value only has intgerpart
    if not integerPart or not fractionPart then
        integerPart = valueAsString
        fractionPart = ""
        noFractions = true
    end

    -- Pad the integer part with zeros if needed
    if #integerPart < integerSpace then
        integerPart = string.rep("0", integerSpace - #integerPart) .. integerPart
    end

    -- Pad the fractional part with zeros if needed
    if #fractionPart < fractionSpace and noFractions == false then
        fractionPart = fractionPart .. string.rep("0", fractionSpace - #fractionPart)
    end

    -- Combine the padded parts if there are two parts
    local zeroPaddedValue = ""

    if noFractions == false then
        zeroPaddedValue = integerPart .. "." .. fractionPart

    else
        zeroPaddedValue = integerPart

    end

    return zeroPaddedValue
end

-- ***********************************************************************************



-- ***********************************************************************************
-- randomize

-- Function to generate random numbers


local function randomize(arrayIndex, maxIntegerPartSize, numberOfDecimals, baseEntropiToUse)

    math.randomseed(baseEntropiToUse+ arrayIndex)

    -- Generate Integer part of random number
    local randomIntegerPart = math.random()
    local integerPart = math.floor(10 ^ maxIntegerPartSize * randomIntegerPart)

    -- Generate Decimal part of random number
    local randomDecimalPart = math.random()
    local decimalPart = 0

    if numberOfDecimals > 0 then
        decimalPart = math.floor(10 ^ numberOfDecimals * randomDecimalPart)
    end

    -- Combine Integer and decimal part into one random number
    local randomNumber = integerPart + 10 ^ (-1 * numberOfDecimals) * decimalPart

    randomNumber = round(randomNumber, numberOfDecimals)

    return randomNumber
end


-- ***********************************************************************************




 -- ***********************************************************************************
-- Fenix_RandomDecimalValue_ArrayValue // Fenix.RandomDecimalValue[n](maxIntegerPartSize, numberOfDecimals)
--
-- Function to generate random value with a specif max number of integer and speciic number of decimals
-- inputArray := [arrayPosition, maxIntegerPartSize, numberOfDecimals, testCaseUuidEntropi]


local function Fenix_RandomDecimalValue_ArrayValue(inputArray)
    local arrayPositionToUse = inputArray[1]
    local maxIntegerPartSize = inputArray[2][1]
    local numberOfDecimals = inputArray[2][2]

    local entropiToUse = inputArray[3]

    local tempValueAsDecimal = randomize(arrayPositionToUse, maxIntegerPartSize, numberOfDecimals, entropiToUse)

    local valueIsBaseFormated =  formatDecimal(tempValueAsDecimal, numberOfDecimals)

    -- No padding should be done
    if #inputArray[2] == 2 then
        return valueIsBaseFormated
    end

    -- extract padding sizes
    local integerSpace = inputArray[2][3]
    local fractionSpace = inputArray[2][4]

    local zeroPaddedValue = padValueWithZeros(valueIsBaseFormated, integerSpace, fractionSpace)

    return zeroPaddedValue

end

-- ***********************************************************************************



-- ***********************************************************************************
-- Fenix_RandomPositiveDecimalValue
--
-- Function to generate random value with a specif max number of integer and speciic number of decimals
-- Always use array value 1, first array position from user perspective
--
-- inputArray := [[arrayindex], [maxIntegerPartSize, numberOfDecimals], testCaseUuidEntropi]

function Fenix_RandomPositiveDecimalValue(inputTable)

    local responseTable = {
        success = true,
        value = "",
        errorMessage = ""
    }

    -- There must be 4 rows in the InputTable
    if #inputTable ~= 4 then

        local error_message = "Error - there should be exactly four rows in InputTable."

        responseTable.success = false
        responseTable.errorMessage = error_message

        return responseTable

    end


    -- Extract ArraysIndexArray
    local arraysIndexTable = inputTable[2]
    local arraysIndexToUse

    -- Secure that ArraysIndexArray is not emtpty or only have one value
    if #arraysIndexTable > 1 then
        -- Have more then 1 value

        -- Convert array to string
        local tableAsString = tableToString (arraysIndexTable, ",")
        local error_message = "Error - array index array can only have a maximum of one value. '" .. tableAsString .. "'"

        responseTable.success = false
        responseTable.errorMessage = error_message

        return responseTable

    elseif  #arraysIndexTable == 0 then
        -- zero array index, so use first index position

        arraysIndexToUse = 1

    else
        -- Array has one value so use that
        arraysIndexToUse = arraysIndexTable[1]

    end


    -- Extract FunctionArgumentsArray
    local functionArgumentsTable = inputTable[3]

    -- Handle if function arguments is not 2 arguments or 4 arguments
    if #functionArgumentsTable ~=  2 and #functionArgumentsTable ~=  4 then

        -- More then 2 arguments
        if #functionArgumentsTable > 2 then

            local tableAsString = tableToString(functionArgumentsTable, ",")
            local error_message = "Error - there must be exact 2 or 4 function parameter. '" .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message

            return responseTable

            -- Exact one argument
        elseif #functionArgumentsTable == 1 then

                local result = "[" .. tostring(functionArgumentsTable[1]) .. "]"

                local error_message = "Error - there must be exact 2 or 4 function parameter. '" .. result .. "'"

                responseTable.success = false
                responseTable.errorMessage = error_message

                return responseTable

            -- Zero values    
        else
                local error_message = "Error - there must be exact 2 or 4 function parameter but it is empty."

                responseTable.success = false
                responseTable.errorMessage = error_message

                return responseTable

        end
    end

    -- verify that each function parameter is a number
    for _, v in ipairs(functionArgumentsTable) do

        -- Must be an integer 
        if type(v) ~=  "number" then
            local tableAsString = tableToString (functionArgumentsTable, ",")
            local error_message = "Error - functions parameters must be of type Integer. '" .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message


        end
    end

    -- Extract entropi
    local entropiTable = inputTable[4]

    -- Verify that content in Entropi is of type 'Table'
    if type(entropiTable) ~= "table" then

        local error_message = "Error - Entropi is not of type 'Table', but is of type '" .. type(entropiTable)  .. "'."

        responseTable.success = false
        responseTable.errorMessage = error_message

        return responseTable
    end

    -- verify that each paramter in Entropi table is a number and sum up all entry values into one
    local entropiValue = 0
    
    for _, v in ipairs(entropiTable) do

        -- Must be an integer 
        if type(v) ~=  "number" then
            local tableAsString = tableToString (entropiTable, ",")
            local error_message = "Error - entropi parameters must be of type 'Integer', expect for first parameter which is a 'Boolean' '" .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message


        else
            entropiValue = entropiValue + v
        end

    end





    -- Make new Array to be send to the function that does stuff
    local inputTableForProcessingen = {arraysIndexToUse, functionArgumentsTable, entropiValue}

    -- Call and process Random Decimal Value
    local respons = Fenix_RandomDecimalValue_ArrayValue(inputTableForProcessingen)

    responseTable.success = true
    responseTable.errorMessage = ""
    responseTable.value = respons

    return responseTable

end




local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{2, 3}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{2, 3}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")



local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{2, 3}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{2, 3}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{1, 2}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{1, 2}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '8.98'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {2},{1, 2}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {2},{1, 2}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '6.48'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{1, 1}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{1, 1}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '8.9'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{1, 1}, 1}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{1, 1}, 1}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '6.4'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{1, 1}, 1}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{1, 1}, 1}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '6.4'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{0, 1}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{0, 1}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0.9'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{1, 0}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{1, {0}}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '8'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{0, 0}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{0, {0}}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{0, 0, 2, 3}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{0, 0, 2, 3}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{0, 2, 3, 4}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{0, 2, 3, 4}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{2, 2, 3, 4}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{2, 2, 3, 4}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{6, 6}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{6, 6}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '815587.986577'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{6, 10}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{6, 10}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.value .. " :: Expected OK - i.e. '815587.9865775100'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{0}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{0}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected ERROR - there must be exact 2 function parameter. '[0]'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected Error - there must be exact 2 function parameter but it is empty.")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{1, 2, 3}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{1, 2, 3}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected Error - there must be exact 2 function parameter. '[1,2,3]'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1, 2},{2, 3}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1, 2},{2, 3}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected Error - array index array can only have a maximum of one value. '[1,2]'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{2, 3}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{2, 3}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected Error - there should be exactly four rows in InputTable.")
print("")





-- ***********************************************************************************
--]]