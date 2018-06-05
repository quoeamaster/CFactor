# CFactor
core configuration module for any system (supports TOML, JSON to be supported soon).
The goal for this project is to provide a library supporting configuration access
(read and write) base on TOML syntax (JSON would be supported soon.

Example: To load a configuration (e.g. toml) and populate the corresponding Struct object's fields
```golang
// create the configuration reader and the targeted Struct object for data population
configReader := TOML.NewTOMLConfigImpl("demo-config.toml", reflect.TypeOf(DemoStruct{}))
configObject := DemoStruct{  }

_, err := configReader.Load(&configObject)
if err != nil {
    return fmt.Errorf("Error in loading the TOML file. %v\n", err)
}

// now the configObject is populated; you can get back the config values by
// referencing the fields directly (e.g. the "FirstName" field which is a string)
fmt.Println(configObject.FirstName)
```

Example: To persist a struct's values back into a configuration (e.g. toml)
```golang
// assume configObject has already been populated
err := configReader.Save("new-demo-config.toml", reflect.TypeOf(configObject), configObject)
if err != nil {
    return fmt.Errorf("something wrong when persisting the toml file~ %v\n", err)
}

// now the values of the configObject are populated to the "new-demo-config.toml"
```

Finally, the Tag setup for the Struct object(s)
```golang
type TransactionRecord struct {
	// toml => indicates when parsed as "toml" config, Amount Field corresponds to the key "amount"
	Amount float32 `toml:"amount"`

	// additional => indicates this Field points to another Struct (hierarchical)
	Client Client `toml:"client" additional:"parent"`

	// struct to describe the "broker" involved
	Broker Broker `toml:"broker" additional:"parent"`
}

...

// the lifeCycle Hook method implementation (check IConfig.go)
// this method is REQUIRED, so that the child Structs (hierarchical) could be
// safely "set".
//
// PS. if the Struct DOES NOT contains child Struct(s), you just need to provide
//  an empty method
func (o *TransactionRecord) SetStructsReferences(structRefMap *map[string]interface{}) (err error) {
	structRefMapVal := *structRefMap
	if len(structRefMapVal)==0 {
		return nil
	}
	for key, structRef := range structRefMapVal {
		switch key {
		case "TOML.Client":
			o.Client = reflect.Indirect(reflect.ValueOf(structRef)).Interface().(Client)
		case "TOML.Broker":
			o.Broker = reflect.Indirect(reflect.ValueOf(structRef)).Interface().(Broker)
		case "TOML.ClientAddress":
			o.Client.Address = reflect.Indirect(reflect.ValueOf(structRef)).Interface().(ClientAddress)
		case "TOML.GeoPoint":
			o.Client.Address.GeoPoint = reflect.Indirect(reflect.ValueOf(structRef)).Interface().(GeoPoint)
		default:
			return fmt.Errorf("unknown struct type! [%v]", key)
		}
	}	// end -- for (structRef)

	// recovery if necessary
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return nil
}

```

A sample toml file
```golang
floatingPoints32 = [12.3,56.9,67.098]
specialDates = ["2016-12-25T14:02:59+08:00","1998-01-01T09:02:59Z"]

author.lastName = ""
author.age = 0
author.height = 0
author.birthday = "0001-01-01T00:00:00Z"
author.attributes64 = []
author.likes = []
author.firstName = ""
author.luckyNumbers = []
author.registrationDates = []

workingHoursDay = 8
activeProfile = false
hobbies = ["badminton","soccer","cooking"]
taskNumbers = [123,345,567]
lastUpdateTime = "2016-12-25T14:02:59+08:00"
```