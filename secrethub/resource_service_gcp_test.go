package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccResourceServiceGCP(t *testing.T) {
	repoPath := testAcc.namespace + "/" + testAcc.repository
	kmsKey := testAcc.gcpKmsKey
	serviceAccount := testAcc.gcpServiceAccount
	description := "TestAccResourceServiceGCP " + acctest.RandString(30)

	config := fmt.Sprintf(`
		resource "secrethub_service_gcp" "test" {
			repo    	          = "%s"
			description	          = "%s"
			service_account_email = "%s"
			kms_key_id	          = "%s"
		}
	`, repoPath, description, serviceAccount, kmsKey)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheckGCP(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkServiceExistsRemotely(repoPath, description),
				),
			},
		},
	})
}
