# Starlark script

# Function to shift the current date by 'n' days and return the new date
# In Starlark, assuming `shift_date` is a provided function to shift dates
def subcustody_today_shift_day(shift_days):
    today = shift_date(0)  # Assuming this gets today's date; implement accordingly
    shifted_date = shift_date(shift_days)
    return shifted_date.strftime("%Y-%m-%d")  # Format the date

# Entry point function
def tengo_script_starting_point(input_array):
    if len(input_array) != 4:
        return "Error - there should be exactly four parameters in InputArray."

    function_name = input_array[0]
    if function_name == "SubCustody_TodayShiftDay":
        shift_days = input_array[2][0] if input_array[2] else 0
        return subcustody_today_shift_day(shift_days)
    else:
        return "ERROR - Unknown function '{}'".format(function_name)

# Example usage
input_array = ["SubCustody_TodayShiftDay", [], [1], 0]  # Example input for "tomorrow"
response = tengo_script_starting_point(input_array)
print(response)
