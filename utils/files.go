package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func ExportJsonFile[T any](filename string, data T) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	err = enc.Encode(data)

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

func _UnescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
