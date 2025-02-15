package utils

import (
	"bytes"
	"io/ioutil"
	"os"
)

func IsStringInFile(s, path string) bool {
	if len(s) <= 0 {
		panic("empty string")
	}
	// use only with small files as we read it fully
	d, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return bytes.Contains(d, []byte(s))
}
