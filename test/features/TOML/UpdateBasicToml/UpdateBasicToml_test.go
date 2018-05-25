package UpdateBasicToml

import (
	"github.com/DATA-DOG/godog"
	"CFactor/TOML"
	"reflect"
	"fmt"
	TOML2 "CFactor/test/features/TOML"
	"strings"
	"CFactor/common"
	"strconv"
	time2 "time"
)

var configReader TOML.TOMLConfigImpl
var configObject TOML2.DemoTOMLConfig

/* ------------------------------------------------------------ */
/*	scenario 1) Load the TOML and then update the field			*/
/*	 <lastUpdateTime>; then retrieve it again to prove 			*/
/*	 it worked													*/
/* ------------------------------------------------------------ */

func gotTomlFileName(tomlFile string) error {
	if len(tomlFile)>0 {
		configReader.Name = tomlFile
		return nil

	} else {
		return fmt.Errorf("the given 'name' is not Valid (%v)", tomlFile)
	}
}

func loadTomlFile(_ string) error {
	_, err := configReader.Load(&configObject)
	if err != nil {
		return fmt.Errorf("Error in loading the TOML file. %v\n", err)
	}
	//configObject = reflect.ValueOf(obj).Elem().Interface().(TOML2.DemoTOMLConfig)
	return nil
}

func theValueForFieldIs(field, value string) error {
	if strings.Compare(field, "version")==0 {
		if strings.Compare(value, configObject.Version) == 0 {
			return nil
		}
	}
	return fmt.Errorf("for field '%v', expected '%v' but got '%v'", field, value, configObject.Version)
}

func setValueForField(field, value string) error {
	if strings.Compare(field, "LastUpdateTime")==0 {
		time, err :=common.ParseStringToTime("", value)
		if err != nil {
			return fmt.Errorf("could NOT convert the given time value '%v' to a valid Time.time", value)
		}
		configObject.LastUpdateTime = time
		return nil
	}
	return fmt.Errorf("unknown error: %v", "")
}

func saveChangesToToml(tomlFile string) error {
	err := configReader.Save(tomlFile, reflect.TypeOf(configObject), configObject)
	if err != nil {
		return fmt.Errorf("could NOT save the config object => %v to file resource '%v'", configObject, tomlFile)
	}
	return nil
}

func reconciliationOnFieldsSet(filename, field, value string) error {
	// reload the config file
	configReader.Name = filename
	// use a new config object to avoid... overwrites
	configObject2 := TOML2.DemoTOMLConfig{ Author: TOML2.Author{} }
	_, err := configReader.Load(&configObject2)

	if err != nil {
		return fmt.Errorf("something wrong when loading the config file %v => %v\n", filename, err)
	}
	// verify the values
	switch field {
	case "version":
		if strings.Compare(configObject2.Version, value) == 0 {
			return nil
		} else {
			return fmt.Errorf("expected value to be [%v] BUT have [%v]\n", value, configObject2.Version)
		}
	case "lastUpdateDate":
		if strings.Compare(common.FormatTimeToString("", configObject2.LastUpdateTime), value) == 0 {
			return nil
		} else {
			return fmt.Errorf("expected value to be [%v] BUT have [%v]\n", value, common.FormatTimeToString("", configObject2.LastUpdateTime))
		}
	}
	return nil
}


/* ------------------------------------------------------------ */
/*	scenario 2) Persist a bunch of fields to the target TOML	*/
/* ------------------------------------------------------------ */

