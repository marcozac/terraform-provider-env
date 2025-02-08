package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"

	_ "embed"
)

const fileDataSourceTestOutputName = "test"

var (
	//go:embed testdata/.env.test
	testFileDataSourceData []byte

	//go:embed testdata/.env.err.test
	testFileDataSourceDataErr []byte
)

func TestFileDataSource_notRequired(t *testing.T) {
	dir := t.TempDir()
	filename := testFileDataSourceWriteFile(filepath.Join(dir, ".env"), testFileDataSourceData)
	errFilename := testFileDataSourceWriteFile(filepath.Join(dir, ".env.err"), testFileDataSourceDataErr)

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   getenvFunctionTerraformVersionChecks,
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testFileDataSourceConfig(filename),
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
				Config: fmt.Sprintf(`
					data "env_file" "test" {}
					output "%s" {
						value = data.env_file.test.result
						sensitive = true
					}`,
					fileDataSourceTestOutputName,
				),
				ExpectError: regexp.MustCompile("Failed to open file"),
			},
		},
	})
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   getenvFunctionTerraformVersionChecks,
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testFileDataSourceConfig(errFilename),
				ExpectError: regexp.MustCompile("Failed to parse file"),
			},
		},
	})
}

// testFileDataSourceWriteFileHelper writes the test data to a file and returns
// the filename as is.
func testFileDataSourceWriteFile(filename string, data []byte) string {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		panic(err)
	}
	return filename
}

func testFileDataSourceConfig(filename string) string {
	return fmt.Sprintf(`
		data "env_file" "test" {
			path = "%s"
		}
		output "%s" {
			value = data.env_file.test.result
			sensitive = true
		}`,
		filename, fileDataSourceTestOutputName,
	)
}
