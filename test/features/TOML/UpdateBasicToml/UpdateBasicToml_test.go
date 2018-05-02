package UpdateBasicToml

import "github.com/DATA-DOG/godog"

func gotTomlFileName(tomlFile string) error {
	return godog.ErrPending
}

func loadTomlFile(_ string) error {
	return godog.ErrPending
}

func theValueForFieldIs(field, value string) error {
	return godog.ErrPending
}

func setValueForField(field, value string) error {
	return godog.ErrPending
}

func saveChangesToToml(tomlFile string) error {
	return godog.ErrPending
}

func reconciliationOnFieldsSet(field, value string) error {
	return godog.ErrPending
}



func FeatureContext(s *godog.Suite) {
	s.Step(`^there is a TOML in the current folder named "([^"]*)"$`, gotTomlFileName)
	s.Step(`^I load the TOML file named "([^"]*)"$`, loadTomlFile)
	s.Step(`^by accessing the toml loaded, the value for field "([^"]*)" is "([^"]*)"$`, theValueForFieldIs)
	s.Step(`^set the "([^"]*)" to the current timestamp "([^"]*)"$`, setValueForField)
	s.Step(`^save changes to the "([^"]*)"$`, saveChangesToToml)
	s.Step(`^finally reload the configuration, "([^"]*)" should equals to "([^"]*)"$`, reconciliationOnFieldsSet)
}
