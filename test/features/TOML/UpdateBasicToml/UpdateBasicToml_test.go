package UpdateBasicToml

import (
	"github.com/DATA-DOG/godog"
	"CFactor/TOML"
	"reflect"
	"fmt"
	TOML2 "CFactor/test/features/TOML"
	"strings"
	"CFactor/common"
)

var configReader TOML.TOMLConfigImpl
var configObject TOML2.DemoTOMLConfig

func gotTomlFileName(tomlFile string) error {
	if len(tomlFile)>0 {
		configReader = TOML.NewTOMLConfigImpl(tomlFile, reflect.TypeOf(TOML2.DemoTOMLConfig{}))
		return nil

	} else {
		return fmt.Errorf("the given 'name' is not Valid (%v)", tomlFile)
	}
}

func loadTomlFile(_ string) error {
	configObject = TOML2.DemoTOMLConfig{ Author: TOML2.Author{} }
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



func FeatureContext(s *godog.Suite) {
	s.Step(`^there is a TOML in the current folder named "([^"]*)"$`, gotTomlFileName)
	s.Step(`^I load the TOML file named "([^"]*)"$`, loadTomlFile)
	s.Step(`^by accessing the toml loaded, the value for field "([^"]*)" is "([^"]*)"$`, theValueForFieldIs)
	s.Step(`^set the "([^"]*)" to the current timestamp "([^"]*)"$`, setValueForField)
	s.Step(`^save changes to the "([^"]*)"$`, saveChangesToToml)
	s.Step(`^finally reload the configuration file "([^"]*)", "([^"]*)" should equals to "([^"]*)"$`, reconciliationOnFieldsSet)
}
