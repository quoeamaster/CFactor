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

	return nil
}

func checkTomlContents(numOfFields int) error {

fmt.Println("**", config)
	if numOfFields == 1 {
		return nil
	}
	return fmt.Errorf("Num of fields available in the toml file MUST be %v\n", numOfFields)
}

func checkFieldValue(field, value string) error {
	if configReader.IsFieldStringValueMatched(config, field, value) {
		return nil
	}
	return fmt.Errorf("field [%v] does not matches with {%v}", field, value)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^there is a TOML in the current folder named "([^"]*)"$`, foundATomlFileLocation)
	s.Step(`^I load the TOML file named "([^"]*)"$`, loadToml)
	s.Step(`^I should be able to access the (\d+) fields from this toml file$`, checkTomlContents)
	s.Step(`^the value for field "([^"]*)" is "([^"]*)"$`, checkFieldValue)
}
