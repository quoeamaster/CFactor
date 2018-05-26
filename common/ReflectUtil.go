package common

import (
	"reflect"
	"strings"
	"errors"
	"strconv"
	"fmt"
	"time"
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
const TypeTime = "time.Time"

const TypeArrayString = "[]string"
const TypeArrayInt = "[]int"
const TypeArrayFloat32 = "[]float32"
const TypeArrayFloat64 = "[]float64"
const TypeArrayBool = "[]bool"
const TypeArrayTime = "[]time.Time"

const TypePartialMap = "map["
const TypeMapStringInterface = "map[string]interface {}"

const TypePointerSymbol = "*"

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

func getStructRefByType(structRefMap map[string]interface{}, structType reflect.Type) interface{} {
	structTypeString := structType.String()
	structRef := structRefMap[structTypeString]

	if structRef == nil {
		structRef = NewStructPointerByType(structType).Interface()
		structRefMap[structTypeString] = structRef
	}
	return structRef
}

/**
 *	helper method to check if the given string is related to an "array" syntax
 */
func isValueAnArray(value string) bool {
	// TODO: might better use regexp...
	return strings.Index(value, "[")==0 && strings.LastIndex(value, "]")==(len(value)-1)
}

func PopulateFieldValues(lines []string, configType string, object interface{}, objectType reflect.Type) (bool, error) {
	if IsValidPointer(object) == true {
		// a map for storing the inner objects / structs
		structRefMap := make(map[string]interface{})

		for _, ln := range lines {
			// trim the lines (spaces removal)
			if len(ln)>0 {
				ln = strings.TrimSpace(ln)
			}
			// if configType is not TOML or JSON, treat it as TOML (default)
			if !(strings.Compare(ConfigTypeTOML, configType) == 0 ||
				strings.Compare(ConfigTypeJSON, configType) == 0) {
				configType = ConfigTypeTOML
			}
			// non empty line (try to process)
			if len(lines)>0 {
				// ignore comments
				if strings.Index(ln, "#") == 0 {
					// * return true, nil
					continue
				}
				kv := strings.Split(ln, "=")
				if len(kv) == 2 {
					k := strings.TrimSpace(kv[0])
					v := strings.TrimSpace(kv[1])

					// check if "v" is an array
					if isValueAnArray(v) {
						// handle array population plus array type policy
						populateStringValByFieldName(object, objectType, k, v, true, &structRefMap)
fmt.Println(k, "=", v, "=>",structRefMap)
					} else {
						populateStringValByFieldName(object, objectType, k, v, false, &structRefMap)
					}	// end -- if (array???)
				}
			}	// end -- if (lines is non empty)
			// * return true, nil
		}	// end -- for (lines)

		// set back the structRef(s) if any
		if err := setStructRefsToInterface(&structRefMap, object); err != nil {
			return false, err
		}
	}
	// * return false, errors.New("object / value provided is non-valid")
	return true, nil
}

func setStructRefsToInterface(structRefMap *map[string]interface{}, object interface{}) (error) {
	structRefMapVal := *structRefMap

	if len(structRefMapVal) > 0 {
		for sKey, structRef := range structRefMapVal {
fmt.Println("ee ",sKey,"-->",structRef)
			objVal := reflect.Indirect(reflect.ValueOf(object))
			objValType := objVal.Type()
			numFields := objValType.NumField()

			for idx:=0; idx<numFields; idx++ {
				objMetaField := objValType.Field(idx)
				objMetaTag := objMetaField.Tag.Get(TagAdditional)
fmt.Println(objMetaField.Tag.Get(TagTOML))
// TODO: check parent + toml tag value for a real "MATCH"
				if strings.Compare(objMetaTag, ConfigTypeParent) == 0 {
					objField := objVal.Field(idx)
					structRefName := reflect.TypeOf(structRef).String()
fmt.Println("$ b4 structName match, ",structRefName, "vs", sKey)
					structMatchInfo, ok := isStructNameMatched(structRefName, sKey)
fmt.Println("$ after structName match", ok)
					// eg. *TOML.Author vs TOML.Author
					if !ok {
						return fmt.Errorf("%v\n", structMatchInfo)
					}
					if objField.CanSet() {
						//fmt.Println("**", structRef, reflect.TypeOf(structRef))
						//fmt.Println("**2", objField.Type().String())
						// ** cannot set a Pointer object directly to the field... (struct is non pointer???)
fmt.Println("# b4 setField with struct ", objField, "==>", reflect.ValueOf(structRef).String(), reflect.ValueOf(structRef).Type())
						objField.Set(reflect.Indirect(reflect.ValueOf(structRef)))
fmt.Println("# after setField with struct ")
					}
				}	// end -- if (found a non primitive field)
			}	// end -- for (numFields)
		}	// end -- for (map content iteration)
	}	// end -- if (len of map is > 0)
	return nil
}

/**
 *	helper method to check if the given names matched (usually either 1 of the names contain a "*"
 */
func isStructNameMatched(actualRefName, targetRefName string) (string, bool) {
	// 4 conditions => a) actual contains * only, b) target contains * only and c) both contains *, both match and no "*"
	actualRefIsPtr := strings.Index(actualRefName, TypePointerSymbol) == 0
	targetRefIsPtr := strings.Index(targetRefName, TypePointerSymbol) == 0
	actualRefIsPtrPrefixedWithPtr := false
	targetRefIsPtrPrefixedWithPtr := false

	if len(actualRefName)>1 && strings.Compare(actualRefName[1: 2], TypePointerSymbol) == 0 {
		actualRefIsPtrPrefixedWithPtr = true
	}
	if len(targetRefName)>1 && strings.Compare(targetRefName[1: 2], TypePointerSymbol) == 0 {
		targetRefIsPtrPrefixedWithPtr = true
	}

	if actualRefIsPtr && !targetRefIsPtr {
		if trimmed := actualRefName[1:]; strings.Compare(trimmed, targetRefName) == 0 {
			return "", true
		}
	} else if targetRefIsPtr && !actualRefIsPtr {
		if trimmed := targetRefName[1:]; strings.Compare(trimmed, actualRefName) == 0 {
			return "", true
		}
	} else if actualRefIsPtr && targetRefIsPtr {
		if actualRefIsPtrPrefixedWithPtr && !targetRefIsPtrPrefixedWithPtr {
			trimmedActualName := actualRefName[2:]
			trimmedTargetName := targetRefName[1:]

			if strings.Compare(trimmedActualName, trimmedTargetName) == 0 {
				return "", true
			}
		} else if targetRefIsPtrPrefixedWithPtr && !actualRefIsPtrPrefixedWithPtr {
			trimmedActualName := actualRefName[1:]
			trimmedTargetName := targetRefName[2:]

			if strings.Compare(trimmedActualName, trimmedTargetName) == 0 {
				return "", true
			}
		} else {
			return fmt.Sprintf("non supported case both name starts with a [%v] though; %v vs %v", TypePointerSymbol, actualRefName, targetRefName), false
		}
	} else {
		// maybe both doesn't contain "*"
		if strings.Compare(actualRefName, targetRefName) == 0 {
			return "", true
		}
	}
	return fmt.Sprintf("non supported case both; %v vs %v", actualRefName, targetRefName), false
}


/*
 *	population of a string field by fieldName
 */
func populateStringValByFieldName(
	object interface{}, objectType reflect.Type, key string, value string,
	isArray bool, structRefMap *map[string]interface{}) {

	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object).Elem()

	// strip the " symbol if any
	value = strings.Replace(value, "\"", "", -1)

	for i:=0; i<fLen; i++ {
		typeField := objectType.Field(i)
		tags := typeField.Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent) == 0 {
			/*
			 *	as the reflected value is not a real object instance...
			 *	NEED to use alternatives (recursively populate...)
			 */
			// check if any related struct reference already there...
			innerObjInterface := getStructRefByType(*structRefMap, objVal.Field(i).Type())
			innerObjVal := reflect.ValueOf(innerObjInterface)
			innerObjValIndirected := reflect.Indirect(innerObjVal)
			innerObjType := innerObjValIndirected.Type()
//fmt.Println("ee structRef => ", objVal.Field(i).Type(), structRefMap)
			numFieldsInner := innerObjType.NumField()

			for i2:=0; i2<numFieldsInner; i2++ {
				innerObjField := innerObjType.Field(i2)
				innerObjTags := innerObjField.Tag.Get(TagTOML)

				// check if the tags match the toml name
				if strings.Compare(innerObjTags, key) == 0 {
					setValueByDataType(
						innerObjValIndirected.Field(i2).Type().String(),
						innerObjValIndirected.Field(i2), key, value, isArray)
					break
				}	// end -- if (tags matched)
			}	// end -- for (innerObject iteration)


			// * innerObjInterface := NewStructPointerByType(objVal.Field(i).Type()).Interface()
			//innerObjType := reflect.TypeOf(innerObjInterface)
			/*
			structObj := NewStructPointerByType(objVal.Field(i).Type())
			paramsMap := populateStringValueByFieldNameUnderChildStruct(structObj.Type().Elem(), key, value)


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
			*/

		} else {
//fmt.Println("ff simple fields - ", key, "vs", value)
			if strings.Compare(tags.Get(TagTOML), key) == 0 {
				// ### reflect.ValueOf(&r).Elem().Field(i).SetInt( i64 )
				setValueByDataType(typeField.Type.String(), objVal.Field(i), key, value, isArray)
				break
			}	// end -- if (key matched)
		}	// end -- if (additional_info == parent)
	}	// end -- for (fLen)
}
func populateStringValueByFieldNameUnderChildStruct(structObjType reflect.Type, k, v string) (map[string]string) {
	// strip the " symbol if any
	v = strings.Replace(v, "\"", "", -1)
	p := make(map[string]string)
	p[k] = v

	return p
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
}

