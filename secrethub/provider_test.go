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
	envAWSRole           = "SECRETHUB_TF_ACC_AWS_ROLE"
	envAWSKMSKey         = "SECRETHUB_TF_ACC_AWS_KMS_KEY"
	envGCPServiceAccount = "SECRETHUB_TF_ACC_GCP_SERVICE_ACCOUNT"
	envGCPKMSKey         = "SECRETHUB_TF_ACC_GCP_KMS_KEY"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAcc *testAccValues

type testAccValues struct {
	namespace         string
	repository        string
	secretName        string
	secondAccountName string
	secretPath        string
	dirName           string
	dirPath           string
	repoPath          string
	awsRole           string
	awsKmsKey         string
	gcpServiceAccount string
	gcpKmsKey         string
}

func (testAccValues) validate() error {
	if testAcc.namespace == "" || testAcc.repository == "" || testAcc.secondAccountName == "" {
		return fmt.Errorf("make sure you set environment variables: %s, %s, %s, %s", envCredential, envNamespace, envRepo, envSecondAccountName)
	}
	return nil
}

func (testAccValues) validateAWS() error {
	if testAcc.awsKmsKey == "" || testAcc.awsRole == "" {
		return fmt.Errorf("make sure you set environment variables: %s, %s", envAWSKMSKey, envAWSRole)
	}
	return nil
}

func (testAccValues) validateGCP() error {
	if testAcc.gcpKmsKey == "" || testAcc.gcpServiceAccount == "" {
		return fmt.Errorf("make sure you set environment variables: %s, %s", envGCPKMSKey, envGCPServiceAccount)
	}
	return nil
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
		dirName:           "test_acc_dir",
		awsKmsKey:         os.Getenv(envAWSKMSKey),
		awsRole:           os.Getenv(envAWSRole),
		gcpKmsKey:         os.Getenv(envGCPKMSKey),
		gcpServiceAccount: os.Getenv(envGCPServiceAccount),
	}

	testAcc.repoPath = secretpath.Join(testAcc.namespace, testAcc.repository)
	testAcc.secretPath = secretpath.Join(testAcc.repoPath, testAcc.secretName)
	testAcc.dirPath = secretpath.Join(testAcc.repoPath, testAcc.dirName)
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

func testAccPreCheckAWS(t *testing.T) func() {
	return func() {
		testAccPreCheck(t)()
		err := testAcc.validateAWS()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testAccPreCheckGCP(t *testing.T) func() {
	return func() {
		testAccPreCheck(t)()
		err := testAcc.validateGCP()
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
