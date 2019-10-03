package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceService_create(t *testing.T) {
	repoPath := testAcc.namespace + "/" + testAcc.repository
	serviceDescription := "TestAccResourceService_create " + acctest.RandString(30)

	config := fmt.Sprintf(`
		resource "secrethub_service" "test" {
			repo = "%s"
			description = "%s"
		}
	`, repoPath, serviceDescription)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkServiceExistsRemotely(repoPath, serviceDescription),
				),
			},
		},
	})
}

func checkServiceExistsRemotely(path string, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := *testAccProvider.Meta().(providerMeta).client

		services, err := client.Services().List(path)
		if err != nil {
			return fmt.Errorf("cannot list services: %s", err)
		}

		for _, service := range services {
			if service.Description == description {
				return nil
			}
		}

		return fmt.Errorf("expected service on repo %s with description \"%s\"", path, description)
	}
}
