package common

import (
	"strings"
	"strconv"
	"time"
)

/**
 *	helper method to split the given string to []string
 */
func CleanseArrayedString(val string) []string {
	if len(val)>2 && len(strings.TrimSpace(val))>2 {
		finalVal := val[1:len(val)-1]

		return strings.Split(finalVal, ",")
	}
	return []string{}
}

/**
 *	handy method to convert the given []string into []int
 */
func ConvertStringArrayToIntArray(sArray []string) ([]int, error)  {
	if sArray != nil && len(sArray)>0 {
		iArray := make([]int, len(sArray))

		for i, v := range sArray {
			iVal, err := strconv.Atoi( strings.TrimSpace(v) )
			if err != nil {
				return nil, err
			}
			iArray[i] = iVal
		}
		return iArray, nil
	}
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []int\n", sArray))
	return nil, nil
}

/**
 *	handy method to convert the given []string into []float32
 */
func ConvertStringArrayToFloat32Array(sArray []string) ([]float32, error)  {
	if sArray != nil && len(sArray)>0 {
		iArray := make([]float32, len(sArray))

		for i, v := range sArray {
			iVal, err := strconv.ParseFloat( strings.TrimSpace(v), 32 )
			if err != nil {
				return nil, err
			}
			iArray[i] = float32(iVal)
		}
		return iArray, nil
	}
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []float32\n", sArray))
	return nil, nil
}
/**
 *	handy method to convert the given []string into []float32
 */
func ConvertStringArrayToFloat64Array(sArray []string) ([]float64, error)  {
	if sArray != nil && len(sArray)>0 {
		iArray := make([]float64, len(sArray))

		for i, v := range sArray {
			iVal, err := strconv.ParseFloat( strings.TrimSpace(v), 64 )
			if err != nil {
				return nil, err
			}
			iArray[i] = float64(iVal)
		}
		return iArray, nil
	}
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []float64\n", sArray))
	return nil, nil
}

/**
 *	handy method to convert the given []string into []bool
 */
func ConvertStringArrayToBoolArray(sArray []string) ([]bool, error)  {
	if sArray != nil && len(sArray)>0 {
		iArray := make([]bool, len(sArray))

		for i, v := range sArray {
			iVal, err := strconv.ParseBool( strings.TrimSpace(v) )
			if err != nil {
				return nil, err
			}
			iArray[i] = iVal
		}
		return iArray, nil
	}
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []bool\n", sArray))
	return nil, nil
}

/**
 *	handy method to convert the given []string into []time.Time
 */
func ConvertStringArrayToTimeArray(sArray []string) ([]time.Time, error)  {
	if sArray != nil && len(sArray)>0 {
		iArray := make([]time.Time, len(sArray))

		for i, v := range sArray {
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
	//return nil, errors.New(fmt.Sprintf("Failed to convert %v to []time.Time\n", sArray))
	return nil, nil
}

/**
 *	trim the given []string
 */
func TrimStringArrayMembers(sArray []string) []string {
	if sArray != nil && len(sArray)>0 {
		for i, v := range sArray {
			sArray[i] = strings.TrimSpace(v)
		}
	}
	return sArray
}


/* ------------------------------------ */
/*	check empty / nil based on type		*/
/* ------------------------------------ */

func IsStringEmptyOrNil(value string) bool {
	if len(strings.TrimSpace(value))==0 {
		return true
	}
	return false
}