# CFactor
core configuration module for any system (supports TOML, JSON to be supported soon).
The goal for this project is to provide a library supporting configuration access
(read and write) base on TOML syntax (JSON would be supported soon.

Example: To load a configuration (e.g. toml) and populate the corresponding Struct object's fields
```golang
// create the configuration reader and the targeted Struct object for data population
configReader := TOML.NewTOMLConfigImpl("demo-config.toml", reflect.TypeOf(DemoStruct{}))
configObject := DemoStruct{  }
```

