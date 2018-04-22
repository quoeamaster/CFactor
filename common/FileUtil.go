package common

import (
	"io/ioutil"
	"strings"
)

/**
 *	method to simply load a file from the given "name"
 */
func LoadFile(name string) ([]byte, error) {
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




