package dataapigenerator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func codeForTerraformConstant(namespace, constantName string, details resourcemanager.ConstantDetails) (*string, error) {
	code := make([]string, 0)

	sortedKeys := make([]string, 0)
	for key := range details.Values {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		value := details.Values[key]
		code = append(code, fmt.Sprintf("\t\t[Description(%q)]\n\t\t%s,", value, key))
	}

	attributes := make([]string, 0)
	constantFieldType, err := mapConstantFieldType(details.Type)
	if err != nil {
		return nil, fmt.Errorf("mapping constant field type %q: %+v", string(details.Type), err)
	}
	attributes = append(attributes, fmt.Sprintf("\t[ConstantType(%s)]", *constantFieldType))

	out := fmt.Sprintf(`using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace %[1]s;

%[4]s
public enum %[2]sConstant
{
%[3]s
}
`, namespace, constantName, strings.Join(code, "\n\n"), strings.Join(attributes, "\n"))
	return &out, nil
}