func setupScenario2() error {
	// remove the existing config file
	//common.RemoveFile(configReader.Name)

	// just setup the values according to the feature file's contents
	configObject.WorkingHoursDay = 8
	configObject.ActiveProfile = false
	configObject.Hobbies = []string{ "badminton", "soccer", "cooking" }
	configObject.TaskNumbers = []int{ 123, 345, 567 }
	configObject.FloatingPoints32 = []float32{ 12.3, 56.90, 67.098 }

	time, err := common.ParseStringToTime("", "2016-12-25T14:02:59+08:00")
	if err == nil {
		configObject.LastUpdateTime = time
	} else {
		return fmt.Errorf("something is wrong on casting the lastUpdateTime to a time.Time object => %v\n", err)
	}

	configObject.SpecialDates = make([]time2.Time, 2, 2)
	time, err = common.ParseStringToTime("", "2016-12-25T14:02:59+08:00")
	if err == nil {
		configObject.SpecialDates[0] = time
	} else {
		return fmt.Errorf("something is wrong on casting the lastUpdateTime to a time.Time object => %v\n", err)
	}
	time, err = common.ParseStringToTime("", "1998-01-01T09:02:59+00:00")
	if err == nil {
		configObject.SpecialDates[1] = time
	} else {
		return fmt.Errorf("something is wrong on casting the lastUpdateTime to a time.Time object => %v\n", err)
	}

	return nil
}

func persistConfigValuesToToml(filename string) error {
	if !common.IsStringEmptyOrNil(filename) {
		err := configReader.Save(filename, reflect.TypeOf(configObject), configObject)
		if err != nil {
			return fmt.Errorf("somethng wrong when persisting the toml file~ %v\n", err)
		}
	}
	return nil
}

func reloadConfigFile(filename string) error {
	if !common.IsStringEmptyOrNil(filename) {
		configReader.Name = filename
		// reset
		configObject = TOML2.DemoTOMLConfig{ Author: TOML2.Author{} }

		_, err := configReader.Load(&configObject)

		if err != nil {
			return fmt.Errorf("something is wrong in loading the given toml file [%v] => %v\n", filename, err)
		}
	}
	return nil
}


/*
"WorkingHoursDay" should yield "8",
And field "ActiveProfile" should yield "false",
And field "LastUpdateTime" should yield "2016-12-25T14:02:59+08:00",
 */
func fieldShouldYield(fieldName, valueInString string) error {
	switch fieldName {
	case "WorkingHoursDay":
		iVal, err := strconv.ParseInt(valueInString, 10, 32)
		if err != nil {
			return fmt.Errorf("could not convert [%v]'s value [%v] to int32\n", fieldName, valueInString)
		}
		if configObject.WorkingHoursDay != int(iVal) {
			return fmt.Errorf("expected [%v] for value [%v] BUT got [%v]\n", fieldName, iVal, configObject.WorkingHoursDay)
		}
	case "ActiveProfile":
		bVal, err := strconv.ParseBool(valueInString)
		if err != nil {
			return fmt.Errorf("could not convert [%v]'s value [%v] to bool\n", fieldName, valueInString)
		}
		if configObject.ActiveProfile != bVal {
			return fmt.Errorf("expected [%v] for value [%v] BUT got [%v]\n", fieldName, bVal, configObject.ActiveProfile)
		}
	case "LastUpdateTime":
		timeString := common.FormatTimeToString("", configObject.LastUpdateTime)
		if strings.Compare(timeString, valueInString) != 0 {
			return fmt.Errorf("expected [%v] for value [%v] BUT got [%v]\n", fieldName, valueInString, timeString)
		}

	default:
		return fmt.Errorf("non support field yet => %v", fieldName)
	}

	return nil
}

/*
And array-field "Hobbies" should yield "badminton,soccer,cooking",
And array-field "TaskNumbers" should yield "123,345,567",
And array-field "FloatingPoints32" should yield "12.3,56.90,67.098",
And array-field "SpecialDates" should yield "2016-12-25T14:02:59+08:00,1998-01-01T09:02:59+00:00"
 */
