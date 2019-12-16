package test

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

// func TestTerraformAwsAlbWebContainersSimple(t *testing.T) {
// 	t.Parallel()

// 	tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/simple")

// 	testName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
// 	loggingBucket := fmt.Sprintf("%s-logs", testName)
// 	awsRegion := "us-west-2"
// 	vpcAzs := aws.GetAvailabilityZones(t, awsRegion)[:3]

// 	terraformOptions := &terraform.Options{
// 		// The path to where our Terraform code is located
// 		TerraformDir: tempTestFolder,

// 		// Variables to pass to our Terraform code using -var options
// 		Vars: map[string]interface{}{
// 			"test_name":   testName,
// 			"logs_bucket": loggingBucket,
// 			"vpc_azs":     vpcAzs,
// 			"region":      awsRegion,
// 		},

// 		// Environment variables to set when running Terraform
// 		EnvVars: map[string]string{
// 			"AWS_DEFAULT_REGION": awsRegion,
// 		},
// 	}

// 	defer terraform.Destroy(t, terraformOptions)
// 	defer aws.EmptyS3Bucket(t, awsRegion, loggingBucket)
// 	terraform.InitAndApply(t, terraformOptions)
// }

func TestTerraformAwsAlbWebContainersSimpleHttp(t *testing.T) {
	t.Parallel()

	tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/simple")

	testName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	loggingBucket := fmt.Sprintf("%s-logs", testName)
	awsRegion := "us-east-2"
	vpcAzs := aws.GetAvailabilityZones(t, awsRegion)[:3]

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"test_name":   testName,
			"logs_bucket": loggingBucket,
			"vpc_azs":     vpcAzs,
			"region":      awsRegion,
		},

		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	defer aws.EmptyS3Bucket(t, awsRegion, loggingBucket)
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	albDomain := terraform.Output(t, terraformOptions, "alb_url")
	albURL := fmt.Sprintf("https://%s/", albDomain)
	expectedText := "Hello World"
	tlsConfig := tls.Config{}
	maxRetries := 20
	timeBetweenRetries := 30 * time.Second

	http_helper.HttpGetWithRetry(t, albURL, &tlsConfig, 200, expectedText, maxRetries, timeBetweenRetries)
}
