package testattributes

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func testAttributeDiff(old string, newBody hclwrite.Body) ([]string, error) {
	oldConfig, diag := hclwrite.ParseConfig([]byte(old), "", hcl.Pos{Line: 1, Column: 1})
	if diag.HasErrors() {
		return []string{}, fmt.Errorf("parsing string config: %+v", diag.Errs())
	}

	diffs := append(findDiff(*oldConfig.Body(), newBody, ""), findBlockDiff(oldConfig.Body().Blocks(), newBody.Blocks(), "")...)

	return diffs, nil
}

func findDiff(body, diffBody hclwrite.Body, prefix string) []string {
	notFoundKeys := make([]string, 0)

	for key := range body.Attributes() {
		found := false
		for newKey := range diffBody.Attributes() {
			if newKey == key {
				found = true
				break
			}
		}
		if !found {
			if prefix == "" {
				notFoundKeys = append(notFoundKeys, key)
			} else {
				notFoundKeys = append(notFoundKeys, fmt.Sprintf("%s.%s", prefix, key))
			}
		}
	}

	return notFoundKeys
}

func findBlockDiff(oldBlocks, newBlocks []*hclwrite.Block, prefix string) []string {
	notFoundBlocks := make([]string, 0)

	for _, oldBlock := range oldBlocks {
		found := false
		for _, newBlock := range newBlocks {
			if oldBlock.Type() == newBlock.Type() {
				nestedDiff := make([]string, 0)
				if prefix == "" {
					nestedDiff = findDiff(*oldBlock.Body(), *newBlock.Body(), oldBlock.Type())
				} else {
					nestedDiff = findDiff(*oldBlock.Body(), *newBlock.Body(), fmt.Sprintf("%s.%s", prefix, oldBlock.Type()))
				}

				found = true
				notFoundBlocks = append(notFoundBlocks, nestedDiff...)
				break
			}
		}
		if !found {
			notFoundBlocks = append(notFoundBlocks, oldBlock.Type())
		}
	}

	return notFoundBlocks
}
