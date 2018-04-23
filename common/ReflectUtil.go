package common

import (
	"reflect"
	"strings"
	"errors"
	"strconv"
	"fmt"
)

const ConfigTypeTOML = "CFG_T"
const ConfigTypeJSON = "CFG_J"
const ConfigTypeParent = "parent"

const TagTOML = "toml"
const TagJSON = "json"
const TagAdditional = "additional"
const TagSet = "set"
//const TagGet = "get"

const TypeInt = "int"
const TypeString = "string"
const TypeFloat32 = "float32"
const TypeFloat64 = "float64"
const TypeBool = "bool"

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
	if IsValidPointer(object) == true {
		// trim the ln (spaces removal)
		if len(ln)>0 {
			ln = strings.TrimSpace(ln)
		}
		// if configType is not TOML or JSON, treat it as TOML (default)
		if !(strings.Compare(ConfigTypeTOML, configType) == 0 ||
			strings.Compare(ConfigTypeJSON, configType) == 0) {
			configType = ConfigTypeTOML
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
		typeField := objectType.Field(i)
		tags := typeField.Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent) == 0 {
			/*
			 *	as the reflected value is not a real object instance...
			 *	NEED to use alternatives (recursively populate...)
			 */
			structObj := NewStructPointerByType(objVal.Field(i).Type())
			paramsMap := populateStringValueByFieldNameUnderChildStruct(structObj.Type().Elem(), k, v)

			structField := objVal.Field(i)
			if objVal.CanSet() && structField.CanSet() {
				methodName := tags.Get(TagSet)
				methodRef := reflect.ValueOf(object).MethodByName(methodName)

				inParams := make([]reflect.Value, methodRef.Type().NumIn())
				inParams[0]=reflect.ValueOf(paramsMap)

				outVals := methodRef.Call(inParams)
				if outVals[0].Bool() == false {
					panic(outVals[1])
				}	// end -- if (have error)
			}	// end -- if (canSet)

		} else {
			if strings.Compare(tags.Get(TagTOML), k) == 0 {
				// ### reflect.ValueOf(&r).Elem().Field(i).SetInt( i64 )
				setValueByDataType(typeField.Type.String(), objVal.Field(i), k, v)
				break
			}	// end -- if (k matched)
		}	// end -- if (additional_info == parent)
	}	// end -- for (fLen)
}
func populateStringValueByFieldNameUnderChildStruct(structObjType reflect.Type, k, v string) (map[string]string) {
	// strip the " symbol if any
	v = strings.Replace(v, "\"", "", -1)

	/*
	fLen := structObjType.NumField()

	structObj := NewStructPointerByType(structObjType) //.Elem()
	//fmt.Println(structObj, structObj.Type(), structObj.CanSet())

	for i:=0; i<fLen; i++ {
		tags := structObjType.Field(i).Tag

		if strings.Compare(tags.Get(TagTOML), k) == 0 {
			rField := structObj.Elem().Field(i)
			if rField.CanSet() {
				rField.SetString(v)
			}	// end -- if (rField can set)
		}	// end -- if (tagStruct.Field == k)
	}	// end -- for (fLen)
	return structObj
	*/

	p := make(map[string]string)
	p[k] = v

	return p
}

/**
 *	handy method to handle set-value operation based on dataType (sharable by TOML and JSON config)
 */
func setValueByDataType(dataType string, targetField reflect.Value, k, v string) {
	if strings.Compare(dataType, TypeInt)==0 {
		iVal, cErr := strconv.Atoi(v)
		if cErr != nil {
			panic(errors.New(fmt.Sprintf("cannot convert [%v] to int type for field [%v]", v, k)))
		}
		targetField.SetInt(int64(iVal))
	} else if strings.Compare(dataType, TypeString)==0 {
		targetField.SetString(v)
	} else if strings.Compare(dataType, TypeFloat32)==0 || strings.Compare(dataType, TypeFloat64)==0 {
		fVal, cErr := strconv.ParseFloat(v, 64)
		if cErr != nil {
			panic(errors.New(fmt.Sprintf("cannot convert [%v] to float32 / 64 type for field [%v]", v, k)))
		}
		targetField.SetFloat(fVal)
	} else if strings.Compare(dataType, TypeBool)==0 {
		bVal, cErr := strconv.ParseBool(v)
		if cErr != nil {
			panic(errors.New(fmt.Sprintf("cannot convert [%v] to bool type for field [%v]", v, k)))
		}
		targetField.SetBool(bVal)
	} else {
		panic(errors.New(fmt.Sprintf("unknown type / value for field [%v] = [%v]", k, v)))
	}
}


