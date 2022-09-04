package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dstotijn/go-notion"

	"github.com/daikichidaze/wagumi-sbt-go/utils"
)

func postProcessing(metadata Metadata, last_exe_log Log, metadata_file_name string) (string, error) {
	var message string

	// Only update the metadata when there are new contribusion data

	if len(metadata.Properties.Contribusions) > 0 {
		if utils.Exists(metadata_file_name) {
			// Add previouse contribution data
			last_metadata, err := utils.ReadJsonFile[Metadata](metadata_file_name)
			if err != nil {
				return "", err
			}

			metadata.Properties.Contribusions =
				append(last_metadata.Properties.Contribusions, metadata.Properties.Contribusions...)

			message = fmt.Sprintf("update metadata/%s", metadata_file_name)

		} else {
			message = fmt.Sprintf("create metadata/%s", metadata_file_name)
		}

		// export metadata json
		err := exportMetadataJsonFile(metadata_file_name, metadata)
		if err != nil {
			return "", err
		}
		fmt.Println(metadata_file_name + " updated")
	} else {
		// case of no new contributions
		message = fmt.Sprintf("no updates in %s", metadata_file_name)
		fmt.Println(message)
	}
	return message, nil

}

func makeUserPageidMap(client *notion.Client, page_map map[string]Contribution) map[string]([]string) {

	result := make(map[string]([]string))

	for pageid, contribusion := range page_map {

		for _, userid := range contribusion.users {
			_, ok := result[userid]

			if !ok {
				result[userid] = make([]string, 0)
			}

			result[userid] = append(result[userid], pageid)
		}
	}
	return result
}

func createSingleUserMetadataFromMap(client *notion.Client,
	user_db_id string, user_id string, pageIds []string) Metadata {

	user_page := getUserPage(client, user_db_id, user_id)
	page_id := user_page.ID

	url := getNotionExternalURL(user_page.URL)
	desp := "He/She is one of wagumi members."

	ctx := context.Background()
	pq := &notion.PaginationQuery{}

	resp_tmp2, err := client.FindPagePropertyByID(ctx, page_id, "name", pq)
	utils.Check(err)
	name := resp_tmp2.Results[0].Title.PlainText

	prop, err := directCallNotionPageProperties(page_id, map_prop_id["icon"])
	utils.Check(err)
	image := prop.Files[0].Name

	cntb := make([]Contribution, 0)
	for _, pageid := range pageIds {
		cntb = append(cntb, page_contribution_map[pageid])
	}

	return Metadata{
		Name:         name,
		Description:  desp,
		Image:        image,
		External_url: url,
		Properties: MetadetaProperty{
			Contribusions: cntb,
		},
		id:       user_id,
		filename: fmt.Sprintf("%s.json", user_id),
	}
}

func createContribution(client *notion.Client,
	pagination *notion.PaginationQuery, page notion.Page, ctx context.Context) Contribution {
	page_id := page.ID
	external_url := getNotionExternalURL(page.URL)

	resp_tmp, err := client.FindPagePropertyByID(ctx, page_id, "name", pagination)
	utils.Check(err)
	name := resp_tmp.Results[0].Title.PlainText

	prop, err := directCallNotionPageProperties(page_id, map_prop_id["image"])
	utils.Check(err)
	image := prop.Files[0].File.Url

	prop, err = directCallNotionPageProperties(page_id, map_prop_id["description"])
	utils.Check(err)
	description := prop.Results[0].RichText.PlainText

	prop, err = directCallNotionPageProperties(page_id, map_prop_id["date"])
	utils.Check(err)
	start := prop.Date.Start
	end := prop.Date.End

	resp_users, err := client.FindPagePropertyByID(ctx, page.ID, map_prop_id["userId"], pagination)
	utils.Check(err)
	users := make([]string, 0)
	for _, user := range resp_users.Results {
		userid_string := user.RichText.PlainText
		users = append(users, userid_string)
	}

	return Contribution{
		Name:        name,
		Description: description,
		Image:       image,
		ExternalUrl: external_url,
		Properties: ContributionProperty{
			PageId: page_id,
			Date: Date{
				Start: start,
				End:   end,
			},
		},
		users: users,
	}

}

