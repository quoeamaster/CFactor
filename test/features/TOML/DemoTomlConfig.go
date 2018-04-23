package TOML

import (
	"fmt"
	"strings"
	"errors"
	"strconv"
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

	WorkingHoursDay int `toml:"workingHoursDay"`

	// TODO: more to come...
}


type Author struct {
	FirstName string `toml:"author.firstName"`
	LastName string `toml:"author.lastName"`
	Age int `toml:"author.age"`
}


/**
 *	override to have a meaningful description of the struct / object / instance
 */
func (d *DemoTOMLConfig) String() string {
	s := fmt.Sprintf("Version => %v, WorkingHoursDay => %v, Author [struct] => %v",
		d.Version, d.WorkingHoursDay, d.Author.String())

	return s
}
/**
 *	override to have a meaningful description of the struct / object / instance
 */
func (a *Author) String() string {
	s := fmt.Sprintf("{FirstName => %v; LastName => %v; Age => %v}", a.FirstName, a.LastName, a.Age)

	return s
}

/**
 *	setter implementation
 */
func (d *DemoTOMLConfig) Set(key string, params map[string]string) (bool, error) {
	// in this case "key" could be "Author"
	if len(key)>0 && strings.Compare(key, "Author")==0 {
		// check if any existing Author struct available
		author := d.Author
		if len(author.FirstName)==0 && len(author.LastName)==0 {
			author = Author{}
			d.Author = author
		}
		// populate
		if len(params["author.firstName"])>0 {
			author.FirstName = params["author.firstName"]
		}
		if len(params["author.lastName"])>0 {
			author.LastName = params["author.lastName"]
		}
		if len(params["author.age"])>0 {
			iVal, cErr := strconv.Atoi(params["author.age"])
			if cErr != nil {
				panic(errors.New(fmt.Sprintf("author.age should be of type integer, given value => [%v]", params["author.age"])))
			}
			author.Age = iVal
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
func (d *DemoTOMLConfig) SetAuthor(p map[string]string) (bool, error) {
	return d.Set("Author", p)

	// casting (due to the design constraints, golang doesn't provide such feature on reflection casting...)
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


	//authorStruct := objVal.Interface().(Author)
	//fmt.Println("checkpoint", authorStruct)
	//d.Author = authorStruct
	*/
}