func arrayfieldShouldYield(fieldName, valueArrayInString string) error {
	switch fieldName {
	case "Hobbies":
		sArr := strings.Split(valueArrayInString, ",")
		if len(sArr) != len(configObject.Hobbies) {
			return fmt.Errorf("the length of the given field %v is [%v] which is NOT the same as the given one [%v]\n", fieldName, len(configObject.Hobbies), len(sArr))
		}
		// element by element
		for idx, val := range sArr {
			sVal := configObject.Hobbies[idx]
			if strings.Compare(sVal, val) != 0 {
				return fmt.Errorf("the values of the field: %v doesn't match. Expected [%v] BUT got [%v]\n", fieldName, val, sVal)
			}
		}
	case "TaskNumbers":
		sArr := strings.Split(valueArrayInString, ",")
		if len(sArr) != len(configObject.TaskNumbers) {
			return fmt.Errorf("the length of the given field %v is [%v] which is NOT the same as the given one [%v]\n", fieldName, len(configObject.Hobbies), len(sArr))
		}
		// element by element
		for idx, val := range sArr {
			sVal := configObject.TaskNumbers[idx]
			val, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("cannot convert %v to int\n", val)
			}
			if val != sVal {
				return fmt.Errorf("the values of the field: %v doesn't match. Expected [%v] BUT got [%v]\n", fieldName, val, sVal)
			}
		}
	case "FloatingPoints32":
		sArr := strings.Split(valueArrayInString, ",")
		if len(sArr) != len(configObject.FloatingPoints32) {
			return fmt.Errorf("the length of the given field %v is [%v] which is NOT the same as the given one [%v]\n", fieldName, len(configObject.Hobbies), len(sArr))
		}
		// element by element
		for idx, val := range sArr {
			sVal := configObject.FloatingPoints32[idx]
			val, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return fmt.Errorf("cannot convert %v to float32\n", val)
			}
			if float32(val) != sVal {
				return fmt.Errorf("the values of the field: %v doesn't match. Expected [%v] BUT got [%v]\n", fieldName, val, sVal)
			}
		}
	case "SpecialDates":
		sArr := strings.Split(valueArrayInString, ",")
		if len(sArr) != len(configObject.SpecialDates) {
			return fmt.Errorf("the length of the given field %v is [%v] which is NOT the same as the given one [%v]\n", fieldName, len(configObject.Hobbies), len(sArr))
		}
		// element by element
		for idx, val := range sArr {
			sVal := configObject.SpecialDates[idx]
			val, err := common.ParseStringToTime("", val)
			if err != nil {
				return fmt.Errorf("cannot convert %v to time.Time\n", val)
			}
			if !sVal.Equal(val) {
				return fmt.Errorf("the values of the field: %v doesn't match. Expected [%v] BUT got [%v]\n", fieldName, val, sVal)
			}
		}

	default:
		return fmt.Errorf("non support field yet => %v", fieldName)
	}
	return nil
}

/* ------------------------------------------------------------ */
/*	scenario 3) Persist a bunch of fields to the target TOML	*/
/*	  (child struct)											*/
/* ------------------------------------------------------------ */

func setupScenario3() error {
	// remove the existing config file
	//common.RemoveFile(configReader.Name)

	// just setup the values according to the feature file's contents
	configObject.WorkingHoursDay = 12

	configObject.Author.LastName = "Wong"
	configObject.Author.Age = 18
	configObject.Author.Height = 166.5

	time, err := common.ParseStringToTime("", "1980-01-30T00:00:00+08:00")
	if err != nil {
		return fmt.Errorf("could not convert %v to time.Time => %v", "1980-01-30T00:00:00+08:00", err)
	}
	configObject.Author.Birthday = time

	// array fields
	configObject.Author.LuckyNumbers = []int { 1, 23, 908 }
	configObject.Author.Attributes64 = []float64 { 12, 990.0009 }
	configObject.Author.Likes = []bool { true,false,true,false,false }

	configObject.Author.RegistrationDates = make([]time2.Time, 2, 2)
	time, err = common.ParseStringToTime("", "1998-01-30T00:00:00+08:00")
	if err != nil {
		return fmt.Errorf("could not convert %v to time.Time => %v", "1998-01-30T00:00:00+08:00", err)
	}
	configObject.Author.RegistrationDates[0] = time

	time, err = common.ParseStringToTime("", "1990-07-28T00:00:00+00:00")
	if err != nil {
		return fmt.Errorf("could not convert %v to time.Time => %v", "1990-07-28T00:00:00+00:00", err)
	}
	configObject.Author.RegistrationDates[1] = time

	return nil
}

