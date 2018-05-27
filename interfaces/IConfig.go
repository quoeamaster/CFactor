package interfaces

import "reflect"

/**
 *	declare the interface for a Config Reader / Writer
 */
type IConfig interface {
	/**
	 *	able to load a configuration resource name
	 *		(provided by the struct that implements this interface);
	 *	return a generic object / interface (casting is possible)
	 *
	 *	TODO: create a generics version???
	 */
	Load(name string, structType reflect.Type) (reflect.Value, error)

	/**
	 *	to save / persist the given configObject to the the given resource name
	 *	(filename); type information is required so that the correct
	 *	translation is performed
	 */
	Save(name string, structType reflect.Type, configObject interface{}) (error)
}

/**
 *	lifecycle hooks for CFactor.
 */
type IConfigLifeCycleHooks interface {
	/**
	 *	when inner fields are struct typed; this method would help to
	 *		set back these references to the field(s)
	 */
	SetStructsReferences(structRefMap *map[string]interface{}) (error)
}

const MethodSetStructsReference = "SetStructsReferences"


/**
 *	include a generic set method.
 *	provide the "key" plus an optional "params" map
 *
 *	implementation would need to handle implicit in the code (eg. date format
 *	handling and type casting into time.Time)
 */
/*
type ISetter interface {
	Set(key string, params map[string]string) (bool, error)
}


type IGetter interface {
	Get(key string) (interface{})
}
*/

