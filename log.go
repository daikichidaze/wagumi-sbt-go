package main

import (
	"errors"
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
		return errors.New("Exection data file exists")
	}

	ini := Log{
		Time:    time.Now(),
		Message: "initialize",
		UserId:  userid,
	}

	err := utils.WriteJsonFile(filename, ini)
	return err

}

func updateExecutionData(filename, message, user_id string) error {
	if !utils.Exists(filename) {
		return errors.New("Execution data file does not exists")
	}

	log := Log{
		Time:    time.Now(),
		Message: message,
		UserId:  user_id,
	}

	err := utils.WriteJsonFile(filename, log)
	return err

}
