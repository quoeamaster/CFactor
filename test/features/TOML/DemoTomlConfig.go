package TOML

import (
	"fmt"
	"strings"
	"errors"
)

/*
 *	basic tag format
 *		toml:field_name
 *
 *		toml => the config type identifier
 *		field_name => corresponding toml's field name
 *		additional => etc type of the field, date-format etc
 *		set => set method to use if necessary (good for non primitive typed fields)
 *
 *	if "additional" => "parent"; means this is a struct instead of
 *		simple type(s); hence would need to dive into it (1 level depth by default)
 *		In this case, the "set" method should be used if provided
 */

type DemoTOMLConfig struct {
	Version string `toml:"version"`

	Author Author `toml:"author" additional:"parent" set:"SetAuthor"`
	//FirstName string `toml:"author.firstName"`	// easiest way to implement "hierarchy"

	// more to come...
}


type Author struct {
	FirstName string `toml:"author.firstName"`
	LastName string `toml:"author.lastName"`
}


/**
 *	override to have a meaningful description of the struct / object / instance
 */
func (d *DemoTOMLConfig) String() string {
	s := fmt.Sprintf("Version => %v", d.Version)

	return s
}
/**
 *	override to have a meaningful description of the struct / object / instance
 */
func (a *Author) String() string {
	s := fmt.Sprintf("FirstName => %v; LastName => %v", a.FirstName, a.LastName)

	return s
}


func (d *DemoTOMLConfig) Set(key string, params map[string]string) (bool, error) {
	// in this case "key" could be "Author" means the value provided
	if len(key)>0 && strings.Compare(key, "Author")==0 {
		author := Author{
			FirstName: params["author.firstName"],
			LastName: params["author.lastName"],
		}
		d.Author = author
		return true, nil
	}
	// TODO: might have more to come...

	return false, errors.New(
		fmt.Sprintf("something wrong on creating the Author struct, key provided => [%v]\n", key))
}

/**
 *	setter for the "Author" member
 */
func (d *DemoTOMLConfig) SetAuthor(object interface{}, p map[string]string) (bool, error) {


	fmt.Println(p["demo"])


	// casting
	/*
	type ABC struct {
		NAme string
	}
	var x1 ABC = ABC{ "fuck" }
	var x2 interface{} = x1
	var x3 reflect.Value = reflect.ValueOf(x2)
		fmt.Println(x3, reflect.TypeOf(x3)) // if x3 is not reflect.Value, all works
	v1 := reflect.ValueOf(x3)
	y1 := v1.Interface().(ABC) // y will have type float64.
	fmt.Println(y1, "type?", reflect.TypeOf(y1))
	*/

fmt.Println("inside setAuthor")


	//author := Author{}


	//authorStruct := objVal.Interface().(Author)
//fmt.Println("checkpoint", authorStruct)
	//d.Author = authorStruct

	return true, nil
}
