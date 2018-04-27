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
				// runtime error, can't help in most cases
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
