package main

import "bytes"

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
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Image       string               `json:"image"`
	ExternalUrl string               `json:"external_url"`
	Properties  ContributionProperty `json:"properties"`
	users       []string
}

type ContributionProperty struct {
	PageId    string   `json:"page_id"`
	Reference []string `json:"reference"`
	Date      Date     `json:"date"`
}

type Date struct {
	Start string  `json:"start"`
	End   DateEnd `json:"end"` //endが存在していた場合、文字列として認識してそうでない場合nullを返したい

}

type DateEnd string

func (c DateEnd) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if len(string(c)) == 0 {
		buf.WriteString(`null`)
	} else {
		buf.WriteString(`"` + string(c) + `"`) // add double quation mark as json format required
	}
	return buf.Bytes(), nil
}

func (c *DateEnd) UnmarshalJSON(in []byte) error {
	str := string(in)
	if str == `null` {
		*c = ""
		return nil
	}
	res := DateEnd(str)
	if len(res) >= 2 {
		res = res[1 : len(res)-1] // remove the wrapped qutation
	}
	*c = res
	return nil
}
