package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dstotijn/go-notion"
)

func processMetadata(client *notion.Client, db_id string, user_id string) {
	filename := fmt.Sprintf("%s.json", user_id)

	if Exists(filename) {

	} else {

	}

}

func createMetadata(client *notion.Client, user_db_id string, contribution_db_id string, user_id string) Metadata {

	user_page := getUserPage(client, user_db_id, user_id)
	page_id := user_page.ID

	url := user_page.URL
	desp := "He/She is one of wagumi members."

	ctx := context.Background()
	pq := &notion.PaginationQuery{}

	resp_tmp2, err := client.FindPagePropertyByID(ctx, page_id, "name", pq)
	check(err)
	name := resp_tmp2.Results[0].Title.PlainText

	prop, err := directCallNotionPageProperties(page_id, map_prop_id["icon"])
	check(err)
	image := prop.Files[0].Name

	cntb := getContributionData(client, contribution_db_id, user_id)

	return Metadata{
		Name:         name,
		Description:  desp,
		Image:        image,
		External_url: url,
		Properties: MetadetaProperty{
			Contribusions: *cntb,
		},
	}

}

func getContributionData(client *notion.Client, db_id string, user_id string) *[]Contribution {
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
	check(err)

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

			prop, err := directCallNotionPageProperties(page.ID, map_prop_id["image"])
			check(err)
			contribution.Image = prop.Files[0].File.Url

			prop, err = directCallNotionPageProperties(page.ID, map_prop_id["description"])
			check(err)
			contribution.Description = prop.Results[0].RichText.PlainText

			prop, err = directCallNotionPageProperties(page.ID, map_prop_id["date"])
			check(err)
			contribution.Properties.Date.Start = prop.Date.Start
			contribution.Properties.Date.End = prop.Date.End

			contribusions = append(contribusions, contribution)
		}
	}
	return &contribusions

}

func directCallNotionPageProperties(page_id, property_id string) (PropertyResponse, error) {

	var resStruct PropertyResponse
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s/properties/%s", page_id, property_id)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", api_key)
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

func getUserPage(client *notion.Client, db_id string, user_id string) notion.Page {
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

	return resp.Results[0]
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
