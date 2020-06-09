package secrethub

import (
	"fmt"
	"os"
	"testing"

	"github.com/secrethub/secrethub-go/pkg/secrethub"
	"github.com/secrethub/secrethub-go/pkg/secretpath"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	envCredential        = "SECRETHUB_CREDENTIAL"
	envNamespace         = "SECRETHUB_TF_ACC_NAMESPACE"
	envRepo              = "SECRETHUB_TF_ACC_REPOSITORY"
	envSecondAccountName = "SECRETHUB_TF_ACC_SECOND_ACCOUNT_NAME"
	envAWSKMSKey         = "SECRETHUB_TF_ACC_AWS_KMS_KEY"
	envAWSRole           = "SECRETHUB_TF_ACC_AWS_ROLE"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAcc *testAccValues

type testAccValues struct {
	namespace         string
	repository        string
	secretName        string
	secondAccountName string
	path              string
	pathErr           error
	awsKmsKey         string
	awsRole           string
}

func (testAccValues) validate() error {
	if testAcc.namespace == "" || testAcc.repository == "" || testAcc.secondAccountName == "" || testAcc.awsKmsKey == "" || testAcc.awsRole == "" {
		return fmt.Errorf("the following environment variables need to be set: %s, %s, %s, %s, %s, %s", envCredential, envNamespace, envRepo, envSecondAccountName, envAWSKMSKey, envAWSRole)
	}
	return testAcc.pathErr
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"secrethub": testAccProvider,
	}

	testAcc = &testAccValues{
		namespace:         os.Getenv(envNamespace),
		repository:        os.Getenv(envRepo),
		secondAccountName: os.Getenv(envSecondAccountName),
		secretName:        "test_acc_secret",
		awsKmsKey:         os.Getenv(envAWSKMSKey),
		awsRole:           os.Getenv(envAWSRole),
	}

	testAcc.path = secretpath.Join(testAcc.namespace, testAcc.repository, testAcc.secretName)
}

func client() *secrethub.Client {
	return testAccProvider.Meta().(providerMeta).client
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
