package common

import (
	"io/ioutil"
	"strings"
	"os"
)

/**
 *	method to simply load a file from the given "name"
 */
func LoadFile(name string) (data []byte, err error) {
	return ioutil.ReadFile(name)
}

/**
 *	return the "lines" of the given []byte; if valid
 */
func GetLinesFromByteArrayContent(data []byte) ([]string) {
	var sLines []string
	if data != nil {
		sData := string(data)
		sLines = strings.Split(sData, "\n")
	}
	return sLines
}

/**
 *	helper method to create a file object
 */
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
 */
func RemoveFile(filename string) error {
	return os.Remove(filename)
}




