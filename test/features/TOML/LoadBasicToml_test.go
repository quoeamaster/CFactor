package TOML

import (
	"github.com/DATA-DOG/godog"
	"fmt"
	"CFactor/TOML"
	"reflect"
	"strconv"
	"strings"
	"CFactor/common"
)

// class level variable
var configReader TOML.TOMLConfigImpl
var config DemoTOMLConfig

func foundATomlFileLocation(name string) error {
	// somehow you need to know the target Config object/struct's type
	if len(name)>0 {
		configReader = TOML.NewTOMLConfigImpl(name, reflect.TypeOf(DemoTOMLConfig{}))
		return nil

	} else {
		return fmt.Errorf("the given 'name' is not Valid (%v)", name)
	}
}

func loadToml(name string) error {
	// create an instance for population
	configObject := DemoTOMLConfig{ Author: Author{} }

	// no overriding parameters supplied
	cfgInterface, err := configReader.Load(&configObject)
	if err != nil {
		return fmt.Errorf("Error in loading the TOML file. %v\n", err)
	}
	config = reflect.ValueOf(cfgInterface).Elem().Interface().(DemoTOMLConfig)
	fmt.Println("\t#",config.String())

	return nil
}

func iShouldBeAbleToAccessTheFieldsFromThisTomlFile() error {
	// really just to add this "feature" line for clarity, no actions are required
	return nil
}

func checkFieldValue(field, value string) error {
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
	/*
	ok, val := configReader.GetFloatValueByKey(config, field)
	if !ok {
		return fmt.Errorf("given %v's value not FOUND", field)
	}
	if val==float64(value) {
		return nil
	}
	return fmt.Errorf("field [%v] does not matches with {%v}; value got is (%v)", field, value, val)
	*/
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
	patterns := []string { common.TIME_DEFAULT, common.TIME_SHORT_DATE_TIME, common.TIME_SHORT_DATE }
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


func FeatureContext(s *godog.Suite) {
	s.Step(`^there is a TOML in the current folder named "([^"]*)"$`, foundATomlFileLocation)
	s.Step(`^I load the TOML file named "([^"]*)"$`, loadToml)
	s.Step(`^I should be able to access the fields from this toml file$`, iShouldBeAbleToAccessTheFieldsFromThisTomlFile)
	s.Step(`^the value for field "([^"]*)" is "([^"]*)"$`, checkFieldValue)
	s.Step(`^the integer value for field "([^"]*)" is (\d+)$`, theIntegerValueForFieldIs)
	s.Step(`^the float value for field "([^"]*)" is (\d+\.\d+)$`, theFloatValueForFieldIs)
	s.Step(`^the bool value for field "([^"]*)" is "([^"]*)"$`, theBoolValueForFieldIs)
	s.Step(`^the time value for field "([^"]*)" is "([^"]*)"$`, theTimeValueForFieldIs)
}
