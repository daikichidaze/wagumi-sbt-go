package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var env_file = ".env"
var api_key2 string

var map_prop_id = map[string]string{
	"name":        "title",
	"userId":      "%7BExC",
	"image":       "F%5C%7Dc",
	"description": "NGz%5D",
	"date":        "PG%3Af",
}

type Metadata struct {
	Name         string
	Description  string
	Image        string
	External_url string
	Properties   MetadetaProperty
}

type MetadetaProperty struct {
	// Sns           map[string]string
	// Tokens        []map[string]string
	Contribusions []Contribution
}

type Contribution struct {
	Name        string
	Description string
	Image       string
	ExternalUrl string
	Properties  ContributionProperty
}

type ContributionProperty struct {
	PageId    string
	Reference []string
	Date      struct {
		Start string
		End   string
	}
}

func main() {
	loadEnv(env_file)

	app := cli.NewApp()
	app.Name = "Wagumi SBT Go client"
	app.Usage = "Generate json metadata from Notion"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "id, i",
			Value: "",
			Usage: "Target user id",
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
		api_key := c.String("api_key")
		contribute_db_id := c.String("contribute_db_id")

		client := notion.NewClient(api_key)
		api_key2 = api_key

		updateMetadata(client, contribute_db_id, user_id)

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

func getUserPageID(client *notion.Client, db_id string, user_id string) string {
	query := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "id",
			Text: &notion.TextDatabaseQueryFilter{
				Equals: user_id,
			},
		},
	}
	resp, err := client.QueryDatabase(context.Background(), db_id, query)

	if err != nil {
		panic(err)
	}

	if len(resp.Results) > 1 {
		panic("More than one user results")
	}

	return resp.Results[0].ID
}

func updateMetadata(client *notion.Client, db_id string, user_id string) {
	ctx := context.Background()
	pq := &notion.PaginationQuery{}

	query := &notion.DatabaseQuery{
		// Filter: &notion.DatabaseQueryFilter{
		// 	Property: "userid",
		// 	Relation: &notion.RelationDatabaseQueryFilter{
		// 		Contains: user_id,
		// 	},
		// },
		Sorts: []notion.DatabaseQuerySort{
			{
				Property:  "last_edited_time",
				Timestamp: notion.SortTimeStampLastEditedTime,
				Direction: notion.SortDirAsc,
			},
		},
	}

	resp, err := client.QueryDatabase(ctx, db_id, query)

	if err != nil {
		panic(err)
	}

	var contribusions []Contribution

	for _, page := range resp.Results {

		resp_users, _ := client.FindPagePropertyByID(ctx, page.ID, map_prop_id["userId"], pq)

		var userSearchResult bool

		for _, user := range resp_users.Results {
			if user.RichText.PlainText == user_id {
				userSearchResult = true
			}
		}

		if userSearchResult {
			var contribution Contribution

			contribution.ExternalUrl = page.URL
			contribution.Properties.PageId = page.ID

			resp_tmp, err := client.FindPagePropertyByID(ctx, page.ID, "name", pq)
			check(err)
			contribution.Name = resp_tmp.Results[0].Title.PlainText

			prop, err := direct_call(page.ID, map_prop_id["image"])
			check(err)
			contribution.Image = prop.Files[0].File.Url

			prop, err = direct_call(page.ID, map_prop_id["description"])
			check(err)
			contribution.Description = prop.Results[0].RichText.PlainText

			prop, err = direct_call(page.ID, map_prop_id["date"])
			check(err)
			contribution.Properties.Date.Start = prop.Date.Start
			contribution.Properties.Date.End = prop.Date.End

			contribusions = append(contribusions, contribution)
		}
	}
	return

}

func direct_call(page_id, property_id string) (PropertyResponse, error) {

	var resStruct PropertyResponse
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s/properties/%s", page_id, property_id)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", api_key2)
	req.Header.Add("Notion-Version", "2022-06-28")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return resStruct, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resStruct, err
	}

	err = json.Unmarshal(body, &resStruct)
	if err != nil {
		return resStruct, err
	}

	return resStruct, err

}
