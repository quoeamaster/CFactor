Feature: TOML Access (Basic use cases)
  for a configuration framework to work;
  it NEEDS to be able to LOAD a configuration file
    (source in this case is a File, but could be database or in-memory caches);
  in some sense it should also be able to MERGE configuration changes from
    different sources (e.g. File and System Cache)

  alas, it should be able to WRITE the changes back to the source(s)

  Assumptions for the feature test:
  - data file is in TOML format
  - data file is in the current folder next to the feature file (can test on absolute path as well)

  Major use cases:
  - basic loading of the config file(s) - local and absolute path
  - access of certain fields from the loaded TOML
  - make changes to the <lastUpdated> field in the TOML (to prove write operations also working)

  Scenario: Load TOML in the current / relative path
    Given there is a TOML in the current folder named "loadBasicToml.toml"
    When I load the TOML file named "loadBasicToml.toml"
    Then I should be able to access the fields from this toml file
    And the value for field "version" is "1.1.0a"
    And the value for field "author.firstName" is "Jason"
    And the integer value for field "workingHoursDay" is 8
    And the integer value for field "author.age" is 25
    And the float value for field "author.height" is 167.5
    And the value for field "role" is "admin"
    And the bool value for field "activeProfile" is "true"
    And the time value for field "lastUpdateTime" is "2016-12-25T14:02:59+08:00"
    And the time value for field "shortDateTime" is "2016-03-13T14:12:56"
    And the time value for field "shortDate" is "2016-02-12"
    And the time value for field "author.birthday" is "1990-02-28"

