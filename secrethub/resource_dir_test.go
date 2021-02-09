package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceDir(t *testing.T) {
	config := fmt.Sprintf(`
		resource "secrethub_dir" "%v" {
			path = "%v"
		}
	`, testAcc.dirName, testAcc.dirPath)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkDirExistsRemotely(testAcc.dirPath),
				),
			},
			{
				ResourceName:            fmt.Sprintf("secrethub_dir.%v", testAcc.dirName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func checkDirExistsRemotely(path string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := *testAccProvider.Meta().(providerMeta).client

		exists, err := client.Dirs().Exists(path)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("expected directory to exist: %s", path)
		}

		return nil
	}
}
