package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/daikichidaze/wagumi-sbt-go/utils"
)

type Log struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	UserId  string    `json:"userid"`
}

func makeExecutionData(filename, userid string) error {

	if utils.Exists(filename) {
		return errors.New("Exection file exists")
	}

	ini := Log{
		Time:    time.Now(),
		Message: "initialize",
		UserId:  userid,
	}

	err := utils.WriteJsonFile(filename, ini)
	return err

}

func readLastExecution(filename string) ([]Log, error) {
	var logs []Log

	if !utils.Exists(filename) {
		return logs, errors.New("Execution file does not exists")
	}

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return logs, err
	}

	json.Unmarshal(raw, &logs)
	return logs, nil

}
