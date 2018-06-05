// package containing common functions and features for CFactor to work
// smoothly. FileUtil contains file I/O functions
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




