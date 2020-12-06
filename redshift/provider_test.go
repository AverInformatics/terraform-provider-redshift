package redshift

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	if os.Getenv(resource.TestEnvVar) != "" {
		testAccProvider = Provider().(*schema.Provider)
		if err := testAccProvider.Configure(&terraform.ResourceConfig{}); err != nil {
			panic(err)
		}
		testAccProviders = map[string]terraform.ResourceProvider{
			"redshift": testAccProvider,
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("REDSHIFT_URL"); v == "" {
		t.Fatal("REDSHIFT_URL must be set for acceptance tests")
	}
	if v := os.Getenv("REDSHIFT_USER"); v == "" {
		t.Fatal("REDSHIFT_USER must be set for acceptance tests")
	}
	if v := os.Getenv("REDSHIFT_PASSWORD"); v == "" {
		t.Fatal("REDSHIFT_PASSWORD must be set for acceptance tests")
	}
}