/* ---------------------------------------- */
/*	get value for TOML based on dataType	*/
/* ---------------------------------------- */

func GetStringValueByTomlField(object interface{}, objectType reflect.Type, k string) (bool, string) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			ok, sVal := GetStringValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), k)
			if ok {
				return ok, sVal
			}
		} else if strings.Compare(tags.Get(TagTOML), k) == 0 {
			return true, objVal.Field(i).String() //return true, fmt.Sprint(objVal.Field(i).Interface())
		}	// end -- if (k matched)
	}	// end -- for (fLen)
	return false, ""
}
func GetStringValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, string) {
	fieldType := field.Type()
	fLen := fieldType.NumField()

	for i:=0; i<fLen; i++ {
		// check if tag contains "k"
		tags := fieldType.Field(i).Tag

		if strings.Compare(tags.Get(TagTOML), k)==0 {
			innerField := field.Field(i)
			return true, innerField.String()
		}
	}	// end -- for(numField len)
	return false, "not found"
}

func GetIntValueByTomlField(object interface{}, objectType reflect.Type, k string) (bool, int64) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			//fmt.Println(tags, "v", k, "b", reflect.ValueOf(object).Field(i))
			ok, iVal := GetIntValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), k)
			if ok {
				return ok, iVal
			}	// end -- if (ok)
		} else if strings.Compare(tags.Get(TagTOML), k) == 0 {
			return true, reflect.ValueOf(objVal.Field(i).Int()).Interface().(int64)
		}	// end -- if (k matched)
	}	// end -- for (fLen)
	return false, -1
}
func GetIntValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, int64) {
	fieldType := field.Type()
	fLen := fieldType.NumField()

	for i:=0; i<fLen; i++ {
		// check if tag contains "k"
		tags := fieldType.Field(i).Tag

		if strings.Compare(tags.Get(TagTOML), k)==0 {
			innerField := field.Field(i)
			return true, innerField.Int()
		}
	}	// end -- for(numField len)
	return false, -1
}

func GetFloatValueByTomlField(object interface{}, objectType reflect.Type, k string) (bool, float64) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			ok, iVal := GetFloatValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), k)
			if ok {
				return ok, iVal
			}	// end -- if (ok)
		} else if strings.Compare(tags.Get(TagTOML), k) == 0 {
			return true, reflect.ValueOf(objVal.Field(i).Float()).Interface().(float64)
		}	// end -- if (k matched)
	}	// end -- for (fLen)
	return false, -1
}
func GetFloatValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, float64) {
	fieldType := field.Type()
	fLen := fieldType.NumField()

	for i:=0; i<fLen; i++ {
		// check if tag contains "k"
		tags := fieldType.Field(i).Tag

		if strings.Compare(tags.Get(TagTOML), k)==0 {
			innerField := field.Field(i)
			return true, innerField.Float()
		}
	}	// end -- for(numField len)
	return false, -1
}

func GetBoolValueByTomlField(object interface{}, objectType reflect.Type, k string) (bool, bool) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			ok, bVal := GetBoolValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), k)
			if ok {
				return ok, bVal
			}	// end -- if (ok)
		} else if strings.Compare(tags.Get(TagTOML), k) == 0 {
			return true, reflect.ValueOf(objVal.Field(i).Bool()).Interface().(bool)
		}	// end -- if (k matched)
	}	// end -- for (fLen)
	return false, false
}
func GetBoolValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, bool) {
	fieldType := field.Type()
	fLen := fieldType.NumField()

	for i:=0; i<fLen; i++ {
		// check if tag contains "k"
		tags := fieldType.Field(i).Tag

		if strings.Compare(tags.Get(TagTOML), k)==0 {
			innerField := field.Field(i)
			return true, innerField.Bool()
		}
	}	// end -- for(numField len)
	return false, false
}

/* ---------------------------- */
/*	check validity of type(s)	*/
/* ---------------------------- */


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