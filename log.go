package main

import (
	"errors"
	"sort"
	"time"

	"github.com/daikichidaze/wagumi-sbt-go/utils"
)

type Log struct {
	Time    time.Time `json:"time"`
	Message string    `json:"log"`
	UserId  string    `json:"userid"`
}

func makeExecutionData(filename string) error {

	if utils.Exists(filename) {
		return errors.New("Exection data file exists")
	}

	ini := Log{
		Time:    time.Now(),
		Message: "initialize",
		// UserId:  userid,
	}

	err := utils.WriteJsonFile(filename, ini)
	return err

}

func updateExecutionData(filename, message, user_id string) error {
	if !utils.Exists(filename) {
		return errors.New("Execution data file does not exists")
	}

	if message == "" {
		return errors.New("Message to execution data is null")
	}

	if user_id == "" {
		return errors.New("User ID on  execution data is null")
	}

	log := Log{
		Time:    time.Now(),
		Message: message,
		UserId:  user_id,
	}

	err := utils.WriteJsonFile(filename, log)
	return err

}

func findLastExecutionResult(logs []Log, user_id string) Log {
	sort.Slice(logs, func(i, j int) bool { return logs[i].Time.After(logs[j].Time) }) //Decending

	for _, log := range logs {
		if log.UserId == user_id {
			return log
		}
	}

	return Log{}

}
