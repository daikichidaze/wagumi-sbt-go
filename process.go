package main


func updateMetadata(client *notion.Client, db_id string, user_id string) []Contribution{
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
