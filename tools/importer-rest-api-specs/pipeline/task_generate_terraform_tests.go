package pipeline

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/discovery"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/schema"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/testattributes"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func (t pipelineTask) generateTerraformTests(input discovery.ServiceInput, details resourcemanager.TerraformDetails, logger hclog.Logger) (*resourcemanager.TerraformDetails, error) {
	// TODO: @mbfrahry: go through and add the tests to each of the existing resources in details.Resources["blag"].Tests

	for resourceName, resourceDetails := range details.Resources {
		if !resourceDetails.Tests.Generate {
			continue
		}
		logger.Trace(fmt.Sprintf("Generating Tests for %q", resourceName))

		basicTest, err := generateBasicTestConfig(resourceDetails)
		if err != nil {
			return nil, err
		}
		if basicTest != nil {
			resourceDetails.Tests.BasicConfiguration = *basicTest
		}

		importTest, err := generateImportTestConfig(resourceDetails)
		if err != nil {
			return nil, err
		}
		if importTest != nil {
			resourceDetails.Tests.RequiresImportConfiguration = *importTest
		}
	}

	return &details, nil
}

func generateBasicTestConfig(input resourcemanager.TerraformResourceDetails) (*string, error) {
	f := hclwrite.NewEmptyFile()
	h := testattributes.TestAttributesHelpers{
		SchemaModels: input.SchemaModels,
	}
	if err := h.GetAttributesForTests(input.SchemaModels[input.SchemaModelName], *f.Body(), true); err != nil {
		return nil, err
	}

	// todo don't hardcode azurerm
	output := fmt.Sprintf(`
resource "azurerm_%[1]s" "test" {
%[2]s
}
`, schema.ConvertToSnakeCase(input.ResourceName), f.Bytes())
	return &output, nil
}

func generateImportTestConfig(input resourcemanager.TerraformResourceDetails) (*string, error) {
	f := hclwrite.NewEmptyFile()
	h := testattributes.TestAttributesHelpers{
		SchemaModels: input.SchemaModels,
	}
	if err := h.GetAttributesForTests(input.SchemaModels[input.SchemaModelName], *f.Body(), true); err != nil {
		return nil, err
	}

	for hclName := range f.Body().Attributes() {
		f.Body().SetAttributeTraversal(hclName, hcl.Traversal{
			hcl.TraverseRoot{
				// todo don't hardcode azurerm
				Name: fmt.Sprintf("azurerm_%s.test.%s", schema.ConvertToSnakeCase(input.ResourceName), hclName),
			},
		})
	}

	// todo don't hardcode azurerm
	output := fmt.Sprintf(`
resource "azurerm_%[1]s" "import" {
%[2]s
}
`, schema.ConvertToSnakeCase(input.ResourceName), f.Bytes())
	return &output, nil
}