/**
 *	handy method to handle set-value operation based on dataType (sharable by TOML and JSON config)
 */
func setValueByDataType(dataType string, targetField reflect.Value, k, v string, isArray bool) {
	var sArray []string

	if isArray {
		sArray = CleanseArrayedString(v)
	}

	if strings.Compare(dataType, TypeInt) == 0 {
		iVal, cErr := strconv.Atoi(v)
		if cErr != nil {
			panic(errors.New(fmt.Sprintf("cannot convert [%v] to int type for field [%v]", v, k)))
		}
		targetField.SetInt(int64(iVal))
	} else if strings.Compare(dataType, TypeString) == 0 {
		targetField.SetString(v)

	} else if strings.Compare(dataType, TypeFloat32) == 0 || strings.Compare(dataType, TypeFloat64) == 0 {
		fVal, cErr := strconv.ParseFloat(v, 64)
		if cErr != nil {
			panic(errors.New(fmt.Sprintf("cannot convert [%v] to float32 / 64 type for field [%v]", v, k)))
		}
		targetField.SetFloat(fVal)

	} else if strings.Compare(dataType, TypeBool) == 0 {
		bVal, cErr := strconv.ParseBool(v)
		if cErr != nil {
			panic(errors.New(fmt.Sprintf("cannot convert [%v] to bool type for field [%v]", v, k)))
		}
		targetField.SetBool(bVal)

	} else if strings.Compare(dataType, TypeTime) == 0 {
		patterns := []string{TimeShortDate, TimeShortDateTime, TimeDefault}
		tVal, _, cErr := ParseStringToTimeWithPatterns(patterns, v)
		if cErr != nil {
			panic(errors.New(fmt.Sprintf("cannot convert [%v] to time.Time type for field [%v]", v, k)))
		}
		// TODO: log by level (info level or debug level)???
		//fmt.Printf("[debug] format matched for time.Time field => [%v]; time.Time value => {%v}\n", format, tVal)
		targetField.Set(reflect.ValueOf(tVal))

	} else if strings.Compare(dataType, TypeArrayString)==0 {
		// easiest... string array, no additional type conversion
		targetField.Set(reflect.ValueOf( TrimStringArrayMembers(sArray) ))

	} else if strings.Compare(dataType, TypeArrayInt)==0 {
		// conversion required
		array, err := ConvertStringArrayToIntArray(sArray)
		if err != nil {
			panic(err)
		}
		targetField.Set(reflect.ValueOf( array ))

	} else if strings.Compare(dataType, TypeArrayFloat32)==0 {
		// conversion required
		array, err := ConvertStringArrayToFloat32Array(sArray)
		if err != nil {
			panic(err)
		}
		targetField.Set(reflect.ValueOf( array ))

	} else if strings.Compare(dataType, TypeArrayFloat64)==0 {
		// conversion required
		array, err := ConvertStringArrayToFloat64Array(sArray)
		if err != nil {
			panic(err)
		}
		targetField.Set(reflect.ValueOf( array ))

	} else if strings.Compare(dataType, TypeArrayBool)==0 {
		// conversion required
		array, err := ConvertStringArrayToBoolArray(sArray)
		if err != nil {
			panic(err)
		}
		targetField.Set(reflect.ValueOf( array ))

	} else if strings.Compare(dataType, TypeArrayTime)==0 {
		// conversion required
		array, err := ConvertStringArrayToTimeArray(sArray)
		if err != nil {
			panic(err)
		}
		targetField.Set(reflect.ValueOf( array ))

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

func GetTimeValueByTomlField(object interface{}, objectType reflect.Type, k string) (bool, time.Time) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			ok, tVal := GetTimeValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), k)
			if ok {
				return ok, tVal
			}	// end -- if (ok)
		} else if strings.Compare(tags.Get(TagTOML), k) == 0 {
			tVal := objVal.Field(i).Interface().(time.Time)
			return true, tVal
		}	// end -- if (k matched)
	}	// end -- for (fLen)
	return false, time.Now()
}
func GetTimeValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, time.Time) {
	fieldType := field.Type()
	fLen := fieldType.NumField()

	for i:=0; i<fLen; i++ {
		// check if tag contains "k"
		tags := fieldType.Field(i).Tag

		if strings.Compare(tags.Get(TagTOML), k)==0 {
			innerField := field.Field(i)
			// innerField.Elem().Interface().(time.Time)
			return true, innerField.Interface().(time.Time)
		}
	}	// end -- for(numField len)
	return false, time.Now()
}

