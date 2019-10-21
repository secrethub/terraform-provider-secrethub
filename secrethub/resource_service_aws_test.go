package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccResourceServiceAWS(t *testing.T) {
	repoPath := testAcc.namespace + "/" + testAcc.repository
	kmsKey := testAcc.awsKmsKey
	role := testAcc.awsRole
	description := "TestAccResourceServiceAWS " + acctest.RandString(30)

	config := fmt.Sprintf(`
		resource "secrethub_service_aws" "test" {
			repo    	= "%s"
			kms_key_arn	= "%s"
			role    	= "%s"
			description	= "%s"
		}
	`, repoPath, kmsKey, role, description)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
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
