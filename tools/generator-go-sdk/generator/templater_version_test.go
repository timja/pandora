package generator

import (
	"testing"
)

func TestTemplateVersion(t *testing.T) {
	input := ServiceGeneratorData{
		baseClientPackage: "pandas-api",
		packageName:       "somepackage",
		apiVersion:        "2022-02-01",
		source:            AccTestLicenceType,
	}

	actual, err := versionTemplater{}.template(input)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := `package somepackage

import "fmt"

// acctests licence placeholder

const defaultApiVersion = "2022-02-01"

func userAgent() string {
	return fmt.Sprintf("hashicorp/go-azure-sdk/pandas-api/somepackage/%s", defaultApiVersion)
}`
	assertTemplatedCodeMatches(t, expected, *actual)
}
