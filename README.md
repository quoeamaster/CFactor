# CFactor
core configuration module for any system (supports TOML, JSON to be supported soon)

goals: a library supporting configuration access (read and write) base on JSON and TOML syntax

Example: To load a configuration (e.g. toml) and populate the corresponding Struct object's fields
```golang
configReader := TOML.NewTOMLConfigImpl(name, reflect.TypeOf(TOML2.DemoTOMLConfig{}))
```