func getContributionDataFromMap(client *notion.Client, user_pages []notion.Page) *[]Contribution {
	ctx := context.Background()
	pq := &notion.PaginationQuery{}

	contributions := make([]Contribution, 0)

	for _, page := range user_pages {

		page_id := page.ID
		external_url := getNotionExternalURL(page.URL)

		resp_tmp, err := client.FindPagePropertyByID(ctx, page_id, "name", pq)
		utils.Check(err)
		name := resp_tmp.Results[0].Title.PlainText

		prop, err := directCallNotionPageProperties(page_id, map_prop_id["image"])
		utils.Check(err)
		image := prop.Files[0].File.Url

		prop, err = directCallNotionPageProperties(page_id, map_prop_id["description"])
		utils.Check(err)
		description := prop.Results[0].RichText.PlainText

		prop, err = directCallNotionPageProperties(page_id, map_prop_id["date"])
		utils.Check(err)
		start := prop.Date.Start
		end := prop.Date.End

		contributions = append(contributions,
			Contribution{
				Name:        name,
				Description: description,
				Image:       image,
				ExternalUrl: external_url,
				Properties: ContributionProperty{
					PageId: page_id,
					Date: Date{
						Start: start,
						End:   end,
					},
				},
			})

	}
	return &contributions

}

// Process for one user id execution

func createSingleUserMetadataFromDB(client *notion.Client,
	user_db_id string, contribution_db_id string, user_id string,
	last_execution_time time.Time) Metadata {

	user_page := getUserPage(client, user_db_id, user_id)
	page_id := user_page.ID

	url := getNotionExternalURL(user_page.URL)
	desp := "He/She is one of wagumi members."

	ctx := context.Background()
	pq := &notion.PaginationQuery{}

	resp_tmp2, err := client.FindPagePropertyByID(ctx, page_id, "name", pq)
	utils.Check(err)
	name := resp_tmp2.Results[0].Title.PlainText

	prop, err := directCallNotionPageProperties(page_id, map_prop_id["icon"])
	utils.Check(err)
	image := prop.Files[0].Name

	cntb := getSingleUserContributionDataFromDB(client, contribution_db_id, user_id, last_execution_time)

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

func getSingleUserContributionDataFromDB(client *notion.Client, db_id string, user_id string, last_execution_time time.Time) *[]Contribution {
	ctx := context.Background()
	pq := &notion.PaginationQuery{}

	query := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "last_edited_time",
			Date: &notion.DateDatabaseQueryFilter{
				After: &last_execution_time,
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

	resp, err := client.QueryDatabase(ctx, db_id, query)
	utils.Check(err)

	contribusions := make([]Contribution, 0)

	for _, page := range resp.Results {

		resp_users, _ := client.FindPagePropertyByID(ctx, page.ID, map_prop_id["userId"], pq)

		var userSearchResult bool

		for _, user := range resp_users.Results {
			if user.RichText.PlainText == user_id {
				userSearchResult = true
			}
		}

		if userSearchResult {
			page_id := page.ID
			external_url := getNotionExternalURL(page.URL)

			resp_tmp, err := client.FindPagePropertyByID(ctx, page.ID, "name", pq)
			utils.Check(err)
			name := resp_tmp.Results[0].Title.PlainText

			prop, err := directCallNotionPageProperties(page.ID, map_prop_id["image"])
			utils.Check(err)
			image := prop.Files[0].File.Url

			prop, err = directCallNotionPageProperties(page.ID, map_prop_id["description"])
			utils.Check(err)
			description := prop.Results[0].RichText.PlainText

			prop, err = directCallNotionPageProperties(page.ID, map_prop_id["date"])
			utils.Check(err)
			start := prop.Date.Start
			end := prop.Date.End

			contribusions = append(contribusions,
				Contribution{
					Name:        name,
					Description: description,
					Image:       image,
					ExternalUrl: external_url,
					Properties: ContributionProperty{
						PageId: page_id,
						Date: Date{
							Start: start,
							End:   end,
						},
					},
				})
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

func getNotionExternalURL(internal_url string) string {
	external_base_url := "https://wagumi-dev.notion.site/"

	split_result1 := strings.Split(internal_url, "/")
	split_result2 := strings.Split(split_result1[len(split_result1)-1], "-")

	result, err := utils.UrlJoin(external_base_url, split_result2[len(split_result2)-1])
	utils.Check(err)
	return result

}

func getMetadataFilename(user_id string) string {
	return fmt.Sprintf("%s.json", user_id)
}

func exportMetadataJsonFile(filename string, data Metadata) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	err = enc.Encode(data)

	if err != nil {
		return err
	}

	return nil

}
