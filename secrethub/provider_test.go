package secrethub

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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
	path       string
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

	testAcc.path = newCompoundSecretPath(testAcc.namespace, testAcc.repository, testAcc.secretName)
}

func testAccPreCheck(t *testing.T) func() {
	return func() {
		err := testAcc.validate()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
