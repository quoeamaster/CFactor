/*
 *  Copyright Project - CFactor, Author - quoeamaster, (C) 2018
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

// ReflectionUtil contains reflection related functions.
package common

import (
	"reflect"
	"strings"
	"errors"
	"strconv"
	"fmt"
	"time"
	"CFactor/interfaces"
)

// the config file is "toml"
const ConfigTypeTOML = "CFG_T"
// the config file is "json" (to be supported)
const ConfigTypeJSON = "CFG_J"
// this value indicates that the field of a Struct reference is pointing
// to another Struct reference (hierarchical)
const ConfigTypeParent = "parent"

// the Tag's value indicating the corresponding toml config's key
// (e.g. toml:"first_name")
const TagTOML = "toml"
//const TagJSON = "json"

// the Tag's key indicating additional information for this Struct's field
const TagAdditional = "additional"
// deprecated => set method; use the lifeCycle hook functions such as
// "SetStructsReferences" instead (check IConfig.go)
const TagSet = "set"
//const TagGet = "get"

// data type int
const TypeInt = "int"
// data type string
const TypeString = "string"
// data type float32
const TypeFloat32 = "float32"
// data type float64
const TypeFloat64 = "float64"
// data type bool
const TypeBool = "bool"
// data type time.Time
const TypeTime = "time.Time"

// data type for array string
const TypeArrayString = "[]string"
// data type for array int
const TypeArrayInt = "[]int"
// data type for array float32
const TypeArrayFloat32 = "[]float32"
// data type for array float64
const TypeArrayFloat64 = "[]float64"
// data type for array bool
const TypeArrayBool = "[]bool"
// data type for array time.Time
const TypeArrayTime = "[]time.Time"

// data type for map (for prefix pattern matching)
const TypePartialMap = "map["
// data type for map[string]interface {}
const TypeMapStringInterface = "map[string]interface {}"

// string presentation for a "pointer"
const TypePointerSymbol = "*"

// wraps a struct field's "Tag"
type TagStructure struct {
    // config type (toml or json)
	CType string
	// name of the struct's field
	Field string
	// additional information of the struct's field (e.g. is it a "parent")
	Additional string
}

// create a new instance based on the given "type".
// The return value is a pointer referencing the new instance.
func NewStructPointerByType(t reflect.Type) reflect.Value {
	if t != nil {
		return reflect.New(t)
	}
	return reflect.Zero(t)
}

// return a Struct reference based on the "type". If not found within the
// given map; a new instance of the targeted "type" would be
// created and returned.
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

// function to populate the targeted Struct reference field(s) based on the
// configuration lines read.
// PS. the lifeCycle hook function "SetStructsReferences" would be invoked here.
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
//fmt.Println(k, "=", v, "=>",structRefMap)
					} else {
						populateStringValByFieldName(object, objectType, k, v, false, &structRefMap)
					}	// end -- if (array???)
				}
			}	// end -- if (lines is non empty)
			// * return true, nil
		}	// end -- for (lines)

		// set back the structRef(s) if any
//fmt.Println("** final structMap", len(structRefMap), "value=>", structRefMap)
		// if err := setStructRefsToInterface(&structRefMap, object); err != nil {
		if err := setStructRefsToInterfaceByLifeCycleHooks(&structRefMap, object); err != nil {
			return false, err
		}
	}
	// * return false, errors.New("object / value provided is non-valid")
	return true, nil
}

func getLifeCycleHookMethodByName(methodName string, object interface{}) reflect.Value {
	objVal := reflect.ValueOf(object)

	return objVal.MethodByName(methodName)
}

// TODO: to remove the sticky dependency of MethodSetStructsReference !!!
func setStructRefsToInterfaceByLifeCycleHooks(structRefMap *map[string]interface{}, object interface{}) (error) {
	// use the ugly approach to setStructs through relection + method invocation
	methodVal := getLifeCycleHookMethodByName(interfaces.MethodSetStructsReference, object)
	if !methodVal.IsValid() {
		return fmt.Errorf("unknown method [%v]", interfaces.MethodSetStructsReference)
	}
	methodValType := methodVal.Type()
	// create num of arguments
	args := make([]reflect.Value, methodValType.NumIn())
	args[0] = reflect.ValueOf(structRefMap)
	outArgs := methodVal.Call(args)

	if len(outArgs) > 0 {
		if !outArgs[0].IsNil() {
			return outArgs[0].Interface().(error)
		}
	}
	return nil
}

/*
func setStructRefsToInterface(structRefMap *map[string]interface{}, object interface{}) (error) {
	structRefMapVal := *structRefMap

	if len(structRefMapVal) == 0 {
		return nil
	}

	for sKey, structRef := range structRefMapVal {
//fmt.Println("ee ",sKey,"-->",structRef, reflect.TypeOf(structRef))
		objVal := reflect.Indirect(reflect.ValueOf(object))
		objValType := objVal.Type()
		numFields := objValType.NumField()

		for idx:=0; idx<numFields; idx++ {
			objMetaField := objValType.Field(idx)
			objMetaTag := objMetaField.Tag.Get(TagAdditional)
//fmt.Println(objMetaField.Tag.Get(TagTOML), " type ->", objVal.Field(idx).Type())
			// check parent + toml tag value for a real "MATCH"
			if strings.Compare(objMetaTag, ConfigTypeParent) == 0 {
				objField := objVal.Field(idx)
				if strings.Compare(objVal.Field(idx).Type().String(), sKey) == 0 {
					structRefName := reflect.TypeOf(structRef).String()
					//fmt.Println("$ b4 structName match, ",structRefName, "vs", sKey)
					structMatchInfo, ok := isStructNameMatched(structRefName, sKey)
					//fmt.Println("$ after structName match", ok)
					// eg. *TOML.Author vs TOML.Author
					if !ok {
						return fmt.Errorf("%v\n", structMatchInfo)
					}
//fmt.Println(objField.Addr())
					if objField.CanSet() {
						// ** cannot set a Pointer object directly to the field... (struct is non pointer???)
						objField.set(reflect.Indirect(reflect.ValueOf(structRef)))
						//fmt.Println("# after setField with struct ")
					}
				} else {
					// if not match, try to check recursively if this struct might contains the required struct
					//objPtr := objField.Interface()
					err := setStructRefsToInterface(structRefMap, objField.Interface())
					if err != nil {
						return err
					}
				}	// end -- if (keys matched)
			}	// end -- if (found a non primitive field)
		}	// end -- for (numFields)
	}	// end -- for (map content iteration)

	return nil
}
*/

