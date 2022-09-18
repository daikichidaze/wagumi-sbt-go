package main

type PropertyResponse struct {
	Object  string `json:"object"`
	Type    string `json:"type"`
	Id      string `json:"id"`
	Results []struct {
		RichText struct {
			Type      string `json:"type"`
			PlainText string `json:"plain_text"`
		} `json:"rich_text"`
	} `json:"results"`
	Files []struct {
		Name string `json:"name"`
		Type string `json:"type"`
		File struct {
			Url        string `json:"url"`
			ExpiryTime string `json:"expiry_time"`
		} `json:"file,omitempty"`
		External struct {
			Url		   string `json:"url"`
		} `json:"external,omitempty"`
	} `json:"files"`
	Date struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"date"`
}
	