/**
 *	generic method to get back values based on the object and fieldName
 */
func GetValueByTomlFieldNType(object interface{}, objectType reflect.Type, fieldName string) interface{} {
	numFields := objectType.NumField()
	objectVal := reflect.ValueOf(object)

	for idx:=0; idx<numFields; idx++ {
		// match the field names....
		fieldRef := objectVal.Field(idx)
		fieldMetaRef := objectType.Field(idx)
		if strings.Compare(fieldName, fieldMetaRef.Name)==0 {
			// found~ based on fieldRef type ... do the casting
			indirectVal := reflect.Indirect(fieldRef)
			indirectType := indirectVal.Type()
			indirectValTypeInString := indirectVal.Type().String()

			if strings.Compare(indirectValTypeInString, TypeString) == 0 {
				// real strings MUST be surrounded by "
				return "\""+indirectVal.String()+"\""

			} else if strings.Compare(indirectValTypeInString, TypeTime) == 0 {
				// real time.Time MUST be surrounded by "
				return "\""+FormatTimeToString("", indirectVal.Interface().(time.Time))+"\""

			} else if strings.Compare(indirectValTypeInString, TypeInt) == 0 {
				return indirectVal.Interface().(int)

			} else if strings.Compare(indirectValTypeInString, TypeBool) == 0 {
				return indirectVal.Interface().(bool)

			} else if strings.Compare(indirectValTypeInString, TypeFloat32) == 0 {
				return indirectVal.Interface().(float32)

			} else if strings.Compare(indirectValTypeInString, TypeFloat64) == 0 {
				return indirectVal.Interface().(float64)

			} else if strings.Compare(indirectValTypeInString, TypeArrayString) == 0 {
				return indirectVal.Interface().([]string)

			} else if strings.Compare(indirectValTypeInString, TypeArrayTime) == 0 {
				return indirectVal.Interface().([]time.Time)

			} else if strings.Compare(indirectValTypeInString, TypeArrayInt) == 0 {
				return indirectVal.Interface().([]int)

			} else if strings.Compare(indirectValTypeInString, TypeArrayBool) == 0 {
				return indirectVal.Interface().([]bool)

			} else if strings.Compare(indirectValTypeInString, TypeArrayFloat32) == 0 {
				return indirectVal.Interface().([]float32)

			} else if strings.Compare(indirectValTypeInString, TypeArrayFloat64) == 0 {
				return indirectVal.Interface().([]float64)
			}

			// non primitive type met, probably "struct"
			return GetValueByTomlFieldNStructType(indirectVal.Interface(), indirectType)

			break
		}
	}	// end -- for (loop of all fields)

	return nil
}

