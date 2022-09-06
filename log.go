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

	ini := []Log{
		{
			Time:    time.Now(),
			Message: "initialize",
			// UserId:  userid,

		},
	}

	err := utils.ExportJsonFile(filename, ini)
	return err

}

func updateExecutionData(filename, message, user_id string, log_time time.Time) error {
	if !utils.Exists(filename) {
		return errors.New("Execution data file does not exists")
	}

	if message == "" {
		return errors.New("Message to execution data is null")
	}

	last, err := utils.ReadJsonFile[[]Log](filename)

	log := Log{
		Time:    log_time,
		Message: message,
		UserId:  user_id,
	}

	last = append(last, log)

	err = utils.ExportJsonFile(filename, last)
	return err

}

func getLastExecutionResultsMap(logs []Log) Log {
	// sort.Slice(logs, func(i, j int) bool { return logs[i].Time.After(logs[j].Time) }) //Decending

	var last_exe_time time.Time
	var result *Log
	for _, log := range logs {
		if log.Time.After(last_exe_time) {
			last_exe_time = log.Time
			result = &log
		}
	}

	return *result

}

func findLastExecutionResultByUserId(logs []Log, user_id string) Log {
	sort.Slice(logs, func(i, j int) bool { return logs[i].Time.After(logs[j].Time) }) //Decending

	for _, log := range logs {
		if log.UserId == user_id {
			return log
		}
	}

	return Log{}

}
