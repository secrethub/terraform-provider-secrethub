package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceAccessRule(t *testing.T) {
	repoPath := testAcc.namespace + "/" + testAcc.repository
	accountName := testAcc.secondAccountName
	permission := "read"

	config := fmt.Sprintf(`
		resource "secrethub_access_rule" "test" {
			dir = "%s"
			account_name = "%s"
			permission = "%s"
		}
	`, repoPath, accountName, permission)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkAccessRuleExistsRemotely(repoPath, accountName, permission),
				),
			},
		},
	})
}

func TestAccAccessRuleForService(t *testing.T) {
	repoPath := testAcc.namespace + "/" + testAcc.repository
	serviceDescription := "TestAccessRuleForService " + acctest.RandString(30)
	permission := "read"

	config := fmt.Sprintf(`
		resource "secrethub_service" "test" {
			repo = "%s"
			description = "%s"
		}

		resource "secrethub_access_rule" "test" {
			dir = "%s"
			account_name = "${secrethub_service.test.id}"
			permission = "%s"
		}
		`,
		repoPath, serviceDescription, repoPath, permission)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkAccessRuleForServiceExistsRemotely(repoPath, serviceDescription, permission),
				),
			},
		},
	})
}

func checkAccessRuleExistsRemotely(path string, account string, permission string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := *testAccProvider.Meta().(providerMeta).client

		accessRule, err := client.AccessRules().Get(path, account)
		if err != nil {
			return fmt.Errorf("cannot get created access rule: %s", err)
		}

		actual := accessRule.Permission.String()
		if actual != permission {
			return fmt.Errorf("expected permission %s but got %s", permission, actual)
		}

		return nil
	}
}

func checkAccessRuleForServiceExistsRemotely(repoPath string, serviceDescription string, permission string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := *testAccProvider.Meta().(providerMeta).client

		services, err := client.Services().List(repoPath)
		if err != nil {
			return fmt.Errorf("cannot list services: %s", err)
		}

		for _, service := range services {
			if service.Description == serviceDescription {
				return checkAccessRuleExistsRemotely(repoPath, service.ServiceID, permission)(s)
			}
		}

		return fmt.Errorf("cannot find service in repo %s with description %s", repoPath, serviceDescription)
	}
}

