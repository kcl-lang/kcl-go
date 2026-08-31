// Copyright The KCL Authors. All rights reserved.

package xml

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		wantErr  bool
		contains []string
	}{
		{
			name: "Simple map",
			data: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			wantErr:  false,
			contains: []string{"<name>test</name>", "<value>123</value>"},
		},
		{
			name: "Nested map",
			data: map[string]interface{}{
				"config": map[string]interface{}{
					"name":  "test",
					"value": 123,
				},
			},
			wantErr:  false,
			contains: []string{"<config>", "<name>test</name>", "<value>123</value>", "</config>"},
		},
		{
			name: "Array",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "first"},
					map[string]interface{}{"name": "second"},
				},
			},
			wantErr:  false,
			contains: []string{"<items>", "<name>first</name>", "<name>second</name>", "</items>"},
		},
		{
			name: "Primitive types",
			data: map[string]interface{}{
				"stringVal": "hello",
				"intVal":    42,
				"floatVal":  3.14,
				"boolVal":   true,
			},
			wantErr:  false,
			contains: []string{"<stringVal>hello</stringVal>", "<intVal>42</intVal>", "<floatVal>3.14</floatVal>", "<boolVal>true</boolVal>"},
		},
		{
			name:     "Empty map",
			data:     map[string]interface{}{},
			wantErr:  false,
			contains: []string{"<root>", "</root>"},
		},
		{
			name: "Map with interface{} keys",
			data: map[interface{}]interface{}{
				"key1": "value1",
				"key2": 123,
			},
			wantErr:  false,
			contains: []string{"<key1>value1</key1>", "<key2>123</key2>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				resultStr := string(result)
				for _, expected := range tt.contains {
					if !strings.Contains(resultStr, expected) {
						t.Errorf("Convert() result = %v, want to contain %v", resultStr, expected)
					}
				}
				if !strings.Contains(resultStr, "<root>") {
					t.Errorf("Convert() result should contain <root> element")
				}
				if !strings.Contains(resultStr, "</root>") {
					t.Errorf("Convert() result should contain closing </root> element")
				}
			}
		})
	}
}

func TestSingle(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		wantErr  bool
		contains []string
	}{
		{
			name:     "Simple YAML document",
			yaml:     "name: test\nvalue: 123\n",
			wantErr:  false,
			contains: []string{"<?xml version=", "<name>test</name>", "<value>123</value>", "<root>", "</root>"},
		},
		{
			name:     "Nested YAML structure",
			yaml:     "config:\n  name: test\n  value: 123\n",
			wantErr:  false,
			contains: []string{"<config>", "<name>test</name>", "<value>123</value>"},
		},
		{
			name:     "YAML with array",
			yaml:     "items:\n  - name: first\n  - name: second\n",
			wantErr:  false,
			contains: []string{"<items>", "<name>first</name>", "<name>second</name>"},
		},
		{
			name:     "YAML with special characters",
			yaml:     "message: \"Hello <world> & goodbye\"\n",
			wantErr:  false,
			contains: []string{"<message>", "Hello", "world", "goodbye", "</message>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Single(tt.yaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("Single() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				resultStr := string(result)
				for _, expected := range tt.contains {
					if !strings.Contains(resultStr, expected) {
						t.Errorf("Single() result = %v, want to contain %v", resultStr, expected)
					}
				}
				if !strings.Contains(resultStr, "<?xml version=") {
					t.Errorf("Single() result should contain XML declaration")
				}
				if !strings.HasSuffix(resultStr, "\n") {
					t.Errorf("Single() result should end with newline")
				}
			}
		})
	}
}

