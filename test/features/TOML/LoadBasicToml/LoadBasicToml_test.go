/*
 *  Copyright Project - CFactor, Author - quoeamaster, (C) 2018
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package LoadBasicToml

import (
	"github.com/DATA-DOG/godog"
	"fmt"
	"CFactor/TOML"
	"reflect"
	"strconv"
	"strings"
	"CFactor/common"
	"time"
	TOML2 "CFactor/test/features/TOML"
)

// class level variable
var configReader TOML.TOMLConfigImpl
var config TOML2.DemoTOMLConfig

func foundATomlFileLocation(name string) error {
	// somehow you need to know the target Config object/struct's type
	if len(name)>0 {
		configReader = TOML.NewTOMLConfigImpl(name, reflect.TypeOf(TOML2.DemoTOMLConfig{}))
		return nil

	} else {
		return fmt.Errorf("the given 'name' is not Valid (%v)", name)
	}
}

func loadToml(name string) error {
	// create an instance for population
	configObject := TOML2.DemoTOMLConfig{ Author: TOML2.Author{} }

	// no overriding parameters supplied
	_, err := configReader.Load(&configObject)
	if err != nil {
		return fmt.Errorf("Error in loading the TOML file. %v\n", err)
	}
	config = configObject
	fmt.Println(config.String())

	return nil
}

func iShouldBeAbleToAccessTheFieldsFromThisTomlFile() error {
	// really just to add this "feature" line for clarity, no actions are required
	return nil
}

func checkFieldValue(field, value string) error {
	switch field {
	case "version":
		if strings.Compare(config.Version, value) != 0 {
			return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, config.Version)
		}
	case "author.firstName":
		if strings.Compare(config.Author.FirstName, value) != 0 {
			return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, config.Author.FirstName)
		}
	case "role":
		if strings.Compare(config.Role, value) != 0 {
			return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, config.Role)
		}

	default:
		return fmt.Errorf("unsupported field [%v]", field)
	}

	ok, val := configReader.GetStringValueByKey(config, field)
	if !ok {
		return fmt.Errorf("given %v's value not FOUND", field)
	}
	if strings.Compare(val, value)==0 {
		// additional check
		if !configReader.IsFieldStringValueMatched(config, field, value) {
			return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, val)
		}
		return nil
	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, val)
}

func theIntegerValueForFieldIs(field string, value int) error {
	ok, val := configReader.GetIntValueByKey(config, field)
	if !ok {
		return fmt.Errorf("given %v's value not FOUND", field)
	}
	if val==int64(value) {
		return nil
	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, val)
}

func theFloatValueForFieldIs(field string, value float32) error {
	/*
	 ** this is the normal way to do checking; however there are cases that
	 **	you need reflection api to check dynamic struct values
	 */
	if strings.Compare("author.height", field)==0 {
		if value == config.Author.Height {
			return nil
		}
		return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, config.Author.Height)
	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, config.Author.Height)
}

func theBoolValueForFieldIs(field, value string) error {
	bValue, _ := strconv.ParseBool(value)
	ok, val := configReader.GetBoolValueByKey(config, field)
	if !ok {
		return fmt.Errorf("given %v's value not FOUND", field)
	}
	if val==bValue {
		return nil
	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, val)
}

func theTimeValueForFieldIs(field, valueInString string) error {
	// parse the valueInString to time.Time
	// if you know the pattern ... use common.ParseStringToTime; else ...
	patterns := []string { common.TimeDefault, common.TimeShortDateTime, common.TimeShortDate}
	t0, _, err := common.ParseStringToTimeWithPatterns(patterns, valueInString)
	if err != nil {
		return fmt.Errorf("the given time (string format) is not valid {%v}", err)
	}
	//fmt.Println("[debug] matched format => ", format)

	// equality check
	ok, val := configReader.GetTimeValueByKey(config, field)
	if !ok {
		return fmt.Errorf("given %v's value not FOUND", field)
	}
	if t0.Equal(val) {
		return nil
	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, t0, val)
}

func theStrArrValueForFieldAtIndexIsCapIs(field string, arrayIdx int, value string, arraySize int) error {
	// semi-hard code test case (for simplicity)
	var actualArrSize int
	var actualVal string

	if strings.Compare("hobbies", field)==0 {
		sArr := config.Hobbies
		actualVal = sArr[arrayIdx]
		actualArrSize = len(sArr)

		if strings.Compare(actualVal, value)==0 && arraySize == actualArrSize {
			return nil
		}

	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v) / size might also not match {%v} vs [%v]", field, value, actualVal, arraySize, actualArrSize)
}
func theIntArrayValueForFieldAtIndexIsCapIs(field string, arrayIdx int, value int, arraySize int) error {
	// semi-hard code test case (for simplicity)
	var actualVal, actualArrSize int

	if strings.Compare("author.luckyNumbers", field)==0 {
		sArr := config.Author.LuckyNumbers
		actualVal = sArr[arrayIdx]
		actualArrSize = len(sArr)

		if actualVal == value && arraySize == actualArrSize {
			return nil
		}

	} else if strings.Compare("taskNumbers", field)==0 {
		sArr := config.TaskNumbers
		actualVal = sArr[arrayIdx]
		actualArrSize = len(sArr)

		if actualVal == value && arraySize == actualArrSize {
			return nil
		}

	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v) / size might also not match {%v} vs [%v]", field, value, actualVal, arraySize, actualArrSize)
}

