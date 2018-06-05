// TypeUtil contains data type related functions.
package common

import (
	"strings"
	"strconv"
	"time"
)

// function to parse a string formatted array back into a real []string
func CleanseArrayedString(val string) []string {
	if len(val)>2 && len(strings.TrimSpace(val))>2 {
		finalVal := val[1:len(val)-1]

		return strings.Split(finalVal, ",")
	}
	return []string{}
}

// function to parse a []string to []int
func ConvertStringArrayToIntArray(stringArray []string) ([]int, error)  {
	if stringArray != nil && len(stringArray)>0 {
		iArray := make([]int, len(stringArray))

		for i, v := range stringArray {
			iVal, err := strconv.Atoi( strings.TrimSpace(v) )
			if err != nil {
				return nil, err
			}
			iArray[i] = iVal
		}
		return iArray, nil
	}
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []int\n", stringArray))
	return nil, nil
}

// function to parse a []string to []float32
func ConvertStringArrayToFloat32Array(stringArray []string) ([]float32, error)  {
	if stringArray != nil && len(stringArray)>0 {
		iArray := make([]float32, len(stringArray))

		for i, v := range stringArray {
			iVal, err := strconv.ParseFloat( strings.TrimSpace(v), 32 )
			if err != nil {
				return nil, err
			}
			iArray[i] = float32(iVal)
		}
		return iArray, nil
	}
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []float32\n", stringArray))
	return nil, nil
}

// function to parse a []string to []float64
func ConvertStringArrayToFloat64Array(stringArray []string) ([]float64, error)  {
	if stringArray != nil && len(stringArray)>0 {
		iArray := make([]float64, len(stringArray))

		for i, v := range stringArray {
			iVal, err := strconv.ParseFloat( strings.TrimSpace(v), 64 )
			if err != nil {
				return nil, err
			}
			iArray[i] = float64(iVal)
		}
		return iArray, nil
	}
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []float64\n", stringArray))
	return nil, nil
}

// function to parse a []string to []bool
func ConvertStringArrayToBoolArray(stringArray []string) ([]bool, error)  {
	if stringArray != nil && len(stringArray)>0 {
		iArray := make([]bool, len(stringArray))

		for i, v := range stringArray {
			iVal, err := strconv.ParseBool( strings.TrimSpace(v) )
			if err != nil {
				return nil, err
			}
			iArray[i] = iVal
		}
		return iArray, nil
	}
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []bool\n", stringArray))
	return nil, nil
}

// function to parse a []string to []time.Time
func ConvertStringArrayToTimeArray(stringArray []string) ([]time.Time, error)  {
	if stringArray != nil && len(stringArray)>0 {
		iArray := make([]time.Time, len(stringArray))

		for i, v := range stringArray {
			iVal, _, err := ParseStringToTimeWithPatterns(
				[]string{TimeDefault, TimeShortDateTime, TimeShortDate},
				strings.TrimSpace(v))

			if err != nil {
				return nil, err
			}
			iArray[i] = iVal
		}
		return iArray, nil
	}
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []time.Time\n", stringArray))
	return nil, nil
}

// trim the contents of the []string
func TrimStringArrayMembers(stringArray []string) []string {
	if stringArray != nil && len(stringArray)>0 {
		for i, v := range stringArray {
			stringArray[i] = strings.TrimSpace(v)
		}
	}
	return stringArray
}


/* ------------------------------------ */
/*	check empty / nil based on type		*/
/* ------------------------------------ */

// function to check if the given string is empty or nil
func IsStringEmptyOrNil(value string) bool {
	if len(strings.TrimSpace(value))==0 {
		return true
	}
	return false
}