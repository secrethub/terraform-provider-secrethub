package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceSecret_absPath(t *testing.T) {
	config := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			data = "secretpassword"
		}

		data "secrethub_secret" "%v" {
			path = "${secrethub_secret.%v.path}"
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
						string(testAcc.path),
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secrethub_secret.%v", testAcc.secretName),
						"data",
						"secretpassword",
					),
				),
			},
		},
	})
}

func TestAccDataSourceSecret_absPathVersioned(t *testing.T) {
	configInit := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			data = "secretpasswordv1"
		}

		data "secrethub_secret" "%v" {
			path = "${secrethub_secret.%v.path}:1"
		}
	`, testAcc.secretName, testAcc.path, testAcc.secretName, testAcc.secretName)

	configVersioned := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			data = "secretpasswordv2"
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
						"data",
						"secretpasswordv1",
					),
				),
			},
			{
				Config: configVersioned,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secrethub_secret.%v", testAcc.secretName),
						"data",
						"secretpasswordv1",
					),
				),
			},
		},
	})
}

func TestAccDataSourceSecret_prefPath(t *testing.T) {
	config := fmt.Sprintf(`
		provider "secrethub" {
			path_prefix = "%v/%v"
		}

		resource "secrethub_secret" "%v" {
			path = "%v"
			data = "secretpassword"
		}

		data "secrethub_secret" "%v" {
			path = "${secrethub_secret.%v.path}"
		}
	`, testAcc.namespace, testAcc.repository, testAcc.secretName, testAcc.secretName, testAcc.secretName, testAcc.secretName)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secrethub_secret.%v", testAcc.secretName),
						"data",
						"secretpassword",
					),
				),
			},
		},
	})
}