func GetValueByTomlFieldNStructType(object interface{}, objectType reflect.Type) (map[string]interface{}) {
	// map with all the non-null values
	valueMap := make(map[string]interface{})
	numFields := objectType.NumField()

	for idx:=0; idx<numFields; idx++ {
		fieldMetaRef := objectType.Field(idx)

		valueMap[fieldMetaRef.Tag.Get(TagTOML)] = GetValueByTomlFieldNType(
			object,
			reflect.TypeOf(object),
			fieldMetaRef.Name)
	}
	return valueMap
}

/* ---------------------------------------- */
/*	access field value through reflection 	*/
/* ---------------------------------------- */

func IsFieldValueEmptyOrNil(object interface{}, idx int, field reflect.StructField) bool {
	valObj := reflect.ValueOf(object)
	fieldTypeString := valObj.Field(idx).Type().String()
	bMatched := false

	if strings.Compare(fieldTypeString, TypeString) == 0 {
		return IsStringEmptyOrNil(valObj.Field(idx).String())

	} else if strings.Compare(fieldTypeString, TypeTime) == 0 {
		/*
		 *	to check if time.Time is ZERO => https://golang.org/pkg/time/#Time.IsZero
		 */
		time := reflect.Indirect(valObj.Field(idx)).Interface().(time.Time)
		if !time.IsZero() {
			return false
		} else {
			return true
		}
	} else if strings.Compare(fieldTypeString, TypeInt) == 0 {
		// int's default is "0" which means never possible to be empty
		bMatched = true
		return false

	} else if strings.Compare(fieldTypeString, TypeBool) == 0 {
		// bool's default is "false" which means never possible to be empty
		bMatched = true
		return false

	} else if strings.Compare(fieldTypeString, TypeFloat32) == 0 {
		// floats default is "0.0" which means never possible to be empty
		bMatched = true
		return false

	} else if strings.Compare(fieldTypeString, TypeFloat64) == 0 {
		// floats default is "0.0" which means never possible to be empty
		bMatched = true
		return false

	}

	// *** arrays ***
	if strings.Compare(fieldTypeString, TypeArrayString) == 0 {
		arr := reflect.Indirect(valObj.Field(idx)).Interface().([]string)
		if len(arr) > 0 {
			return false
		}
		bMatched = true

	} else if strings.Compare(fieldTypeString, TypeArrayTime) == 0 {
		arr := reflect.Indirect(valObj.Field(idx)).Interface().([]time.Time)
		if len(arr) > 0 {
			return false
		}
		bMatched = true

	} else if strings.Compare(fieldTypeString, TypeArrayInt) == 0 {
		arr := reflect.Indirect(valObj.Field(idx)).Interface().([]int)
		if len(arr) > 0 {
			return false
		}
		bMatched = true

	} else if strings.Compare(fieldTypeString, TypeArrayBool) == 0 {
		arr := reflect.Indirect(valObj.Field(idx)).Interface().([]bool)
		if len(arr) > 0 {
			return false
		}
		bMatched = true

	} else if strings.Compare(fieldTypeString, TypeArrayFloat32) == 0 {
		arr := reflect.Indirect(valObj.Field(idx)).Interface().([]float32)
		if len(arr) > 0 {
			return false
		}
		bMatched = true

	} else if strings.Compare(fieldTypeString, TypeArrayFloat64) == 0 {
		arr := reflect.Indirect(valObj.Field(idx)).Interface().([]float64)
		if len(arr) > 0 {
			return false
		}
		bMatched = true
	}

	// *** non primitive types such as struct(s) ***
	if !bMatched {
		//fA, _ := reflect.TypeOf(object).FieldByName("Author")
		//fmt.Println(fieldTypeString,"$$", fA.Tag)
		//fmt.Println(fieldTypeString,"$$", field.Tag)
		if object == nil {
			return true
		} else {
			return false
		}
	}
	return true
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



