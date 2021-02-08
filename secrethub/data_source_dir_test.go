package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceDir(t *testing.T) {
	cases := map[string]struct {
		config string
	}{
		"repo root directory": {
			config: fmt.Sprintf(`
				data "secrethub_dir" "repo" {
					path = "%v"
				}
				`, testAcc.repoPath),
		},
		"subdirectory": {
			config: fmt.Sprintf(`
				data "secrethub_dir" "repo" {
					path = "%v"
				}

				resource "secrethub_dir" "subdir" {
					path = "${data.secrethub_dir.repo.path}/subdir"
				}

				data "secrethub_dir" "subdir" {
					path = secrethub_dir.subdir.path
				}
				`, testAcc.repoPath),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				Providers: testAccProviders,
				PreCheck:  testAccPreCheck(t),
				Steps: []resource.TestStep{
					{
						Config: tc.config,
					},
				},
			})
		})
	}
}
