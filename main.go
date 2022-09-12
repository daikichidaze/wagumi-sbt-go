package main

import (
	"context"
	"errors"
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
var metadata_directory = "metadata"
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
			EnvVars: []string{"WAGUMI_SAMURAI_API_TOKEN"},
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

		start := time.Now()
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

		end := time.Now()
		fmt.Println("Execution time:", end.Sub(start))

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

	if !utils.Exists(metadata_directory) {
		return errors.New(fmt.Sprintf("Directory: \"%s\" does not exists. Please create the directory for metadata.", metadata_directory))
	}

	ctx := context.Background()
	pq := &notion.PaginationQuery{}
	execution_timestamp := time.Now()

	if last_exe_log.Message != "initialize" {
	}
  checkbox := new(bool)
  *checkbox = true
	query := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "publish",
			Checkbox: &notion.CheckboxDatabaseQueryFilter{
				Equals: checkbox,
			},
		},
		Sorts: []notion.DatabaseQuerySort{
			{
				Property:  "date",
				Direction: notion.SortDirDesc,
			},
		},
	}

	fmt.Println("Call Notion API")
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

	fmt.Println("Metadata processing...")
	page_contribution_map = make(map[string]Contribution)
	for _, page := range resp.Results {
		page_contribution_map[page.ID] = createContribution(client, pq, page, ctx)
	}

	user_contribution_map := makeUserPageidMap(client, page_contribution_map)

	fmt.Println("Export json files...")
	for key, value := range user_contribution_map {

		md := createSingleUserMetadataFromMap(client, user_db_id, key, value)
		msg, err := postProcessingMetadata(md, last_exe_log, md.filename)
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
