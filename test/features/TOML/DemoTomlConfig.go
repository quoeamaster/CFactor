package TOML

import "fmt"

/*
 *	basic tag format
 *		toml:field_name:additional_info;
 *
 *		toml => the config type identifier
 *		field_name => corresponding toml's field name
 *		additional_info => etc type of the field, date-format etc
 *		";" => multiple settings are separated by ";"
 *
 *	if "additional_info" => "parent"; means this is a struct instead of
 *		simple type(s); hence would need to dive into it (1 level depth by default)
 */

type DemoTOMLConfig struct {
	Version string `toml:version`
	//FirstName string `toml:author.firstName`	// easiest way to implement "hierarchy"

	Author Author `toml::parent`
	// more to come...
}


type Author struct {
	FirstName string `toml:author.firstName`
	LastName string `toml:author.lastName`
}


/**
 *	override to have a meaningful description of the struct / object / instance
 */
func (d *DemoTOMLConfig) String() string {
	s := fmt.Sprintf("Version => %v", d.Version)

	return s
}
