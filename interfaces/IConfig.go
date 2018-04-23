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
}

/**
 *	include a generic set method.
 *	provide the "key" plus an optional "params" map
 *
 *	implementation would need to handle implicit in the code (eg. date format
 *	handling and type casting into time.Time)
 */
type ISetter interface {
	Set(key string, params map[string]string) (bool, error)
}


