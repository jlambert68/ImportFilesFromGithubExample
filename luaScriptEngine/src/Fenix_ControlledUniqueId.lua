function ProcessControlledUniqueId(input, inputTable)

    if type(input) ~= "string" then
        error("Input must be a string, got " .. type(input))
    end

    local function randomString(length, upper, seedValue)
        if seedValue ~= nil then
            math.randomseed(seedValue)
        end

        local chars = upper and 'ABCDEFGHIJKLMNOPQRSTUVWXYZ' or 'abcdefghijklmnopqrstuvwxyz'
        local result = ''
        for i = 1, length do
            local randIndex = math.random(#chars)
            result = result .. chars:sub(randIndex, randIndex)
        end
        return result
    end

    local function replaceDateTime(pattern, format)
        input = input:gsub(pattern, os.date(format))
    end

    local function replaceRandomNumber(pattern, seedValue)
        
        math.randomseed(seedValue)

        input = input:gsub(pattern, function(n) 
            return string.format("%0"..#n.."d", math.random(10^#n-1))
        end)
    end

    local function replaceRandomString(pattern, upper, seedValue)

        math.randomseed(seedValue)

        input = input:gsub(pattern, function(n, seed) 
            return randomString(tonumber(n), upper, tonumber(seed))
        end)
    end

    -- Extract from imput table
    local arrayPositionTable  = inputTable[1]
    local seedValueTable = inputTable[2]

    -- Replace date patterns
    replaceDateTime("%%YYYY%-MM%-DD%%", "%Y-%m-%d")
    replaceDateTime("%%YYYYMMDD%%", "%Y%m%d")
    replaceDateTime("%%YYMMDD%%", "%y%m%d")

    -- Replace time patterns
    replaceDateTime("%%hh:mm:ss%%", "%H:%M:%S")
    replaceDateTime("%%hh%.mm%.ss%%", "%H.%M.%S")
    replaceDateTime("%%hhmmss%%", "%H%M%S")
    replaceDateTime("%%hhmm%%", "%H%M")


    -- Replace time with milliseconds
    input = input:gsub("%%hh:mm:ss%%", function()
        return os.date("%H:%M:%S")
    end)
    input = input:gsub("%%hh.mm.ss%%", function()
        return os.date("%H.%M.%S")
    end)

    -- Create seed value
    local arrayPosition  = arrayPositionTable[1]
    local seedValue = 0
    for _, value in ipairs(seedValueTable) do
        seedValue = seedValue + value
    end

    -- Add Array posistion
    seedValue = seedValue + arrayPosition


    -- Replace random number pattern
    replaceRandomNumber("%%(n+)%%", seedValue)

    -- Replace random string patterns with seeding
    replaceRandomString("%%a%((%d+),%s*(%d+)%)%%", false, seedValue)
    replaceRandomString("%%A%((%d+),%s*(%d+)%)%%", true, seedValue)

    return input
end


function Fenix_ControlledUniqueId(inputTable)
    
    -- ExtractInput
    local arrayPositionTable  = inputTable[2]
    local textToProcess = inputTable[3][1]
    local seedValueTable = inputTable[4]

    -- Secure that 'textToProcess' is of string type 
    if type(textToProcess) ~= "string" then
        error("textToProcess must be a string, got " .. type(textToProcess))
    end
    
    -- Create a new Input array
    local newInputTable = {arrayPositionTable, seedValueTable}

    local result = ProcessControlledUniqueId(textToProcess, newInputTable)

    return result

end
-- Example usage
local inputString = "Date: %YYYY-MM-DD%, Date: %YYYYMMDD%, Date: %YYMMDD%, Time: %hh:mm:ss%, Time: %hhmmss%, Time: %hhmm%, Random Number: %nnnnn%, Random String: %a(5, 11)%, Random String Uppercase: %A(5, 10)%, Time: %hh:mm:ss%, Time: %hh.mm.ss% "
local inputTable = {"ConstrolledUniqueId", {0}, {inputString}, {0}}
local result = Fenix_ControlledUniqueId(inputTable)
print(inputString)
print(result)
