/**
 *	implement the IConfig.go interface method(s); TOML version
 */
package TOML

import (
	"reflect"
	"CFactor/common"
	"strings"
	"runtime"
	"time"
	"bufio"
	"fmt"
	"errors"
)

type TOMLConfigImpl struct {
	Name string
	StructType reflect.Type
}

func NewTOMLConfigImpl(name string, structType reflect.Type) TOMLConfigImpl {
	impl := TOMLConfigImpl{
		Name: name,
		StructType: structType,
	}
	return impl
}


/**
 *	load the given toml file / resource
 *
 *	return => pointer to the populated instance (created from the Type)
 */
func (t *TOMLConfigImpl) Load(ptrConfigObject interface{}) (ptr interface{}, err error) {
	// defer
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				// runtime error, check if anything could be helped to continue the program
				panic(r)
			}
			err = r.(error)
		}
	}()

	// load the contents of the given "name"
	bBytes, err := common.LoadFile(t.Name)

	if err == nil {
		// build the object based on the given Type plus populate the contents loaded into bBytes
		lines := common.GetLinesFromByteArrayContent(bBytes)
		for _, v := range lines {
			ok, err := common.PopulateFieldValues(v, common.ConfigTypeTOML, ptrConfigObject, t.StructType)
			if !ok && err!=nil {
				return ptrConfigObject, err
			}
		}
		return ptrConfigObject, nil
	}
	return reflect.Zero(t.StructType), err
}

/**
 *	to save / persist the given configObject to the the given resource name
 *	(filename); type information is required so that the correct
 *	translation is performed
 */
func (t *TOMLConfigImpl) Save(name string, structType reflect.Type, configObject interface{}) (err error) {
	err = nil
	// create a Map[string]object structure for the available config tags
	configMap := make(map[string]interface{})
	numFields := structType.NumField()

	for idx := 0; idx < numFields; idx++ {
		fieldMeta := structType.Field(idx)
		ok := common.IsFieldValueEmptyOrNil(configObject, idx, fieldMeta)
		if ok==false {
			tagValue := fieldMeta.Tag.Get(common.TagTOML)
			configMap[tagValue] = common.GetValueByTomlFieldNType(configObject, structType, fieldMeta.Name)
		}	// end -- if (handle the non-nil values)
	}	// end -- for (numFields loop)

	if len(configMap) > 0 {
		cfgFile := common.CreateFile(name)
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
				cfgLine = fmt.Sprintf("%v = %v\n", key, value)
			}

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

/* ------------------------------------ */
/*	GETTERs based on key and dataType	*/
/* ------------------------------------ */


func (t *TOMLConfigImpl) GetStringValueByKey(object interface{}, fieldName string) (bool, string) {
	return common.GetStringValueByTomlField(object, t.StructType, fieldName)
}
func (t *TOMLConfigImpl) GetIntValueByKey(object interface{}, fieldName string) (bool, int64) {
	return common.GetIntValueByTomlField(object, t.StructType, fieldName)
}
func (t *TOMLConfigImpl) GetFloatValueByKey(object interface{}, fieldName string) (bool, float64) {
	return common.GetFloatValueByTomlField(object, t.StructType, fieldName)
}
func (t *TOMLConfigImpl) GetBoolValueByKey(object interface{}, fieldName string) (bool, bool) {
	return common.GetBoolValueByTomlField(object, t.StructType, fieldName)
}
func (t *TOMLConfigImpl) GetTimeValueByKey(object interface{}, fieldName string) (bool, time.Time) {
	return common.GetTimeValueByTomlField(object, t.StructType, fieldName)
}


/**
 *	check if the given field:value pair matches the given object instance (string)
 */
func (t *TOMLConfigImpl) IsFieldStringValueMatched(object interface{}, fieldName, value string) bool {
	ok, sVal := common.GetStringValueByTomlField(object, t.StructType, fieldName)

	if ok && strings.Compare(sVal, value) == 0 {
		return true
	}
	return false
}
