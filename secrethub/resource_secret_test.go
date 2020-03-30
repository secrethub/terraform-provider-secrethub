package secrethub

import (
	"fmt"
	"testing"

	"github.com/secrethub/secrethub-go/pkg/randchar"

	"github.com/secrethub/secrethub-go/internals/assert"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceSecret_writeAbsPath(t *testing.T) {
	config := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			value = "secretpassword"
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
			value = "secretpassword"
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
			value = "secretpassword"
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
				use_symbols = true
			}
		}
	`, testAcc.secretName, testAcc.path)

	configLengthUpdate := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			generate {
				length = 32
				use_symbols = true
			}
		}
	`, testAcc.secretName, testAcc.path)

	configSpecificCharsets := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			generate {
				length = 16
				charsets = ["numbers", "symbols"]
			}
		}
	`, testAcc.secretName, testAcc.path)

	configOneMinRule := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			generate {
				length = 16
				charsets = ["numbers", "letters"]
				min = {
					numbers = 15
				}
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
						if len(s.Attributes["value"]) != 16 {
							return fmt.Errorf("expected 'value' to contain a 16 char secret")
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
						if len(s.Attributes["value"]) != 32 {
							return fmt.Errorf("expected 'value' to contain newly generated 32 char secret")
						}
						return nil
					}),
					checkSecretExistsRemotely(testAcc),
				),
			},
			{
				Config: configSpecificCharsets,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(testAcc, func(s *terraform.InstanceState) error {
						if len(s.Attributes["value"]) != 16 {
							return fmt.Errorf("expected 'value' to contain newly generated 16 char secret")
						}
						if !containsOnly(s.Attributes["value"], randchar.Numeric.Add(randchar.Symbols)) {
							return fmt.Errorf("expected 'value' to only contain numbers and symbols")
						}
						return nil
					}),
					checkSecretExistsRemotely(testAcc),
				),
			},
			{
				Config: configOneMinRule,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(testAcc, func(s *terraform.InstanceState) error {
						if len(s.Attributes["value"]) != 16 {
							return fmt.Errorf("expected 'value' to contain newly generated 16 char secret")
						}
						if !containsOnly(s.Attributes["value"], randchar.Numeric.Add(randchar.Letters)) {
							return fmt.Errorf("expected 'value' to only contain numbers and letters")
						}
						if !containsAtLeast(s.Attributes["value"], randchar.Numeric, 15) {
							return fmt.Errorf("expected 'value' to contain at least 15 numbers")
						}
						return nil
					}),
					checkSecretExistsRemotely(testAcc),
				),
			},
		},
	})
}

func containsOnly(value string, charset randchar.Charset) bool {
	return randchar.NewCharset(value).IsSubset(charset)
}

func containsAtLeast(str string, charset randchar.Charset, count int) bool {
	for _, chr := range str {
		if randchar.NewCharset(string(chr)).IsSubset(charset) {
			count--
		}
	}
	return count <= 0
}

func TestAccResourceSecret_deleteDetection(t *testing.T) {
	config := fmt.Sprintf(`
		resource "secrethub_secret" "%v" {
			path = "%v"
			value = "secretpassword"
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
					err := client().Secrets().Delete(testAcc.path)
					assert.OK(t, err)
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
			value = "secretpassword"
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
