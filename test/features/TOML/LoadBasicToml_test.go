package TOML

import (
	"github.com/DATA-DOG/godog"
	"fmt"
	"CFactor/TOML"
	"reflect"
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
	if configReader.IsFieldStringValueMatched(config, field, value) {
		return nil
	}
	return fmt.Errorf("field [%v] does not matches with {%v}", field, value)
}

func theIntegerValueForFieldIs(field string, value int) error {
	if configReader.IsFieldIntValueMatched(config, field, value) {
		return nil
	}
	return fmt.Errorf("field [%v] does not matches with {%v}", field, value)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^there is a TOML in the current folder named "([^"]*)"$`, foundATomlFileLocation)
	s.Step(`^I load the TOML file named "([^"]*)"$`, loadToml)
	s.Step(`^I should be able to access the fields from this toml file$`, iShouldBeAbleToAccessTheFieldsFromThisTomlFile)
	s.Step(`^the value for field "([^"]*)" is "([^"]*)"$`, checkFieldValue)
	s.Step(`^the integer value for field "([^"]*)" is (\d+)$`, theIntegerValueForFieldIs)
}
