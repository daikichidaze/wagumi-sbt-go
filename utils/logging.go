package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Log struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

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
