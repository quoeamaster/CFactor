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

// package TOML includes an implementation of IConfig interface for toml
// config files.
package TOML

import (
	"reflect"
	"strings"
	"runtime"
	"time"
	"bufio"
	"fmt"
	"errors"
	"bytes"
	"github.com/quoeamaster/CFactor/common"
)

// "CFactor/common"

// struct wrapping the meta data for configuration loading / persisting
type TOMLConfigImpl struct {
    // filename or filepath of the config file
	Name string

	// the Struct's type in which the contents of the config file would be
	// translated into. Simply the corresponding fields of the
	// supplied Struct would be populated accordingly.
	StructType reflect.Type
}

// create a new TOMLConfigImpl instance.
func NewTOMLConfigImpl(name string, structType reflect.Type) TOMLConfigImpl {
	impl := TOMLConfigImpl{
		Name: name,
		StructType: structType,
	}
	return impl
}


// load the toml config file based on TOMLConfigImpl.Name property.
// A reference of the targeted Struct Type is given; this reference's fields
// would be populated accordingly based on the targeted Struct's Tag setup.
// Returns the same reference plus any Error occurred during the
// loading operation.
func (t *TOMLConfigImpl) Load(ptrConfigObject interface{}) (ptr interface{}, err error) {
	// defer
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				// runtime error, check if anything could be helped to continue the program
				panic(r)
			}
			//err = r.(error)
			err = fmt.Errorf("%v", r)
		}
	}()

	// load the contents of the given "name"
	bBytes, err := common.LoadFile(t.Name)

	if err == nil {
		// build the object based on the given Type plus populate the contents loaded into bBytes
		lines := common.GetLinesFromByteArrayContent(bBytes)
		ok, err := common.PopulateFieldValues(lines, common.ConfigTypeTOML, ptrConfigObject, t.StructType)
		if !ok && err!=nil {
			return ptrConfigObject, err
		}
		/*
		for _, v := range lines {
			ok, err := common.PopulateFieldValues(v, common.ConfigTypeTOML, ptrConfigObject, t.StructType)
			if !ok && err!=nil {
				return ptrConfigObject, err
			}
		}*/
		return ptrConfigObject, nil
	}
	return reflect.Zero(t.StructType), err
}

// persist the provided Struct reference's fields value back to the
// config file. Return the error occurred during the operation.
func (t *TOMLConfigImpl) Save(configFilenameOrPath string, structType reflect.Type, configObject interface{}) (err error) {
	err = nil
	// create a Map[string]object structure for the available config tags
	configMap := make(map[string]interface{})
	numFields := structType.NumField()

	for idx := 0; idx < numFields; idx++ {
		fieldMeta := structType.Field(idx)
		ok := common.IsFieldValueEmptyOrNil(configObject, idx) // , fieldMeta
		if ok==false {
			tagValue := fieldMeta.Tag.Get(common.TagTOML)
			configMap[tagValue] = common.GetValueByTomlFieldNType(configObject, structType, fieldMeta.Name)
		}	// end -- if (handle the non-nil values)
	}	// end -- for (numFields loop)

	if len(configMap) > 0 {
		cfgFile := common.CreateFile(configFilenameOrPath)
		cfgWriter := bufio.NewWriter(cfgFile)
		// sort of finally clause
		defer func() {
			// display some info about the error???
			cfgWriter.Flush()
			cfgFile.Close()
			/*
			 *	try to catch the err and return out to the caller instead of panic =>
			 *	https://stackoverflow.com/questions/19934641/go-returning-from-defer
			 */
			if r := recover(); r != nil {
				switch r.(type) {
				case error:
					err = r.(error)
				default:
					errLine := fmt.Sprintf("unknown error => %v", r)
					err = errors.New(errLine)
				}
			}
		}()

		for key, value := range configMap {
			cfgLine := ""
			// check if it is array (has different format)
			cfgLine, bMatched := translateArrayValueToStringFormat(value, key)
			if !bMatched {
				// check if it is non primitive type such as struct
				cfgLine, bMatched, err = translateNonPrimitiveValueToString(value)
				if err != nil {
					return err
				}
				if !bMatched {
					cfgLine = fmt.Sprintf("%v = %v\n", key, value)
				}	// end -- if (non array + non primitive)
			}	// end -- if (non array)

			_, err := cfgWriter.WriteString(cfgLine)
			if err != nil {
				return err
			}
		}	// end -- for (all entries inside the config map
	}	// end -- if (configMap has some elements)
	return err
}

