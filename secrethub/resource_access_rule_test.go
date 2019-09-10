package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceAccessRule_create(t *testing.T) {
	repoPath := testAcc.namespace + "/" + testAcc.repository
	accountName := testAcc.secondAccountName
	permission := "read"

	config := fmt.Sprintf(`
		resource "secrethub_access_rule" "test" {
			dir_path = "%s"
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

func checkAccessRuleExistsRemotely(path string, account string, permission string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := *testAccProvider.Meta().(providerMeta).client

		accessRule, err := client.AccessRules().Get(path, account)
		if err != nil {
			return fmt.Errorf("cannot get created error: %s", err)
		}

		actual := accessRule.Permission.String()
		if actual != permission {
			return fmt.Errorf("expected permission %s but got %s", permission, actual)
		}

		return nil
	}
}
