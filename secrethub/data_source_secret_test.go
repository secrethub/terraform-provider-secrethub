package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceSecret_PathUnversioned(t *testing.T) {
	config := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			value = "secretpassword"
		}

		data "secrethub_secret" "%v" {
			path = secrethub_secret.%v.path
		}
	`, testAcc.secretName, testAcc.path, testAcc.secretName, testAcc.secretName)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secrethub_secret.%v", testAcc.secretName),
						"path",
						testAcc.path,
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secrethub_secret.%v", testAcc.secretName),
						"value",
						"secretpassword",
					),
				),
			},
		},
	})
}

func TestAccDataSourceSecret_PathVersioned(t *testing.T) {
	configInit := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			value = "secretpasswordv1"
		}

		data "secrethub_secret" "%v" {
			path = "${secrethub_secret.%v.path}:1"
		}
	`, testAcc.secretName, testAcc.path, testAcc.secretName, testAcc.secretName)

	configVersioned := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			value = "secretpasswordv2"
		}

		data "secrethub_secret" "%v" {
			path = "${secrethub_secret.%v.path}:1"
		}
	`, testAcc.secretName, testAcc.path, testAcc.secretName, testAcc.secretName)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: configInit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secrethub_secret.%v", testAcc.secretName),
						"value",
						"secretpasswordv1",
					),
				),
			},
			{
				Config: configVersioned,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secrethub_secret.%v", testAcc.secretName),
						"value",
						"secretpasswordv1",
					),
				),
			},
		},
	})
}
