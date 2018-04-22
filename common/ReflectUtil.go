package common

import (
	"reflect"
	"strings"
	"errors"
	"fmt"
)

const CONFIG_TYPE_TOML = "CFG_T"
const CONFIG_TYPE_JSON = "CFG_J"
const CONFIG_TYPE_PARENT = "parent"


type TagStructure struct {
	CType string
	Field string
	Additional string
}

/**
 *	create a pointer instance of the given "Type"
 */
func NewStructPointerByType(t reflect.Type) reflect.Value {
	if t != nil {
		return reflect.New(t)
	}
	return reflect.Zero(t)
}

func PopulateFieldValues(ln string, configType string, object interface{}, objectType reflect.Type) (bool, error) {
	// if IsValueValid(object) == true {
	if IsValidPointer(object) == true {
		// trim the ln (spaces removal)
		if len(ln)>0 {
			ln = strings.TrimSpace(ln)
		}
		// if configType is not TOML or JSON, treat it as TOML (default)
		if !(strings.Compare(CONFIG_TYPE_TOML, configType) == 0 ||
			strings.Compare(CONFIG_TYPE_JSON, configType) == 0) {
			configType = CONFIG_TYPE_TOML
		}

		// non empty line (try to process)
		if len(ln)>0 {
			kv := strings.Split(ln, "=")
			if len(kv) == 2 {
// TODO: check if this is a root level key or key under some hierarchy
				k := strings.TrimSpace(kv[0])
// TODO: might be string data type or something else
				v := strings.TrimSpace(kv[1])

				populateStringValByFieldName(object, objectType, k, v)
			}
		}	// end -- if (ln is non empty)
		return true, nil
	}
	return false, errors.New("object / value provided is non-valid")
}

/*
 *	population of a string field by fieldName
 */
func populateStringValByFieldName(object interface{}, objectType reflect.Type, k string, v string) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object).Elem()

	// strip the " symbol if any
	v = strings.Replace(v, "\"", "", -1)

	for i:=0; i<fLen; i++ {
		tag := fmt.Sprintln(objectType.Field(i).Tag)
		tagStruct := ParseTagToTagStructure(tag)

		if strings.Compare(tagStruct.Additional, CONFIG_TYPE_PARENT) == 0 {
			// TODO: do not work at the moment, as the reflected value is not a real object instance... (recursively populate...)
			/*
			structObjType := objVal.Field(i).Type()
			structObj := objVal.Field(i)
			populateStringValByFieldName(&structObj, structObjType, k, v)
			*/

		} else {
			if strings.Compare(tagStruct.Field, k) == 0 {
				objVal.Field(i).SetString(v)
				break
			}	// end -- if (k matched)
		}	// end -- if (additional_info == parent)
	}	// end -- for (fLen)
}

func GetStringValueByTomlField(object interface{}, objectType reflect.Type, k string) (bool, string) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tag := fmt.Sprintln(objectType.Field(i).Tag)
		tagStruct := ParseTagToTagStructure(tag)

		if strings.Compare(tagStruct.Field, k) == 0 {
			return true, objVal.Field(i).String() //return true, fmt.Sprint(objVal.Field(i).Interface())
		}	// end -- if (k matched)
	}	// end -- for (fLen)
	return false, ""
}


/**
 *	check if the given "value" is valid or not (sort of nil check)
 */
func IsValueValid(object reflect.Value) bool {
	return !strings.Contains(object.String(), "invalid")
}

func IsValidPointer(object interface{}) bool {
	return object != nil
}

/**
 *	method to parse a Tag (from struct) to a TagStructure instance
 */
func ParseTagToTagStructure(tag string) TagStructure {
	parts := strings.Split(tag, ":")
	s := TagStructure{}

	if len(parts) >= 2 {
		s.CType = strings.TrimSpace(parts[0])
		s.Field = strings.TrimSpace(parts[1])
	}
	if len(parts) > 2 {
		s.Additional = strings.TrimSpace(parts[2])
	}
	return s
}