/**
 *	helper method to check if the given names matched (usually either 1 of the names contain a "*"
 *
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
*/

/*
 *	population of a string field by fieldName
 */

func populateStringValByFieldName(
	object interface{}, objectType reflect.Type, key string, value string,
	isArray bool, structRefMap *map[string]interface{}) {

	fLen := objectType.NumField()
	//objVal := reflect.ValueOf(object).Elem()
	objVal := reflect.Indirect(reflect.ValueOf(object))
//fmt.Println("aa", objVal, objectType, " -> ", key, value, fLen)
	// strip the " symbol if any
	value = strings.Replace(value, "\"", "", -1)

	for i:=0; i<fLen; i++ {
		typeField := objectType.Field(i)
		tags := typeField.Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent) == 0 {
//fmt.Println("bb", key,"=", value," tagToml =>", tags.Get(TagTOML), " FIELDNAME => ", typeField.Name, " type =>", objectType)
			/*
			 *	as the reflected value is not a real object instance...
			 *	NEED to use alternatives (recursively populate...)
			 */
			// check if any related struct reference already there...
			innerObjInterface := getStructRefByType(*structRefMap, objVal.Field(i).Type())
			innerObjVal := reflect.ValueOf(innerObjInterface)
			innerObjValIndirected := reflect.Indirect(innerObjVal)
			innerObjType := innerObjValIndirected.Type()
			numFieldsInner := innerObjType.NumField()

			for i2:=0; i2<numFieldsInner; i2++ {
				innerObjField := innerObjType.Field(i2)
				innerObjTags := innerObjField.Tag.Get(TagTOML)
				// TODO: check if it is parent as well??? recursive calling method to populate values (print key, value for confirm)
				if strings.Compare(ConfigTypeParent, innerObjField.Tag.Get(TagAdditional)) == 0 {
					//innerStructObj := innerObjValIndirected.Field(i2).Interface()
					innerStructObj := innerObjValIndirected.Interface()
//fmt.Println("ff", reflect.TypeOf(innerStructObj))
					populateStringValByFieldName(innerStructObj, reflect.TypeOf(innerStructObj), key, value, isValueAnArray(value), structRefMap)
//fmt.Println("cc", innerStructObj, "typeof-", reflect.TypeOf(innerStructObj), "map =", structRefMap)
				} else {
					//fmt.Println("bb inner", innerObjField.Tag.Get(TagAdditional), innerObjField.Name)
					// check if the tags match the toml name
					if strings.Compare(innerObjTags, key) == 0 {
						setValueByDataType(
							innerObjValIndirected.Field(i2).Type().String(),
							innerObjValIndirected.Field(i2), key, value, isArray)
						break
					}	// end -- if (tags matched)
				}	// end -- if (parent found ??)
			}	// end -- for (innerObject iteration)
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

/*
func populateStringValueByFieldNameUnderChildStruct(structObjType reflect.Type, k, v string) (map[string]string) {
	// strip the " symbol if any
	v = strings.Replace(v, "\"", "", -1)
	p := make(map[string]string)
	p[k] = v

	return p

	// old approch
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
}
*/

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

// return the string value of the Struct reference's field (identified by "key")
func GetStringValueByTomlField(object interface{}, objectType reflect.Type, key string) (bool, string) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			ok, sVal := getStringValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), key)
			if ok {
				return ok, sVal
			}
		} else if strings.Compare(tags.Get(TagTOML), key) == 0 {
			return true, objVal.Field(i).String() //return true, fmt.Sprint(objVal.Field(i).Interface())
		}	// end -- if (key matched)
	}	// end -- for (fLen)
	return false, ""
}
func getStringValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, string) {
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

// return the int value of the Struct reference's field (identified by "key")
func GetIntValueByTomlField(object interface{}, objectType reflect.Type, key string) (bool, int64) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			//fmt.Println(tags, "v", key, "b", reflect.ValueOf(object).Field(i))
			ok, iVal := getIntValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), key)
			if ok {
				return ok, iVal
			}	// end -- if (ok)
		} else if strings.Compare(tags.Get(TagTOML), key) == 0 {
			return true, reflect.ValueOf(objVal.Field(i).Int()).Interface().(int64)
		}	// end -- if (key matched)
	}	// end -- for (fLen)
	return false, -1
}
func getIntValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, int64) {
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

// return the float value of the Struct reference's field (identified by "key")
func GetFloatValueByTomlField(object interface{}, objectType reflect.Type, key string) (bool, float64) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			ok, iVal := getFloatValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), key)
			if ok {
				return ok, iVal
			}	// end -- if (ok)
		} else if strings.Compare(tags.Get(TagTOML), key) == 0 {
			return true, reflect.ValueOf(objVal.Field(i).Float()).Interface().(float64)
		}	// end -- if (key matched)
	}	// end -- for (fLen)
	return false, -1
}
func getFloatValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, float64) {
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

// return the bool value of the Struct reference's field (identified by "key")
func GetBoolValueByTomlField(object interface{}, objectType reflect.Type, key string) (bool, bool) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			ok, bVal := getBoolValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), key)
			if ok {
				return ok, bVal
			}	// end -- if (ok)
		} else if strings.Compare(tags.Get(TagTOML), key) == 0 {
			return true, reflect.ValueOf(objVal.Field(i).Bool()).Interface().(bool)
		}	// end -- if (key matched)
	}	// end -- for (fLen)
	return false, false
}
func getBoolValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, bool) {
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

// return the time.Time value of the Struct reference's field (identified by "key")
func GetTimeValueByTomlField(object interface{}, objectType reflect.Type, key string) (bool, time.Time) {
	fLen := objectType.NumField()
	objVal := reflect.ValueOf(object)

	for i:=0; i<fLen; i++ {
		tags := objectType.Field(i).Tag

		if strings.Compare(tags.Get(TagAdditional), ConfigTypeParent)==0 {
			// met a "parent" level field
			ok, tVal := getTimeValueByTomlFieldUnderChildStruct(reflect.ValueOf(object).Field(i), key)
			if ok {
				return ok, tVal
			}	// end -- if (ok)
		} else if strings.Compare(tags.Get(TagTOML), key) == 0 {
			tVal := objVal.Field(i).Interface().(time.Time)
			return true, tVal
		}	// end -- if (key matched)
	}	// end -- for (fLen)
	return false, time.Now()
}
func getTimeValueByTomlFieldUnderChildStruct(field reflect.Value, k string) (bool, time.Time) {
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

// generic method to get back values based on the object and fieldName
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
			return getValueByTomlFieldNStructType(indirectVal.Interface(), indirectType)

			break
		}
	}	// end -- for (loop of all fields)

	return nil
}

func getValueByTomlFieldNStructType(object interface{}, objectType reflect.Type) (map[string]interface{}) {
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

// function to check if the struct object's field at index "idx"
// is empty or nil
func IsFieldValueEmptyOrNil(object interface{}, idx int) bool {
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
/*
func IsValueValid(object reflect.Value) bool {
	return !strings.Contains(object.String(), "invalid")
}
*/

// function to check if the given "object" is a valid pointer
func IsValidPointer(object interface{}) bool {
	return object != nil
}



