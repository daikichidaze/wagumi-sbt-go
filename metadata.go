package main

type Metadata struct {
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Image        string           `json:"image"`
	External_url string           `json:"external_url"`
	Properties   MetadetaProperty `json:"properties"`
}

type MetadetaProperty struct {
	// Sns           map[string]string   `json:"sns"`
	// Tokens        []map[string]string `json:"tokens"`
	Contribusions []Contribution `json:"contributions"`
}

type Contribution struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Image       string               `json:"image"`
	ExternalUrl string               `json:"external_url"`
	Properties  ContributionProperty `json:"properties"`
}

type ContributionProperty struct {
	PageId    string   `json:"page_id"`
	Reference []string `json:"reference"`
	Date      struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"date"`
}
