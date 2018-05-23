package TOML

import (
	"fmt"
	"strings"
	"errors"
	"strconv"
	"time"
	"CFactor/common"
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
	Role string `toml:"role"`

	Author Author `toml:"author" additional:"parent" set:"SetAuthor" get:"GetAuthor"`
	//FirstName string `toml:"author.firstName"`	// easiest way to implement "hierarchy"

	WorkingHoursDay int  `toml:"workingHoursDay"`
	ActiveProfile   bool `toml:"activeProfile"`

	// slice / array of hobbies (in string)
	Hobbies []string `toml:"hobbies"`

	TaskNumbers []int `toml:"taskNumbers"`

	LastUpdateTime time.Time `toml:"lastUpdateTime"`
	ShortDate time.Time `toml:"shortDate"`
	ShortDateTime time.Time `toml:"shortDateTime"`

	FloatingPoints32 []float32 `toml:"floatingPoints32"`
	SpecialDates []time.Time `toml:"specialDates"`
}


type Author struct {
	FirstName string `toml:"author.firstName"`
	LastName string `toml:"author.lastName"`
	Age int `toml:"author.age"`
	Height float32 `toml:"author.height"`
	Birthday time.Time `toml:"author.birthday"`
	LuckyNumbers []int `toml:"author.luckyNumbers"`
	Attributes64 []float64 `toml:"attributes64"`
	Likes []bool `toml:"author.likes"`
	RegistrationDates []time.Time `toml:"author.registrationDates"`
}


/**
 *	override to have a meaningful description of the struct / object / instance
 */
func (d *DemoTOMLConfig) String() string {
	s := fmt.Sprintf("Version => %v, WorkingHoursDay => %v, Role => %v, ActiveProfile => %v, LastUpdateTime => %v, ShortDate => %v, ShortDateTime => %v, Hobbies => %v(%v), TaskNumbers => %v(%v), FloatingPoints32 => %v(%v), SpecialDates => %v,(%v) # Author [struct] => %v",
		d.Version, d.WorkingHoursDay, d.Role,
		d.ActiveProfile, d.LastUpdateTime.String(),
		d.ShortDate, d.ShortDateTime,
		d.Hobbies, len(d.Hobbies),
		d.TaskNumbers, len(d.TaskNumbers),
		d.FloatingPoints32, len(d.FloatingPoints32),
		d.SpecialDates, len(d.SpecialDates),
		d.Author.String())

	return s
}
/**
 *	override to have a meaningful description of the struct / object / instance
 */
func (a *Author) String() string {
	s := fmt.Sprintf("{FirstName => %v; LastName => %v; Age => %v; Height => %v, birthday => %v, luckyNumbers => %v(%v), attributes => %v[%v], likes => %v(%v), registrationDates => %v(%v) }",
		a.FirstName, a.LastName, a.Age,
		a.Height, a.Birthday, a.LuckyNumbers, len(a.LuckyNumbers),
		a.Attributes64, len(a.Attributes64),
		a.Likes, len(a.Likes),
		a.RegistrationDates, len(a.RegistrationDates) )

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

		} else if len(params["author.lastName"])>0 {
			author.LastName = params["author.lastName"]

		} else if len(params["author.age"])>0 {
			iVal, cErr := strconv.Atoi(params["author.age"])
			if cErr != nil {
				panic(errors.New(fmt.Sprintf("author.age should be of type integer, given value => [%v]", params["author.age"])))
			}
			author.Age = iVal

		} else if len(params["author.height"])>0 {
			fVal, cErr := strconv.ParseFloat(params["author.height"], 32)
			if cErr != nil {
				panic(errors.New(fmt.Sprintf("author.height should be of type float32, given value => [%v]", params["author.height"])))
			}
			author.Height = float32(fVal)

		} else if len(params["author.birthday"])>0 {
			patterns := []string{ common.TimeDefault, common.TimeShortDate, common.TimeShortDateTime}
			tVal, _, cErr := common.ParseStringToTimeWithPatterns(patterns, params["author.birthday"])
			if cErr != nil {
				panic(errors.New(fmt.Sprintf("author.birthday should be of type time.Time, given value => [%v]", params["author.birthday"])))
			}
			author.Birthday = tVal

		} else if len(params["author.luckyNumbers"])>0 {
			sValAll := params["author.luckyNumbers"]
			sVals := common.CleanseArrayedString(sValAll)
			iArr, cErr := common.ConvertStringArrayToIntArray(sVals)

			if cErr != nil {
				panic(cErr)
			}
			author.LuckyNumbers = iArr

		} else if len(params["author.attributes64"])>0 {
			sValAll := params["author.attributes64"]
			sVals := common.CleanseArrayedString(sValAll)
			f64Arr, cErr := common.ConvertStringArrayToFloat64Array(sVals)

			if cErr != nil {
				panic(cErr)
			}
			author.Attributes64 = f64Arr

		} else if len(params["author.likes"])>0 {
			sValAll := params["author.likes"]
			sVals := common.CleanseArrayedString(sValAll)
			bArr, cErr := common.ConvertStringArrayToBoolArray(sVals)

			if cErr != nil {
				panic(cErr)
			}
			author.Likes = bArr

		} else if len(params["author.registrationDates"])>0 {
			sValAll := params["author.registrationDates"]
			sVals := common.CleanseArrayedString(sValAll)
			tArr, cErr := common.ConvertStringArrayToTimeArray(sVals)

			if cErr != nil {
				panic(cErr)
			}
			author.RegistrationDates = tArr

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
