package main

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


