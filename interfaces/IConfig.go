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

// package defining the common interface(s)
package interfaces

import "reflect"

// declare the interface of a configurable object.
// able to Load config data from the config file to a
// targeted Struct instance;
// able to Save changes to a config file
type IConfig interface {
	// able to load a configuration file and populate the values into the
	// targeted Struct reference.
	// Return a generic reflect.Value (casting is possible) and
	// the error occurred during the Load operation.
	Load(configFilenameOrPath string, structType reflect.Type) (reflect.Value, error)

	// able to persist the given object reference's values back into the
	// targeted configuration file.
	Save(configFilenameOrPath string, structType reflect.Type, configObject interface{}) (error)
}


// declare the interface for the lifecycle hook functions.
type IConfigLifeCycleHooks interface {
	// for Structs that are hierarchical
	// (containing fields pointing to another Struct).
	// The "parent" Struct would need to handle the logic to safely set back
	// the child Struct(s).
	// This function acts as the lifecycle hook.
	SetStructsReferences(structRefMap *map[string]interface{}) (error)
}

// declaring the lifecycle hook function's name on
// "hierarchical Struct setting"
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

