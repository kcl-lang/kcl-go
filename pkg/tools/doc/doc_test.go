package doc

import (
	"github.com/stretchr/testify/assert"
	kcl "kcl-lang.io/kcl-go"
	"runtime"
	"strings"
	"testing"
)

func TestDocRender(t *testing.T) {
	tcases := [...]struct {
		source *kcl.KclType
		expect string
	}{
		{
			source: &kcl.KclType{
				SchemaName: "Person",
				SchemaDoc:  "Description of Schema Person",
				Properties: map[string]*kcl.KclType{"name": {
					Type: "string",
				}},
				Required: []string{"name"},
			},
			expect: `## Schema Person

Description of Schema Person

### Attributes

**name** *required*

` + "`" + `string` + "`" + `

todo: The description of the property

### Examples

todo: The example section

## Source Files

- [Person](todo: filepath)
`,
		},
	}

	context := GenContext{
		Format:           Markdown,
		IgnoreDeprecated: true,
	}

	for _, tcase := range tcases {
		content, err := context.renderContent(tcase.source)
		if err != nil {
			t.Errorf("render failed, err: %s", err)
		}
		expect := tcase.expect
		if runtime.GOOS == "windows" {
			expect = strings.ReplaceAll(tcase.expect, "\n", "\r\n")
		}
		assert.Equal(t, expect, string(content), "render content mismatch")
	}
}
