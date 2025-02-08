package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	getenvFunctionTestVarName    = "TEST_GETENV_VALUE"
	getenvFunctionTestVarValue   = "testvalue"
	getenvFunctionTestOutputName = "test"
)

var getenvFunctionTerraformVersionChecks = []tfversion.TerraformVersionCheck{
	tfversion.SkipBelow(minTfVersion),
}

func TestGetenvFunction_notRequired(t *testing.T) {
	testGetenvSetenvHelper(t, getenvFunctionTestVarValue)
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   getenvFunctionTerraformVersionChecks,
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testGetenvFunctionConfig(getenvFunctionTestVarName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						getenvFunctionTestOutputName,
						knownvalue.StringExact(getenvFunctionTestVarValue),
					),
				},
			},
		},
	})
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   getenvFunctionTerraformVersionChecks,
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testGetenvFunctionConfig("TEST_GETENV_NOT_EXISTING_VALUE", false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						getenvFunctionTestOutputName,
						knownvalue.StringExact(""),
					),
				},
			},
		},
	})
}

func TestGetenvFunction_required(t *testing.T) {
	testGetenvSetenvHelper(t, getenvFunctionTestVarValue)
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   getenvFunctionTerraformVersionChecks,
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testGetenvFunctionConfig(getenvFunctionTestVarName, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						getenvFunctionTestOutputName,
						knownvalue.StringExact(getenvFunctionTestVarValue),
					),
				},
			},
		},
	})
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   getenvFunctionTerraformVersionChecks,
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testGetenvFunctionConfig("TEST_GETENV_NOT_EXISTING_VALUE", true),
				ExpectError: regexp.MustCompile("Environment variable not found"),
			},
		},
	})
}

func testGetenvSetenvHelper(t *testing.T, v string) {
	t.Helper()
	t.Setenv(getenvFunctionTestVarName, v)
}

func testGetenvFunctionConfig(name string, required bool) string {
	return fmt.Sprintf(`
		output "%s" {
			value = provider::env::getenv("%s", %t)
		}`,
		getenvFunctionTestOutputName, name, required,
	)
}
