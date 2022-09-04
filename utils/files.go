package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func WriteJsonFile(fileName string, object interface{}) error {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0600)
	defer file.Close()
	fi, _ := file.Stat()
	leng := fi.Size()

	by := make([]byte, leng)
	file.Read(by)

	json_, _ := json.Marshal(object)

	var err error
	if leng == 0 {
		_, err = file.Write([]byte(fmt.Sprintf(`[%s]`, json_)))
	} else if by[leng-1] == 0xa {
		_, err = file.WriteAt([]byte(fmt.Sprintf(`,%s]`, json_)), leng-2)
	} else {
		_, err = file.WriteAt([]byte(fmt.Sprintf(`,%s]`, json_)), leng-1)
	}
	return err
}

func ReadJsonFile[T any](filename string) (T, error) {
	var data T

	if !Exists(filename) {
		return data, errors.New("Target file does not exists: " + filename)
	}

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return data, err
	}

	json.Unmarshal(raw, &data)
	return data, nil

}
