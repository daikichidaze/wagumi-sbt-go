package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dstotijn/go-notion"
)

func TestCreateNewMetadata(t *testing.T) {
	loadEnv(env_file)
	api_key = os.Getenv("NOTION_API_TOKEN")
	contribution_db_id := os.Getenv("WAGUMI_DATABASE_ID")
	user_db_id := os.Getenv("WAGUMI_USER_DATABASE_ID")
	test_user_id := os.Getenv("TEST_USER_ID")

	client := notion.NewClient(api_key)
	ti := time.Date(2022, time.August, 1, 0, 0, 0, 0, time.UTC)

	md := createMetadata(client, user_db_id, contribution_db_id, test_user_id, ti)

	err := exportMetadataJsonFile(fmt.Sprintf("%s.json", test_user_id), md)
	//TODO finalize test code
	if err != nil {
		t.Error("Faild for json output")
	}
	t.Log("Finished TestCreateNewMetadata")

}

func TestGetNotionExternalURL(t *testing.T) {
	internal_url := "https://www.notion.so/daikichi-9f2134fdd56246859220a0d551174783"
	output_url := getNotionExternalURL(internal_url)

	expected_result := "https://wagumi-dev.notion.site/9f2134fdd56246859220a0d551174783"

	if output_url != expected_result {
		t.Error("Notion url conversion error")
	}
	t.Log("Finished TestGetNotionExternalURL")

}

func TestGetContributionData(t *testing.T) {

	loadEnv(env_file)
	api_key = os.Getenv("NOTION_API_TOKEN")
	contribution_db_id := os.Getenv("WAGUMI_DATABASE_ID")
	// user_db_id := os.Getenv("WAGUMI_USER_DATABASE_ID")
	test_user_id := os.Getenv("TEST_USER_ID")

	client := notion.NewClient(api_key)

	ti := time.Date(2022, time.August, 1, 0, 0, 0, 0, time.UTC)
	cntr := getContributionData(client, contribution_db_id, test_user_id, ti)
	fmt.Println(*cntr)

}
