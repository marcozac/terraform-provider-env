package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"

	_ "embed"
)

const fileDataSourceTestOutputName = "test"

//go:embed testdata/.env.test
var testFileDataSourceData []byte

func TestFileDataSource_notRequired(t *testing.T) {
	filename := testFileDataSourceWriteFileHelper(t)
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   getenvFunctionTerraformVersionChecks,
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testFileDataSourceConfig(filename, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						fileDataSourceTestOutputName,
						knownvalue.MapExact(map[string]knownvalue.Check{
							"TEST_ENV_FILE_FOO": knownvalue.StringExact("foo"),
							"TEST_ENV_FILE_BAR": knownvalue.StringExact("bar"),
						}),
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
				Config: testFileDataSourceConfig("not_existing_file", false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						fileDataSourceTestOutputName,
						knownvalue.MapSizeExact(0),
					),
				},
			},
		},
	})
}

// testFileDataSourceWriteFileHelper writes the test data to a file and returns the filename.
func testFileDataSourceWriteFileHelper(t *testing.T) string {
	t.Helper()
	filename := filepath.Join(t.TempDir(), ".env")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.Write(testFileDataSourceData); err != nil {
		panic(err)
	}
	return filename
}

func testFileDataSourceConfig(filename string, required bool) string {
	return fmt.Sprintf(`
		data "env_file" "test" {
			path = "%s"
			required = %t
		}
		output "%s" {
			value = data.env_file.test.result
			sensitive = true
		}`,
		filename, required, fileDataSourceTestOutputName,
	)
}
