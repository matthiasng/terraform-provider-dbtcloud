package data_sources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDbtCloudDatabricksCredentialDataSource(t *testing.T) {

	randomProjectName := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	config := databricks_credential(randomProjectName, "moo", "baa", "maa", 64)

	var check resource.TestCheckFunc

	check = resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.dbtcloud_databricks_credential.test", "credential_id"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_databricks_credential.test", "project_id"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_databricks_credential.test", "adapter_id"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_databricks_credential.test", "target_name"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_databricks_credential.test", "schema"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_databricks_credential.test", "num_threads"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_databricks_credential.test", "catalog"),
	)

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  check,
			},
		},
	})
}

func databricks_credential(projectName string, defaultSchema string, username string, password string, numThreads int) string {
	return fmt.Sprintf(`
    resource "dbtcloud_project" "test_credential_project" {
        name = "%s"
    }

	resource "dbtcloud_connection" "databricks" {
		project_id = dbtcloud_project.test_credential_project.id
		type       = "adapter"
		name       = "Databricks"
		database   = ""
		host_name  = "databricks.com"
		http_path  = "/my/path"
		catalog    = "moo"
	  }

    data "dbtcloud_databricks_credential" "test" {
        project_id = dbtcloud_project.test_credential_project.id
        credential_id = dbtcloud_databricks_credential.test_cred.credential_id
    }
    

	resource "dbtcloud_databricks_credential" "test_cred" {
		project_id = dbtcloud_project.test_credential_project.id
		adapter_id = dbtcloud_connection.databricks.adapter_id
		token = "abcdefg"
		schema = "my_schema"
		adapter_type = "databricks"
		catalog = "my_catalog"
	}
	`, projectName)

}
