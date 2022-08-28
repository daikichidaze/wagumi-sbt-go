package main

import (
	"os"
	"testing"

	"github.com/dstotijn/go-notion"
)

func TestCreateNewMetadata(t *testing.T) {
	loadEnv(env_file)
	api_key = os.Getenv("NOTION_API_TOKEN")
	contribution_db_id := os.Getenv("WAGUMI_DATABASE_ID")
	user_db_id := os.Getenv("WAGUMI_USER_DATABASE_ID")
	test_user_id := os.Getenv("TEST_USER_ID")

	client := notion.NewClient(api_key)

	createMetadata(client, user_db_id, contribution_db_id, test_user_id)

}
