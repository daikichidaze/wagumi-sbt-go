package main

import (
	"fmt"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var env_file = ".env"
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

		createMetadata(client, user_db_id, contribute_db_id, user_id)

		return nil

	}

	app.Run(os.Args)

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loadEnv(file_path string) {
	err := godotenv.Load(file_path)

	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
		panic(err)
	}
}
