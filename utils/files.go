package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

func ExportJsonFile[T any](filename string, data T) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	b, err = _UnescapeUnicodeCharactersInJSON(b)
	var out bytes.Buffer
	err = json.Indent(&out, b, "", strings.Repeat(" ", 2))
	if err != nil {
		return err
	}

	_, err = f.Write(out.Bytes())

	return nil

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

func _UnescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
