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

const TAG_TOML = "toml"
const TAG_JSON = "json"
const TAG_ADDITIONAL = "additional"
const TAG_SET = "set"

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

type TestDemo struct {
	Msg string
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
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TAG_ADDITIONAL), CONFIG_TYPE_PARENT) == 0 {
			// TODO: do not work at the moment, as the reflected value is not a real object instance... (recursively populate...)

			structObj := NewStructPointerByType(objVal.Field(i).Type())
			popStructobj := populateStringValueByFieldNameUnderChildStruct(structObj.Type().Elem(), k, v)

			structField := objVal.Field(i)
			if objVal.CanSet() && structField.CanSet() {
				methodName := tags.Get(TAG_SET)
				methodRef := reflect.ValueOf(object).MethodByName(methodName)

				inParams := make([]reflect.Value, methodRef.Type().NumIn())
				inParams[0]=reflect.ValueOf(popStructobj)
				p := make(map[string]string)
				p["demo"]="demo-value"
				inParams[1]=reflect.ValueOf(p)

				outVals := methodRef.Call(inParams)
				fmt.Println("m name", methodName, "m2", methodRef, "p1",popStructobj, "output", outVals)
/*
				s := TestDemo{ "happy"}
// TODO: t.SetAuthor (method... at the Struct side... do the casting)


				structField.Set(reflect.ValueOf(&popStructobj))
				structField.Set(reflect.ValueOf(s))
*/
			}


		} else {
			if strings.Compare(tags.Get(TAG_TOML), k) == 0 {
				// ### reflect.ValueOf(&r).Elem().Field(i).SetInt( i64 )
				objVal.Field(i).SetString(v)
				break
			}	// end -- if (k matched)
		}	// end -- if (additional_info == parent)
	}	// end -- for (fLen)
}

func populateStringValueByFieldNameUnderChildStruct(structObjType reflect.Type, k, v string) (interface{}) {
	fLen := structObjType.NumField()
	// strip the " symbol if any
	v = strings.Replace(v, "\"", "", -1)

	structObj := NewStructPointerByType(structObjType) //.Elem()
	//fmt.Println(structObj, structObj.Type(), structObj.CanSet())

	for i:=0; i<fLen; i++ {
		tags := structObjType.Field(i).Tag

		if strings.Compare(tags.Get(TAG_TOML), k) == 0 {
			rField := structObj.Elem().Field(i)
			if rField.CanSet() {
				rField.SetString(v)
			}	// end -- if (rField can set)
		}	// end -- if (tagStruct.Field == k)
	}	// end -- for (fLen)
	return structObj
}

func GetStringValueByTomlField(object interface{}, objectType reflect.Type, k string) (bool, string) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TAG_TOML), k) == 0 {
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
 *
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
*/