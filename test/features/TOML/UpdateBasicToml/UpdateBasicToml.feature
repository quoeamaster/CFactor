Feature: TOML Access (Basic use cases)
  for a configuration framework to work;
  it should be able to WRITE the changes back to the source(s)

  Assumptions for the feature test:
  - data file is in TOML format
  - data file is in the current folder next to the feature file (can test on absolute path as well)

  Major use cases:
  - make changes to the <lastUpateTime> field in the TOML (to prove write operations also working)

  Scenario: Load the TOML and then update the field <lastUpdateTime>; then retrieve it again to prove if worked
    Given there is a TOML in the current folder named "updateBasicToml.toml"
    When I load the TOML file named "updateBasicToml.toml"
    Then by accessing the toml loaded, the value for field "version" is "1.1"
    And set the "lastUpdateTime" to the current timestamp "2018-05-01T11:59:59+08:00"
    And save changes to the "updateBasicToml.toml"
    And finally reload the configuration, "lastUpdateTime" should equals to "2018-05-01T11:59:59+08:00"
