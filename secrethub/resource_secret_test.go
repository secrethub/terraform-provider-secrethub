package secrethub

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceSecret_writeAbsPath(t *testing.T) {
	config := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			data = "secretpassword"
		}
	`, testAcc.secretName, testAcc.path)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(testAcc),
				),
			},
		},
	})
}

func TestAccResourceSecret_writePrefPath(t *testing.T) {
	config := fmt.Sprintf(`
		provider "secrethub" {
			path_prefix = "%v"
		}

		resource "secrethub_secret" "%v" {
			path = "%v/%v"
			data = "secretpassword"
		}
	`, testAcc.namespace, testAcc.secretName, testAcc.repository, testAcc.secretName)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(testAcc),
				),
			},
		},
	})
}

func TestAccResourceSecret_writePrefPathOverride(t *testing.T) {
	config := fmt.Sprintf(`
		provider "secrethub" {
			path_prefix = "override_me"
		}
		
		resource "secrethub_secret" "%v" {
			path_prefix = "%v"
			path = "%v/%v"
			data = "secretpassword"
		}
	`, testAcc.secretName, testAcc.namespace, testAcc.repository, testAcc.secretName)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(testAcc),
				),
			},
		},
	})
}

func TestAccResourceSecret_generate(t *testing.T) {
	configInit := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			generate {
				length = 16
				symbols = true
			}
		}
	`, testAcc.secretName, testAcc.path)

	configLengthUpdate := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			generate {
				length = 32
				symbols = true
			}
		}
	`, testAcc.secretName, testAcc.path)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: configInit,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(testAcc, func(s *terraform.InstanceState) error {
						if len(s.Attributes["data"]) != 16 {
							return fmt.Errorf("expected 'data' to contain a 16 char secret")
						}
						return nil
					}),
					checkSecretExistsRemotely(testAcc),
				),
			},
			{
				Config: configLengthUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(testAcc, func(s *terraform.InstanceState) error {
						if len(s.Attributes["data"]) != 32 {
							return fmt.Errorf("expected 'data' to contain newly generated 32 char secret")
						}
						return nil
					}),
					checkSecretExistsRemotely(testAcc),
				),
			},
		},
	})
}

func TestAccResourceSecret_deleteDetection(t *testing.T) {
	config := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			data = "secretpassword"
		}
	`, testAcc.secretName, testAcc.path)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				PreConfig: func() {
					// Delete secret outside of Terraform workspace
					client().Secrets().Delete(testAcc.path)
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true, // The externally deleted secrethub_secret should be planned in for recreation
			},
		},
	})
}

func TestAccResourceSecret_import(t *testing.T) {
	config := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			data = "secretpassword"
		}
	`, testAcc.secretName, testAcc.path)

	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      fmt.Sprintf("secrethub_secret.%v", testAcc.secretName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getSecretResourceState(s *terraform.State, values *testAccValues) (*terraform.InstanceState, error) {
	resourceState := s.Modules[0].Resources[fmt.Sprintf("secrethub_secret.%v", values.secretName)]
	if resourceState == nil {
		return nil, fmt.Errorf("resource '%v' not in tf state", values.secretName)
	}

	state := resourceState.Primary
	if state == nil {
		return nil, fmt.Errorf("resource has no primary instance")
	}

	return state, nil
}

func checkSecretExistsRemotely(values *testAccValues) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := *testAccProvider.Meta().(providerMeta).client

		_, err := client.Secrets().Get(values.path)
		if err != nil {
			return err
		}

		return nil
	}
}

func checkSecretResourceState(values *testAccValues, check func(s *terraform.InstanceState) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState := s.RootModule().Resources[fmt.Sprintf("secrethub_secret.%v", values.secretName)]
		if resourceState == nil {
			return fmt.Errorf("resource '%v' not in tf state", values.secretName)
		}

		state := resourceState.Primary
		if state == nil {
			return fmt.Errorf("resource has no primary instance")
		}

		return check(state)
	}
}

func TestMergeSecretPath(t *testing.T) {
	type args struct {
		prefix string
		path   string
	}
	cases := []struct {
		name string
		args args
		want string
	}{
		{
			"prefixed path",
			args{"myorg/db_passwords", "postgres"},
			"myorg/db_passwords/postgres",
		},
		{
			"abs path",
			args{"", "myorg2/database/postgres"},
			"myorg2/database/postgres",
		},
		{
			"path with redundant slashes",
			args{"myorg/db_passwords/", "/postgres"},
			"myorg/db_passwords/postgres",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := newCompoundSecretPath(c.args.prefix, c.args.path)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("newCompoundSecretPath() = %v, want %v", got, c.want)
			}
		})
	}
}
