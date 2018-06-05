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

// package containing common functions and features for CFactor to work smoothly.
// FileUtil contains file I/O functions.
package common

import (
	"io/ioutil"
	"strings"
	"os"
)


// function to load a file. Return the data in []byte and
// error occurred during the operation
func LoadFile(filenameOrPath string) (data []byte, err error) {
	return ioutil.ReadFile(filenameOrPath)
}

// fucntion to parse the given []byte into "lines" ([]string)
func GetLinesFromByteArrayContent(data []byte) ([]string) {
	var sLines []string
	if data != nil {
		sData := string(data)
		sLines = strings.Split(sData, "\n")
	}
	return sLines
}

// function to create a file. Returns a file reference (*os.File)
func CreateFile(filename string) (*os.File) {
	if !IsStringEmptyOrNil(filename) {
		filePtr, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		return filePtr
	}
	return nil
}

/**
 *	helper method to remove a file by the "filename"
 *
func RemoveFile(filename string) error {
	return os.Remove(filename)
}
*/




