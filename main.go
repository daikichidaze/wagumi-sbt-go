package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"github.com/daikichidaze/wagumi-sbt-go/utils"
)

var env_file = ".env"
var log_file_name = "executionData.json"
var api_key string
var page_contribution_map map[string]Contribution

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
			Required: false,
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
		// user_id := c.String("id")
		api_key = c.String("api_key")
		user_db_id := c.String("user_db_id")
		contribute_db_id := c.String("contribute_db_id")

		client := notion.NewClient(api_key)

		if !utils.Exists(log_file_name) { // First execution
			fmt.Println("This is first execution. Creating " + log_file_name)
			makeExecutionData(log_file_name)
		} else { // From second execution
			fmt.Println("Load " + log_file_name)

		}

		logs, err := utils.ReadJsonFile[[]Log](log_file_name)
		utils.Check(err)
		last_exe_log := getLastExecutionResultsMap(logs)

		err = processMetadata(client, user_db_id, contribute_db_id, last_exe_log)
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

func processMetadata(client *notion.Client,
	user_db_id string, contribution_db_id string, last_exe_log Log) error {

	ctx := context.Background()
	pq := &notion.PaginationQuery{}
	execution_timestamp := time.Now()

	var filer_timestamp time.Time
	if last_exe_log.Message != "initialize" {
		filer_timestamp = last_exe_log.Time
	}

	query := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "last_edited_time",
			Date: &notion.DateDatabaseQueryFilter{
				After: &filer_timestamp,
			},
		},
		Sorts: []notion.DatabaseQuerySort{
			{
				Property:  "last_edited_time",
				Timestamp: notion.SortTimeStampLastEditedTime,
				Direction: notion.SortDirAsc,
			},
		},
	}

	resp, err := client.QueryDatabase(ctx, contribution_db_id, query)
	if err != nil {
		return err
	}

	if len(resp.Results) == 0 {
		err = updateExecutionData(log_file_name, "no updates", "", execution_timestamp)
		if err != nil {
			return err
		}

		fmt.Println("No updates")
		return nil
	}

	page_contribution_map = make(map[string]Contribution)
	for _, page := range resp.Results {
		page_contribution_map[page.ID] = createContribution(client, pq, page, ctx)
	}

	user_contribution_map := makeUserPageidMap(client, page_contribution_map)

	for key, value := range user_contribution_map {
		md := createSingleUserMetadataFromMap(client, user_db_id, key, value)
		msg, err := postProcessing(md, last_exe_log, md.filename)
		if err != nil {
			return err
		}

		err = updateExecutionData(log_file_name, msg, key, execution_timestamp)
		if err != nil {
			return err
		}

	}

	fmt.Println(log_file_name + " updated")
	return nil
}
