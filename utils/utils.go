package utils

import (
	"net/url"
	"os"
	"path"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func UrlJoin(rawurl string, elem ...string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(append([]string{u.Path}, elem...)...)
	return u.String(), nil
}