func translateArrayValueToStringFormat(value interface{}, key string) (string, bool) {
	var cfgLine string
	bMatched := false
	sType := reflect.TypeOf(value).String()

	if strings.Compare(sType, common.TypeArrayString) == 0 {
		// cast
		sArr := value.([]string)
		sArrLine := "["
		for idx2, sVal := range sArr {
			if idx2 > 0 {
				sArrLine += ","
			}
			sArrLine += "\"" + sVal + "\""
		}
		sArrLine += "]"
		cfgLine = fmt.Sprintf("%v = %v\n", key, sArrLine)
		bMatched = true

	} else if strings.Compare(sType, common.TypeArrayInt) == 0 {
		// cast
		sArr := value.([]int)
		sArrLine := "["
		for idx2, iVal := range sArr {
			if idx2 > 0 {
				sArrLine += ","
			}
			sArrLine += fmt.Sprintf("%v", iVal)
		}
		sArrLine += "]"
		cfgLine = fmt.Sprintf("%v = %v\n", key, sArrLine)
		bMatched = true

	} else if strings.Compare(sType, common.TypeArrayTime) == 0 {
		// cast
		sArr := value.([]time.Time)
		sArrLine := "["
		for idx2, iVal := range sArr {
			if idx2 > 0 {
				sArrLine += ","
			}
			sArrLine += "\"" + common.FormatTimeToString("", iVal) + "\""
		}
		sArrLine += "]"
		cfgLine = fmt.Sprintf("%v = %v\n", key, sArrLine)
		bMatched = true

	} else if strings.Compare(sType, common.TypeArrayBool) == 0 {
		// cast
		sArr := value.([]bool)
		sArrLine := "["
		for idx2, iVal := range sArr {
			if idx2 > 0 {
				sArrLine += ","
			}
			sArrLine += fmt.Sprintf("%v", iVal)
		}
		sArrLine += "]"
		cfgLine = fmt.Sprintf("%v = %v\n", key, sArrLine)
		bMatched = true

	} else if strings.Compare(sType, common.TypeArrayFloat32) == 0 {
		// cast
		sArr := value.([]float32)
		sArrLine := "["
		for idx2, iVal := range sArr {
			if idx2 > 0 {
				sArrLine += ","
			}
			sArrLine += fmt.Sprintf("%v", iVal)
		}
		sArrLine += "]"
		cfgLine = fmt.Sprintf("%v = %v\n", key, sArrLine)
		bMatched = true

	} else if strings.Compare(sType, common.TypeArrayFloat64) == 0 {
		// cast
		sArr := value.([]float64)
		sArrLine := "["
		for idx2, iVal := range sArr {
			if idx2 > 0 {
				sArrLine += ","
			}
			sArrLine += fmt.Sprintf("%v", iVal)
		}
		sArrLine += "]"
		cfgLine = fmt.Sprintf("%v = %v\n", key, sArrLine)
		bMatched = true

	}
	return cfgLine, bMatched
}

func translateNonPrimitiveValueToString(value interface{}) (string, bool, error) {
	var err error
	bMatched := false
	sType := reflect.TypeOf(value).String()
	var bBuffer bytes.Buffer	// create a bytes.Buffer for string concatenation

	// is it a map?
	if strings.Index(sType, common.TypePartialMap) != -1 {
		if strings.Compare(common.TypeMapStringInterface, sType) == 0 {
			// translation (we only handle map[string]interface{} type for now)
			valueMap := value.(map[string]interface{})

			for mKey, mVal := range valueMap {
				cfgLine, bMatched2 := translateArrayValueToStringFormat(mVal, mKey)
				if !bMatched2 {
					// handle non array + non primitive (struct)
					cfgLine, bMatched2, err = translateNonPrimitiveValueToString(mVal)
					if err != nil {
						return "", false, err
					} else {
						_, err := bBuffer.WriteString(cfgLine)
						if err != nil {
							return "", false, err
						}
					}
					// primitive probably
					if !bMatched2 {
						_, err := bBuffer.WriteString(fmt.Sprintf("%v = %v\n", mKey, mVal))
						if err != nil {
							return "", false, err
						}
					}	// end -- if (primitive probably)
				} else {
					_, err := bBuffer.WriteString(cfgLine)
					if err != nil {
						return "", false, err
					}
				}	// end -- if (non array)
			}	// end -- for (valueMap)
			bMatched = true

		} else {
			panic(fmt.Sprintf("currently we only support map types of => %v\n", common.TypeMapStringInterface))
		}
	}	// end -- if (map type)
	return bBuffer.String(), bMatched, nil
}

/* ------------------------------------ */
/*	GETTERs based on key and dataType	*/
/* ------------------------------------ */

// deprecated method => get the string value based on a given key and
// then extract the value corresponding to the key at runtime.
func (t *TOMLConfigImpl) GetStringValueByKey(object interface{}, fieldName string) (bool, string) {
	return common.GetStringValueByTomlField(object, t.StructType, fieldName)
}
// deprecated method => get the int value based on a given key and
// then extract the value corresponding to the key at runtime.
func (t *TOMLConfigImpl) GetIntValueByKey(object interface{}, fieldName string) (bool, int64) {
	return common.GetIntValueByTomlField(object, t.StructType, fieldName)
}
// deprecated method => get the float value based on a given key and
// then extract the value corresponding to the key at runtime.
func (t *TOMLConfigImpl) GetFloatValueByKey(object interface{}, fieldName string) (bool, float64) {
	return common.GetFloatValueByTomlField(object, t.StructType, fieldName)
}
// deprecated method => get the bool value based on a given key and
// then extract the value corresponding to the key at runtime.
func (t *TOMLConfigImpl) GetBoolValueByKey(object interface{}, fieldName string) (bool, bool) {
	return common.GetBoolValueByTomlField(object, t.StructType, fieldName)
}
// deprecated method => get the time.Time value based on a given key and
// then extract the value corresponding to the key at runtime.
func (t *TOMLConfigImpl) GetTimeValueByKey(object interface{}, fieldName string) (bool, time.Time) {
	return common.GetTimeValueByTomlField(object, t.StructType, fieldName)
}


// check if the field's value of the reference object equals to the given "value" (string)
func (t *TOMLConfigImpl) IsFieldStringValueMatched(object interface{}, fieldName, value string) bool {
	ok, sVal := common.GetStringValueByTomlField(object, t.StructType, fieldName)

	if ok && strings.Compare(sVal, value) == 0 {
		return true
	}
	return false
}