func theFloat32ArValueForFieldAtIndexIsCapIs(field string, arrayIdx int, value float32, arraySize int) error {
	// semi-hard code test case (for simplicity)
	var actualArrSize int
	var actualVal float32

	if strings.Compare("floatingPoints32", field)==0 {
		sArr := config.FloatingPoints32
		actualVal = sArr[arrayIdx]
		actualArrSize = len(sArr)

		if actualVal == value && arraySize == actualArrSize {
			return nil
		}

	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v) / size might also not match {%v} vs [%v]", field, value, actualVal, arraySize, actualArrSize)
}
func theFloat64ArValueForFieldAtIndexIsCapIs(field string, arrayIdx int, value float64, arraySize int) error {
	// semi-hard code test case (for simplicity)
	var actualArrSize int
	var actualVal float64

	if strings.Compare("author.attributes64", field)==0 {
		sArr := config.Author.Attributes64
		actualVal = sArr[arrayIdx]
		actualArrSize = len(sArr)

		if actualVal == value && arraySize == actualArrSize {
			return nil
		}

	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v) / size might also not match {%v} vs [%v]", field, value, actualVal, arraySize, actualArrSize)
}

func theBoolArrayValueForFieldAtIndexIsCapIs(field string, arrayIdx int, value string, arraySize int) error {
	// semi-hard code test case (for simplicity)
	var actualArrSize int
	var actualVal bool

	if strings.Compare("author.likes", field)==0 {
		sArr := config.Author.Likes
		actualVal = sArr[arrayIdx]
		actualArrSize = len(sArr)

		// parse value to bool
		bVal, cErr := strconv.ParseBool(value)
		if cErr == nil && actualVal == bVal && arraySize == actualArrSize {
			return nil
		}

	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v) / size might also not match {%v} vs [%v]", field, value, actualVal, arraySize, actualArrSize)
}

func theTimeArrayValueForFieldAtIndexIsCapIs(field string, arrayIdx int, value string, arraySize int) error {
	// semi-hard code test case (for simplicity)
	var actualArrSize int
	var actualVal time.Time

	if strings.Compare("author.registrationDates", field)==0 {
		sArr := config.Author.RegistrationDates
		actualVal = sArr[arrayIdx]
		actualArrSize = len(sArr)

		// parse value to bool
		tVal, _, cErr := common.ParseStringToTimeWithPatterns(
			[]string{
				common.TimeDefault,
				common.TimeShortDate,
				common.TimeShortDateTime},
			value)
		if cErr == nil && actualVal.Equal(tVal) && arraySize == actualArrSize {
			return nil
		}

	} else if strings.Compare("specialDates", field)==0 {
		sArr := config.SpecialDates
		actualVal = sArr[arrayIdx]
		actualArrSize = len(sArr)

		// parse value to bool
		tVal, _, cErr := common.ParseStringToTimeWithPatterns(
			[]string{
				common.TimeDefault,
				common.TimeShortDate,
				common.TimeShortDateTime},
			value)
		if cErr == nil && actualVal.Equal(tVal) && arraySize == actualArrSize {
			return nil
		}
	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v) / size might also not match {%v} vs [%v]", field, value, actualVal, arraySize, actualArrSize)
}

// testing the features of this BDD story use case
func FeatureContext(s *godog.Suite) {
	s.Step(`^there is a TOML in the current folder named "([^"]*)"$`, foundATomlFileLocation)
	s.Step(`^I load the TOML file named "([^"]*)"$`, loadToml)
	s.Step(`^I should be able to access the fields from this toml file$`, iShouldBeAbleToAccessTheFieldsFromThisTomlFile)
	s.Step(`^the value for field "([^"]*)" is "([^"]*)"$`, checkFieldValue)
	s.Step(`^the integer value for field "([^"]*)" is (\d+)$`, theIntegerValueForFieldIs)
	s.Step(`^the float value for field "([^"]*)" is (\d+\.\d+)$`, theFloatValueForFieldIs)
	s.Step(`^the bool value for field "([^"]*)" is "([^"]*)"$`, theBoolValueForFieldIs)
	s.Step(`^the time value for field "([^"]*)" is "([^"]*)"$`, theTimeValueForFieldIs)
	s.Step(`^the array value for field "([^"]*)" at index "(\d+)" is "(\d+)" cap is "(\d+)"$`, theIntArrayValueForFieldAtIndexIsCapIs)
	s.Step(`^the array value for field "([^"]*)" at index "(\d+)" is "([^"]*)" cap is "(\d+)"$`, theStrArrValueForFieldAtIndexIsCapIs)
	s.Step(`^the array value for field "32" bit "([^"]*)" at index "(\d+)" is "(\d+\.\d+)" cap is "(\d+)"$`, theFloat32ArValueForFieldAtIndexIsCapIs)
	s.Step(`^the array value for field "64" bit "([^"]*)" at index "(\d+)" is "(\d+\.\d+)" cap is "(\d+)"$`, theFloat64ArValueForFieldAtIndexIsCapIs)
	s.Step(`^the array value for field "bool" "([^"]*)" at index "(\d+)" is "([^"]*)" cap is "(\d+)"$`, theBoolArrayValueForFieldAtIndexIsCapIs)
	s.Step(`^the array value for field "time" "([^"]*)" at index "(\d+)" is "([^"]*)" cap is "(\d+)"$`, theTimeArrayValueForFieldAtIndexIsCapIs)
}
