// Copyright The KCL Authors. All rights reserved.

//go:build cgo

package gen

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	kcl "kcl-lang.io/kcl-go"
)

// TestGenKclFromJsonSchemaNullableIfThen goes beyond the golden-file
// comparison in genkcl_test.go: it executes the schema generated for a
// nullable, conditionally constrained property (the shape reported in
// https://github.com/kcl-lang/kcl/issues/1876) and asserts the validation
// semantics of the emitted `check` guard.
func TestGenKclFromJsonSchemaNullableIfThen(t *testing.T) {
	input := filepath.Join("testdata", "jsonschema", "nullable-types", "input.json")
	var buf bytes.Buffer
	if err := GenKcl(&buf, input, nil, &GenKclOptions{Mode: ModeJsonSchema}); err != nil {
		t.Fatal(err)
	}
	schema := buf.String()
	run := func(instance string) error {
		_, err := kcl.Run("test.k", kcl.WithCode(schema+"\n"+instance))
		return err
	}

	// The `if: {type: "string"}` / `then: {maxLength: 255}` triple only
	// constrains string values: the `["string", "null"]` property still
	// accepts None, while an over-long string must be rejected.
	cases := []struct {
		instance string
		wantErr  bool
	}{
		{"x = Input{pullSecret = None}", false},
		{`x = Input{pullSecret = ""}`, false},
		{`x = Input{pullSecret = "ok"}`, false},
		{`x = Input{pullSecret = "` + strings.Repeat("a", 256) + `"`, true},
		// The `$ref`'d `empty.string` default applies on its own.
		{"x = Input{podPriority = {}}", false},
	}
	for _, c := range cases {
		err := run(c.instance)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: err=%v, wantErr=%v", c.instance, err, c.wantErr)
		}
	}
}