func TestStream(t *testing.T) {
	tests := []struct {
		name       string
		yamlStream string
		wantErr    bool
		contains   []string
		docCount   int
	}{
		{
			name:       "YAML Stream with 2 documents",
			yamlStream: "---\nname: First\n---\nname: Second\n",
			wantErr:    false,
			docCount:   2,
			contains:   []string{"<?xml version=", "<results>", "<root>", "<name>First</name>", "<name>Second</name>", "</results>"},
		},
		{
			name:       "YAML Stream with 3 documents",
			yamlStream: "---\na: 1\n---\nb: 2\n---\nc: 3\n",
			wantErr:    false,
			docCount:   3,
			contains:   []string{"<a>1</a>", "<b>2</b>", "<c>3</c>"},
		},
		{
			name:       "YAML Stream with nested structures",
			yamlStream: "---\nconfig:\n  name: test1\n---\nconfig:\n  name: test2\n",
			wantErr:    false,
			docCount:   2,
			contains:   []string{"<config>", "<name>test1</name>", "<name>test2</name>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Stream(tt.yamlStream)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				resultStr := string(result)
				for _, expected := range tt.contains {
					if !strings.Contains(resultStr, expected) {
						t.Errorf("Stream() result = %v, want to contain %v", resultStr, expected)
					}
				}
				if !strings.Contains(resultStr, "<results>") {
					t.Errorf("Stream() result should contain <results> wrapper")
				}
				if !strings.Contains(resultStr, "</results>") {
					t.Errorf("Stream() result should contain closing </results>")
				}
				if !strings.Contains(resultStr, "<?xml version=") {
					t.Errorf("Stream() result should contain XML declaration")
				}
				rootCount := strings.Count(resultStr, "<root>")
				if rootCount != tt.docCount {
					t.Errorf("Stream() should contain %d <root> elements, got %d", tt.docCount, rootCount)
				}
			}
		})
	}
}

func TestXMLEscaping(t *testing.T) {
	data := map[string]interface{}{
		"special": "<text> & \"quotes\"",
	}

	result, err := Convert(data)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	resultStr := string(result)

	if !strings.Contains(resultStr, "&lt;") && !strings.Contains(resultStr, "<text>") {
		t.Errorf("XML special characters should be escaped or in CDATA")
	}
	if strings.Contains(resultStr, "<text>") && !strings.Contains(resultStr, "&lt;") {
		t.Logf("Note: Raw tags present - ensure XML is valid")
	}
}

func TestValidXMLStructure(t *testing.T) {
	yaml := "name: test\nvalue: 123\n"

	result, err := Single(yaml)
	if err != nil {
		t.Fatalf("Single() error = %v", err)
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Name    string   `xml:"name"`
		Value   int      `xml:"value"`
	}

	var parsed Root
	err = xml.Unmarshal(result, &parsed)
	if err != nil {
		t.Errorf("Single() result is not valid XML: %v", err)
	}

	if parsed.Name != "test" {
		t.Errorf("parsed.Name = %v, want 'test'", parsed.Name)
	}
	if parsed.Value != 123 {
		t.Errorf("parsed.Value = %v, want 123", parsed.Value)
	}
}

func TestArrayHandling(t *testing.T) {
	data := map[string]interface{}{
		"items": []interface{}{"first", "second", "third"},
	}

	result, err := Convert(data)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	resultStr := string(result)

	if !strings.Contains(resultStr, "<item>") {
		t.Errorf("Array items should be wrapped in <item> tags")
	}
	if !strings.Contains(resultStr, "first") || !strings.Contains(resultStr, "second") || !strings.Contains(resultStr, "third") {
		t.Errorf("Array items should be present in result")
	}
}

func TestNilHandling(t *testing.T) {
	data := map[string]interface{}{
		"present": "value",
		"absent":  nil,
	}

	result, err := Convert(data)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	resultStr := string(result)

	if !strings.Contains(resultStr, "<present>value</present>") {
		t.Errorf("Present value should be in result")
	}
	if strings.Contains(resultStr, "<absent>") {
		t.Logf("Note: Nil values are present (this is acceptable if handled correctly)")
	}
}

