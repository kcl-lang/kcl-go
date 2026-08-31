// Copyright The KCL Authors. All rights reserved.

// Package yaml provides YAML stream parsing helpers used by the format
// converters (xml, json, toml). It is vendored from kcl-lang.io/cli to avoid
// pulling the CLI as a dependency.
package yaml

import (
	"strings"

	"github.com/goccy/go-yaml"
)

// ParseStream parses a YAML Stream (multiple documents separated by ---) into a list of documents.
func ParseStream(yamlResult string) ([]any, error) {
	decoder := yaml.NewDecoder(strings.NewReader(yamlResult))
	var docs []any
	for {
		var doc any
		err := decoder.Decode(&doc)
		if err != nil {
			break
		}
		docs = append(docs, doc)
	}
	if len(docs) == 0 {
		// If no stream documents found, treat as single document
		var singleDoc any
		if err := yaml.UnmarshalWithOptions([]byte(yamlResult), &singleDoc); err != nil {
			return nil, err
		}
		docs = append(docs, singleDoc)
	}
	return docs, nil
}

// IsStream checks if the result is a YAML Stream (contains multiple documents separated by ---).
func IsStream(yamlResult string) bool {
	return strings.Contains(yamlResult, "---\n")
}
