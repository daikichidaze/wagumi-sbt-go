package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"github.com/daikichidaze/wagumi-sbt-go/utils"
)

var env_file = ".env"
var log_file_name = "executionData.json"
var api_key string

var map_prop_id = map[string]string{
	"name":        "title",
	"userId":      "%7BExC",
	"image":       "F%5C%7Dc",
	"description": "NGz%5D",
	"date":        "PG%3Af",
	"icon":        "a%3EG%5B",
}

func main() {
	loadEnv(env_file)

	app := cli.NewApp()
	app.Name = "Wagumi SBT Go client"
	app.Usage = "Generate json metadata from Notion"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "id, i",
			Value:    "",
			Usage:    "Target user id",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "api_key",
			Value:   "",
			Usage:   "Notion API token",
			EnvVars: []string{"NOTION_API_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "contribute_db_id",
			Value:   "",
			Usage:   "Contribute DB's id",
			EnvVars: []string{"WAGUMI_DATABASE_ID"},
		},
		&cli.StringFlag{
			Name:    "user_db_id",
			Value:   "",
			Usage:   "User DB's id",
			EnvVars: []string{"WAGUMI_USER_DATABASE_ID"},
		},
	}
	app.Authors = []*cli.Author{
		{
			Name:  "daikichi",
			Email: "dkch.yoshioka.t@gmail.com",
		},
	}
	app.Action = func(c *cli.Context) error {
		user_id := c.String("id")
		api_key = c.String("api_key")
		user_db_id := c.String("user_db_id")
		contribute_db_id := c.String("contribute_db_id")

		client := notion.NewClient(api_key)

		var last_exe_date time.Time
		var last_exe_json_name string
		if !utils.Exists(log_file_name) { // First execution
			makeExecutionData(log_file_name, user_id)
		} else { // From second execution
			logs, _ := utils.ReadMetadata[[]Log](log_file_name)
			sort.Slice(logs, func(i, j int) bool { return logs[i].Time.Before(logs[j].Time) })
			last_exe_log := logs[len(logs)-1]
			last_exe_date = last_exe_log.Time

			last_exe_json_name = fmt.Sprintf("%s.json", last_exe_log.UserId)

		}

		metadata := createMetadata(client, user_db_id, contribute_db_id, user_id, last_exe_date)

		var message string
		if len(metadata.Properties.Contribusions) > 0 {
			// Only update the metadata when there are new contribusion data
			if last_exe_json_name != "" {
				// Add previouse contribution data
				last_metadata, err := utils.ReadMetadata[Metadata](last_exe_json_name)
				utils.Check(err)

				metadata.Properties.Contribusions =
					append(last_metadata.Properties.Contribusions, metadata.Properties.Contribusions...)
			}

			// export metadata json
			err := exportMetadataJsonFile(fmt.Sprintf("%s.json", user_id), metadata)
			utils.Check(err)

			message = fmt.Sprintf("update %s", log_file_name)
		} else {
			message = fmt.Sprintf("no updates in %s", log_file_name)
		}

		err := updateExecutionData(log_file_name, message, user_id)
		utils.Check(err)

		return nil

	}

	app.Run(os.Args)

}

func loadEnv(file_path string) {
	err := godotenv.Load(file_path)

	if err != nil {
		fmt.Printf("Faild to read: %v", err)
		panic(err)
	}
}