/*
And child field "LastName" should yield "Wong",
And child field "Age" should yield "18",
And child field "Height" should yield "166.5",
And child field "Birthday" should yield "1980-01-30T00:00:00+08:00",
*/
func childFieldShouldYield(fieldName, valueInString string) error {
	//fmt.Println(fieldName, "=",valueInString)
	authorRef := configObject.Author

	switch fieldName {
	case "LastName":
		if strings.Compare(valueInString, authorRef.LastName) != 0 {
			return fmt.Errorf("child field [%v] expects [%v] BUT yielded [%v]\n", fieldName, valueInString, authorRef.LastName)
		}
	case "Age":
		sVal := strconv.Itoa(authorRef.Age)
		if strings.Compare(valueInString, sVal) != 0 {
			return fmt.Errorf("child field [%v] expects [%v] BUT yielded [%v]\n", fieldName, valueInString, authorRef.Age)
		}
	case "Height":
		sVal := fmt.Sprintf("%v", authorRef.Height)
		if strings.Compare(valueInString, sVal) != 0 {
			return fmt.Errorf("child field [%v] expects [%v] BUT yielded [%v]\n", fieldName, valueInString, authorRef.Height)
		}
	case "Birthday":
		tVal := common.FormatTimeToString("", authorRef.Birthday)
		if strings.Compare(valueInString, tVal) != 0 {
			return fmt.Errorf("child field [%v] expects [%v] BUT yielded [%v]\n", fieldName, valueInString, tVal)
		}


	default:
		return fmt.Errorf("non supported field yet [%v]\n", fieldName)
	}
	return nil
}
/*
And child array-field "LuckyNumbers" should yield "1,23,908",
And child array-field "Attributes64" should yield "12,990.0009",
And child array-field "Likes" should yield "true,false,true,false,false",
And child array-field "RegistrationDates" should yield "1998-01-30T00:00:00+08:00,1990-07-28T00:00:00+00:00",
 */
