/**
 *	implement the IConfig.go interface method(s); TOML version
 */
package TOML

import (
	"reflect"
	"CFactor/common"
	"strings"
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
func (t *TOMLConfigImpl) Load(ptrConfigObject interface{}) (interface{}, error) {
	// load the contents of the given "name"
	bBytes, err := common.LoadFile(t.Name)

	if err == nil {
		// build the object based on the given Type plus populate the contents loaded into bBytes
		lines := common.GetLinesFromByteArrayContent(bBytes)
		for _, v := range lines {
			ok, err := common.PopulateFieldValues(v, common.CONFIG_TYPE_TOML, ptrConfigObject, t.StructType)
			if !ok && err!=nil {
				return ptrConfigObject, err
			}
		}
		return ptrConfigObject, nil
	}
	return reflect.Zero(t.StructType), err
}

/**
 *	check if the given field:value pair matches the given object instance
 */
func (t *TOMLConfigImpl) IsFieldStringValueMatched(object interface{}, fieldName, value string) bool {
	ok, sVal := common.GetStringValueByTomlField(object, t.StructType, fieldName)

	if ok && strings.Compare(sVal, value) == 0 {
		return true
	}
	return false
}


