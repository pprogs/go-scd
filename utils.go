package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//CheckIfError ee
func CheckIfError(err error) {
	if err != nil {
		lg.Printf("%v\n", err)
		panic(err)
	}
}

//PathExists f
func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

//ReadJSON f
func ReadJSON(filePath string, obj interface{}) (interface{}, error) {

	var data []byte
	var err error
	var f *os.File

	if f, err = os.OpenFile(filePath, os.O_RDWR|os.O_SYNC, 0755); err != nil {
		return nil, err
	}
	defer f.Close()

	if data, err = ioutil.ReadAll(f); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	return obj, nil
}
