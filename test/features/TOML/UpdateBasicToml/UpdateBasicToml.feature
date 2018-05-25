Feature: TOML Access (Basic use cases)
  for a configuration framework to work;
  it should be able to WRITE the changes back to the source(s)

  Assumptions for the feature test:
  - data file is in TOML format
  - data file is in the current folder next to the feature file (can test on absolute path as well)

  Major use cases:
  - make changes to the <lastUpateTime> field in the TOML (to prove write operations also working)

  Scenario: 1) Load the TOML and then update the field <lastUpdateTime>; then retrieve it again to prove if worked
    Given there is a TOML in the current folder named "updateBasicToml.toml"
    When I load the TOML file named "updateBasicToml.toml"
    Then by accessing the toml loaded, the value for field "version" is "1.1 alpha"
    And set the "LastUpdateTime" to the current timestamp "2018-05-01T11:59:59+08:00"
    And save changes to the "updateBasicToml_test.toml"
    And finally reload the configuration file "updateBasicToml_test.toml", "lastUpdateTime" should equals to "2018-05-01T11:59:59+08:00"

  Scenario: 2) Persist a bunch of fields to the target TOML
    Given an in-memory configuration object;
    When persisted the changes to the toml file named "updateBasicToml_test2.toml";
    Then reload the "updateBasicToml_test2.toml" ...
    And field "WorkingHoursDay" should yield "8",
    And field "ActiveProfile" should yield "false",
    And array-field "Hobbies" should yield "badminton,soccer,cooking",
    And array-field "TaskNumbers" should yield "123,345,567",
    And field "LastUpdateTime" should yield "2016-12-25T14:02:59+08:00",
    And array-field "FloatingPoints32" should yield "12.3,56.90,67.098",
    And array-field "SpecialDates" should yield "2016-12-25T14:02:59+08:00,1998-01-01T09:02:59+00:00",

  Scenario: 3) Persist a bunch of fields to the target TOML (child struct)
    Given an in-memory configuration object with child struct;
    When persisted the changes to the toml file named "updateBasicToml_test3.toml";
    Then reload the "updateBasicToml_test3.toml" ...
    And field "WorkingHoursDay" should yield "12",
    And child field "LastName" should yield "Wong",
    And child field "Age" should yield "18",
    And child array-field "LuckyNumbers" should yield "1,23,908",
    And child array-field "Attributes64" should yield "12,990.0009",
    And child array-field "Likes" should yield "true,false,true,false,false",
    And child array-field "RegistrationDates" should yield "1998-01-30T00:00:00+08:00,1990-07-28T00:00:00+00:00",
    And child field "Height" should yield "166.5",
    And child field "Birthday" should yield "1980-01-30T00:00:00+08:00",

  Scenario: 4) Persist a bunch of fields to the target TOML (multiple levels of structs)
    Given an in-memory configuration object with multile levels of struct;
    When persisted the changes to the toml file named "updateBasicToml_test4.toml";
    Then reload the "updateBasicToml_test4.toml" ...
    And field "WorkingHoursDay" should yield "17",
    And child field "LastName" should yield "Wong",
    And child field "Age" should yield "18",
    And child field "Birthday" should yield "1980-01-30T00:00:00+08:00",
    And child array-field "LuckyNumbers" should yield "1,23,908",
    And child array-field "Likes" should yield "true,false,true,false,false",
    And child array-field "RegistrationDates" should yield "1998-01-30T00:00:00+08:00,1990-07-28T00:00:00+00:00",
    And multi child field "abc.def" should yield "2018",
    And multi child array-field "def.ghi" should yield "1784-07-28T04:16:52+00:00",