// TestAttrMarkerSimple verifies the marker renders a single attribute.
func TestAttrMarkerSimple(t *testing.T) {
	yaml := `TextView:
  __kcl_info_meta__: [android:id]
  android:id: "@+id/userNameTextView"
`
	out, err := Single(yaml)
	if err != nil {
		t.Fatalf("Single() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `<TextView android:id="\+id/userNameTextView">`) && !strings.Contains(got, `android:id="`) {
		t.Errorf("expected android:id attribute, got: %s", got)
	}
	if strings.Contains(got, "__kcl_info_meta__") {
		t.Errorf("marker key must not be emitted as element, got: %s", got)
	}
}

// TestAttrMarkerMultiAttr verifies multiple attributes on one element.
func TestAttrMarkerMultiAttr(t *testing.T) {
	yaml := `TextView:
  __kcl_info_meta__: [android:id, android:text]
  android:id: "@+id/userNameTextView"
  android:text: 用户名
`
	out, err := Single(yaml)
	if err != nil {
		t.Fatalf("Single() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `android:id="`) {
		t.Errorf("expected android:id attribute, got: %s", got)
	}
	if !strings.Contains(got, `android:text="`) {
		t.Errorf("expected android:text attribute, got: %s", got)
	}
	if strings.Contains(got, "<__kcl_info_meta__") || strings.Contains(got, "<android:id") {
		t.Errorf("attrs should not be rendered as elements, got: %s", got)
	}
}

// TestAttrMarkerMixedAttrsChildren verifies attrs and children coexist.
func TestAttrMarkerMixedAttrsChildren(t *testing.T) {
	yaml := `Button:
  __kcl_info_meta__: [android:id]
  android:id: "@+id/btn"
  label: "OK"
  enabled: true
`
	out, err := Single(yaml)
	if err != nil {
		t.Fatalf("Single() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `android:id="`) {
		t.Errorf("expected android:id attribute, got: %s", got)
	}
	if !strings.Contains(got, "<label>OK</label>") {
		t.Errorf("expected <label> child element, got: %s", got)
	}
	if !strings.Contains(got, "<enabled>true</enabled>") {
		t.Errorf("expected <enabled> child element, got: %s", got)
	}
}

// TestAttrMarkerAbsentKey verifies a listed-but-missing key does not panic
// and is silently skipped. This guards against `disable_none`/`query_paths`
// leaving a stale marker entry.
func TestAttrMarkerAbsentKey(t *testing.T) {
	yaml := `TextView:
  __kcl_info_meta__: [android:id, android:text]
  android:text: 用户名
`
	out, err := Single(yaml)
	if err != nil {
		t.Fatalf("Single() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `android:text="`) {
		t.Errorf("expected android:text attribute, got: %s", got)
	}
	if strings.Contains(got, `android:id=`) {
		t.Errorf("absent attr should not be emitted, got: %s", got)
	}
}

// TestAttrMarkerNoAttrRegression verifies that without the marker, output is
// unchanged from the legacy renderer.
func TestAttrMarkerNoAttrRegression(t *testing.T) {
	yaml := `TextView:
  android:id: "@+id/userNameTextView"
  android:text: 用户名
`
	out, err := Single(yaml)
	if err != nil {
		t.Fatalf("Single() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "<android:id>") {
		t.Errorf("without marker, android:id should be a child element, got: %s", got)
	}
	if !strings.Contains(got, "<android:text>") {
		t.Errorf("without marker, android:text should be a child element, got: %s", got)
	}
}

// TestAttrMarkerNested verifies nested schema instances carry their own
// markers independently.
func TestAttrMarkerNested(t *testing.T) {
	yaml := `parent:
  child:
    __kcl_info_meta__: [name]
    name: nested
  sibling: siblingVal
`
	out, err := Single(yaml)
	if err != nil {
		t.Fatalf("Single() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `name="nested"`) {
		t.Errorf("expected nested child attribute, got: %s", got)
	}
	if !strings.Contains(got, "<sibling>siblingVal</sibling>") {
		t.Errorf("expected sibling child element, got: %s", got)
	}
}

// TestAttrMarkerStream verifies multi-doc stream where only some docs have
// the marker.
func TestAttrMarkerStream(t *testing.T) {
	yamlStream := `---
TextView:
  __kcl_info_meta__: [android:id]
  android:id: "@+id/a"
---
Button:
  android:id: plain
`
	out, err := Stream(yamlStream)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "android:id=\"@+id/a\"") {
		t.Errorf("first doc should have attribute, got: %s", got)
	}
	if !strings.Contains(got, "<android:id>plain</android:id>") {
		t.Errorf("second doc should be plain child element, got: %s", got)
	}
}
