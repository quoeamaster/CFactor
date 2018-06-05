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

// Timeutil contains time.Time related functions.
package common

import (
	"time"
	"strings"
	"errors"
	"fmt"
)

/*
 * 	ISO_8601
 * 	format description from https://en.wikipedia.org/wiki/ISO_8601
 *	yyyy-mm-dd{ T - time part }hh:mm:ss{ Z => timezone }hh:mm
 *
 *	FORMAT is fixed in someways => " 2006-01-02T15:04:05Z07:00 "
 * 	2016-12-25T01:02:59+08:00 => HKT
 *	2016-12-25T01:02:59Z => UTC (* Z means UTC, special handling)
 */

// default time format
const TimeDefault = "2006-01-02T15:04:05Z07:00"
// time short format
const TimeShortDate = "2006-01-02"
// time short + date format
const TimeShortDateTime = "2006-01-02T15:04:05"

// parse a given string-formatted datetime to time.Time.
// If format is non valid (empty string); the default time format (TimeDefault) is used
func ParseStringToTime(format string, dateInString string) (time.Time, error) {
	finalFormat := validateTimeFormat(format)
	// parse
	return time.Parse(finalFormat, dateInString)
}

// parse a given string-formatted datetime to time.Time based on the
// list of patterns.
func ParseStringToTimeWithPatterns(formats []string, valueInString string) (time.Time, string, error) {
	for _, format := range formats {
		t, err := time.Parse(format, valueInString)
		if err == nil {
			return t, format, err
		} 	// end -- if (parse is valid)
	}	// end -- for (formats)
	return  time.Now(), "", errors.New(fmt.Sprintf("non matchable {%v} on the given patterns {%v}", valueInString, formats))
}

// function to parse the given time.Time reference to string according
// to the given format.
func FormatTimeToString(format string, valueInTime time.Time) string {
	finalFormat := validateTimeFormat(format)
	// format
	return valueInTime.Format(finalFormat)
}


/**
 *	simply check if the given "format" is valid or not
 *	validation is based on if "format" is non empty; no intelligent checks
 *	on the date format in general
 */

func validateTimeFormat(format string) string {
	if len(format) == 0 || len(strings.TrimSpace(format))==0 {
		return TimeDefault
	} else {
		return format
	}
}