func childArrayfieldShouldYield(fieldName, valueInString string) error {
	//fmt.Println(fieldName,"=",valueInString)
	authorRef := configObject.Author

	switch fieldName {
	case "LuckyNumbers":
		sArr := strings.Split(valueInString, ",")
		if len(sArr) != len(authorRef.LuckyNumbers) {
			return fmt.Errorf("[%v] length DOES-NOT match; %v vs %v\n", fieldName, len(sArr), len(authorRef.LuckyNumbers))
		}
		for idx, sVal := range sArr {
			iVal, err := strconv.Atoi(sVal)
			if err != nil {
				return fmt.Errorf("[%v] COULD NOT be converted to numbers [%v]\n", fieldName, sVal)
			}
			if iVal != authorRef.LuckyNumbers[idx] {
				return fmt.Errorf("child Array field [%v] expects [%v] BUT yielded [%v]\n",
					fieldName, iVal, authorRef.LuckyNumbers[idx])
			}
		}	// end -- for (sArr iteration)
	case "Attributes64":
		sArr := strings.Split(valueInString, ",")
		if len(sArr) != len(authorRef.Attributes64) {
			return fmt.Errorf("[%v] length DOES-NOT match; %v vs %v\n", fieldName, len(sArr), len(authorRef.Attributes64))
		}
		for idx, sVal := range sArr {
			fVal, err := strconv.ParseFloat(sVal, 64)
			if err != nil {
				return fmt.Errorf("[%v] COULD NOT be converted to numbers [%v]\n", fieldName, sVal)
			}
			if float64(fVal) != authorRef.Attributes64[idx] {
				return fmt.Errorf("child Array field [%v] expects [%v] BUT yielded [%v]\n",
					fieldName, fVal, authorRef.Attributes64[idx])
			}
		}	// end -- for (sArr iteration)
	case "Likes":
		sArr := strings.Split(valueInString, ",")
		if len(sArr) != len(authorRef.Likes) {
			return fmt.Errorf("[%v] length DOES-NOT match; %v vs %v\n", fieldName, len(sArr), len(authorRef.Likes))
		}
		for idx, sVal := range sArr {
			bVal, err := strconv.ParseBool(sVal)
			if err != nil {
				return fmt.Errorf("[%v] COULD NOT be converted to bool [%v]\n", fieldName, sVal)
			}
			if bVal != authorRef.Likes[idx] {
				return fmt.Errorf("child Array field [%v] expects [%v] BUT yielded [%v]\n",
					fieldName, bVal, authorRef.Likes[idx])
			}
		}	// end -- for (sArr iteration)
	case "RegistrationDates":
		sArr := strings.Split(valueInString, ",")
		if len(sArr) != len(authorRef.RegistrationDates) {
			return fmt.Errorf("[%v] length DOES-NOT match; %v vs %v\n", fieldName, len(sArr), len(authorRef.RegistrationDates))
		}
		for idx, sVal := range sArr {
			tVal, err := common.ParseStringToTime("", sVal)
			if err != nil {
				return fmt.Errorf("[%v] COULD NOT be converted to bool [%v]\n", fieldName, sVal)
			}
			if !tVal.Equal(authorRef.RegistrationDates[idx]) {
				return fmt.Errorf("child Array field [%v] expects [%v] BUT yielded [%v]\n",
					fieldName, tVal, authorRef.RegistrationDates[idx])
			}
		}	// end -- for (sArr iteration)
	default:
		fmt.Println(authorRef)
		return fmt.Errorf("non supported field yet [%v]\n", fieldName)
	}
	return nil
}


func FeatureContext(s *godog.Suite) {
	// before anything is running
	s.BeforeSuite(func() {
		configReader = TOML.NewTOMLConfigImpl("", reflect.TypeOf(TOML2.DemoTOMLConfig{}))
	})
	// lifecycle hooks for scenario
	s.BeforeScenario(func(i interface{}) {
		configReader.Name = ""
		configObject = TOML2.DemoTOMLConfig{ Author: TOML2.Author{} }
	})

	// scenario 1
	s.Step(`^there is a TOML in the current folder named "([^"]*)"$`, gotTomlFileName)
	s.Step(`^I load the TOML file named "([^"]*)"$`, loadTomlFile)
	s.Step(`^by accessing the toml loaded, the value for field "([^"]*)" is "([^"]*)"$`, theValueForFieldIs)
	s.Step(`^set the "([^"]*)" to the current timestamp "([^"]*)"$`, setValueForField)
	s.Step(`^save changes to the "([^"]*)"$`, saveChangesToToml)
	s.Step(`^finally reload the configuration file "([^"]*)", "([^"]*)" should equals to "([^"]*)"$`, reconciliationOnFieldsSet)

	// scenario 2
	s.Step(`^an in-memory configuration object;$`, setupScenario2)
	s.Step(`^persisted the changes to the toml file named "([^"]*)";$`, persistConfigValuesToToml)
	s.Step(`^reload the "([^"]*)" \.\.\.$`, reloadConfigFile)
	s.Step(`^field "([^"]*)" should yield "([^"].*)",$`, fieldShouldYield)
	s.Step(`^array-field "([^"]*)" should yield "([^"].*)",$`, arrayfieldShouldYield)

	// scenario 3
	s.Step(`^an in-memory configuration object with child struct;$`, setupScenario3)
	s.Step(`^child field "([^"]*)" should yield "([^"]*)",$`, childFieldShouldYield)
	s.Step(`^child array-field "([^"]*)" should yield "([^"]*)",$`, childArrayfieldShouldYield)
}
