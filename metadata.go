package main

type Metadata struct {
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Image        string           `json:"image"`
	External_url string           `json:"external_url"`
	Properties   MetadetaProperty `json:"properties"`
	id           string
	filename     string
}

type MetadetaProperty struct {
	// Sns           map[string]string   `json:"sns"`
	// Tokens        []map[string]string `json:"tokens"`
	Contribusions []Contribution `json:"contributions"`
}

type Contribution struct {
	PageId    	string   			 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Image       string               `json:"image"`
	ExternalUrl string               `json:"external_url"`
	Date      	Date     			 `json:"date"`
	users       []string
}

type Date struct {
	Start string `json:"start"`
	End   interface{} `json:"end"`
	// TimeZone interface{} `json:"time_zone"`
}
