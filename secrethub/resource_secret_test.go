package secrethub

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/keylockerbv/secrethub-go/pkg/api"
)

const (
	envNamespace = "SECRETHUB_TF_ACC_NAMESPACE"
	envRepo      = "SECRETHUB_TF_ACC_REPOSITORY"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAcc *testAccValues

type testAccValues struct {
	namespace  string
	repository string
	secretName string
	path       api.SecretPath
	pathErr    error
}

func (testAccValues) validate() error {
	if testAcc.namespace == "" || testAcc.repository == "" {
		return fmt.Errorf("the following environment variables need to be set: %v, %v", envNamespace, envRepo)
	}
	return testAcc.pathErr
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"secrethub": testAccProvider,
	}

	testAcc = &testAccValues{
		namespace:  os.Getenv(envNamespace),
		repository: os.Getenv(envRepo),
		secretName: "test_acc_secret",
	}

	testAcc.path, testAcc.pathErr = newCompoundSecretPath(testAcc.namespace, testAcc.repository, testAcc.secretName)
}

func testAccPreCheck(t *testing.T) func() {
	return func() {
		err := testAcc.validate()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAccSecret_writeAbsPath(t *testing.T) {
	config := fmt.Sprintf(`
		provider "secrethub" {
			credential = "${file("~/.secrethub/credential")}"
		}

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
			credential = "${file("~/.secrethub/credential")}"
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
			credential = "${file("~/.secrethub/credential")}"
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
		provider "secrethub" {
			credential = "${file("~/.secrethub/credential")}"
		}
		
		resource "secrethub_secret" "%v" {
			path = "%v"
			generate {
				length = 16
				symbols = true
			}
		}
	`, testAcc.secretName, testAcc.path)

	configLengthUpdate := fmt.Sprintf(`
		provider "secrethub" {
			credential = "${file("~/.secrethub/credential")}"
		}
		
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
		if _, err := getSecretResourceState(s, values); err != nil {
			return err
		}

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
		state, err := getSecretResourceState(s, values)
		if err != nil {
			return err
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
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"prefixed path",
			args{"myorg/db_passwords", "postgres"},
			"myorg/db_passwords/postgres",
			false,
		},
		{
			"abs path",
			args{"", "myorg2/database/postgres"},
			"myorg2/database/postgres",
			false,
		},
		{
			"path with redundant slashes",
			args{"myorg/db_passwords/", "/postgres"},
			"myorg/db_passwords/postgres",
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path, err := newCompoundSecretPath(c.args.prefix, c.args.path)
			got := string(path)
			if (err != nil) != c.wantErr {
				t.Errorf("newCompoundSecretPath() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("newCompoundSecretPath() = %v, want %v", got, c.want)
			}
		})
	}
}
