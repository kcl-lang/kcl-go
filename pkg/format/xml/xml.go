// Copyright The KCL Authors. All rights reserved.

// Package xml provides XML serialization for KCL YAML/JSON results.
//
// When a KCL schema instance carries the `__kcl_info_meta__` marker
// (a sibling key whose value is a list of attribute names), the listed keys
// are emitted as `name="value"` attributes on the parent element instead of
// as child elements. The marker key is consumed and never rendered.
//
// This file is vendored from kcl-lang.io/cli with the marker support
// layered on top.
package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/goccy/go-yaml"
	yamlformat "kcl-lang.io/kcl-go/pkg/format/yaml"
)

// InfoMetaAttr is the sibling key that lists which map keys carry the
// `@info(type="attr")` role for the parent element. The marker is generic:
// downstream emitters translate the role into their target format (XML
// attributes today, but the marker itself is named after the `@info`
// decorator so future roles — CDATA, comments, protobuf tags — can reuse
// the same channel).
//
// For XML rendering, the listed map keys are emitted as `name="value"`
// attributes on the parent element; their values come from those keys in
// the same map.
const InfoMetaAttr = "__kcl_info_meta__"

// Convert converts arbitrary data structures to XML format with a root element.
//
// If a map contains the `__kcl_info_meta__` key, its listed entries are
// emitted as XML attributes on the parent element. The marker is consumed
// (not emitted as a child element).
func Convert(data any) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("<root>")
	if err := encodeValue(&buf, data, ""); err != nil {
		return nil, err
	}
	buf.WriteString("</root>")
	return buf.Bytes(), nil
}

// Single converts a single YAML result to XML format.
func Single(yamlResult string) ([]byte, error) {
	var yamlData any
	if err := yaml.UnmarshalWithOptions([]byte(yamlResult), &yamlData); err != nil {
		return nil, err
	}
	out, err := Convert(yamlData)
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(out) + "\n"), nil
}

// Stream converts a YAML Stream to XML format with multiple root elements.
func Stream(yamlResult string) ([]byte, error) {
	docs, err := yamlformat.ParseStream(yamlResult)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	out.WriteString(xml.Header)
	out.WriteString("<results>\n")
	for _, doc := range docs {
		xmlData, err := Convert(doc)
		if err != nil {
			return nil, err
		}
		xmlStr := string(xmlData)
		out.WriteString("  ")
		out.WriteString(xmlStr)
		out.WriteString("\n")
	}
	out.WriteString("</results>\n")
	return out.Bytes(), nil
}

// encodeValue renders `data` either as a sequence of child elements (when it
// is a map or list at the top level) or as bare text (scalars). When the
// caller has already supplied a wrapper element name, `key` is non-empty and
// encodeElement is used; otherwise the children are emitted bare.
func encodeValue(buf *bytes.Buffer, data any, key string) error {
	if key != "" {
		return encodeElement(buf, key, data)
	}
	switch v := data.(type) {
	case map[string]any:
		return encodeMapChildren(buf, v)
	case map[any]any:
		stringMap := make(map[string]any, len(v))
		for k, val := range v {
			stringMap[fmt.Sprintf("%v", k)] = val
		}
		return encodeMapChildren(buf, stringMap)
	case []any:
		for _, item := range v {
			if err := encodeElement(buf, "item", item); err != nil {
				return err
			}
		}
		return nil
	case string:
		buf.WriteString(escapeString(v))
		return nil
	case int, int64, float64, bool:
		buf.WriteString(fmt.Sprintf("%v", v))
		return nil
	case nil:
		return nil
	default:
		buf.WriteString(fmt.Sprintf("%v", v))
		return nil
	}
}

// encodeMapChildren renders each (k, v) entry of `m` as a child element,
// honouring the `__kcl_info_meta__` marker on `m` itself if it carries one.
// The marker is consumed (not emitted as a child). When a child value is a
// map that itself carries the marker, the marker drives attribute emission
// on that child's element.
func encodeMapChildren(buf *bytes.Buffer, m map[string]any) error {
	for k, v := range m {
		if err := encodeElement(buf, k, v); err != nil {
			return err
		}
	}
	return nil
}

// encodeElement encodes `<key ...attrs...>children</key>` where attrs come
// from the `__kcl_info_meta__` marker on the value map (if any) and the
// marker key itself is suppressed from the children. Attributes are
// emitted sorted by name for deterministic output.
//
// encodeElement is the only place that reads the marker, so the marker is
// always honoured at the schema-instance level (one element), never leaked
// across siblings.
func encodeElement(buf *bytes.Buffer, key string, value any) error {
	stringMap, isMap := asStringMap(value)
	attrs, children := splitAttrs(stringMap, isMap)

	buf.WriteString("<")
	buf.WriteString(key)
	for _, attrName := range sortedKeys(attrs) {
		buf.WriteString(" ")
		buf.WriteString(attrName)
		buf.WriteString(`="`)
		xml.Escape(buf, []byte(attrs[attrName]))
		buf.WriteString(`"`)
	}
	buf.WriteString(">")

	if isMap {
		// Render each remaining child as its own element. The marker is
		// already removed from `children` (it is never re-checked here)
		// and the attribute keys are excluded so we don't double-emit.
		for _, childKey := range children {
			childVal := stringMap[childKey]
			if err := encodeElement(buf, childKey, childVal); err != nil {
				return err
			}
		}
	} else {
		if err := encodeValue(buf, value, ""); err != nil {
			return err
		}
	}

	buf.WriteString("</")
	buf.WriteString(key)
	buf.WriteString(">")
	return nil
}

// splitAttrs reads the `__kcl_info_meta__` marker from `m` and returns:
//   - attrs: a map of attribute-name to attribute-value, sourced from the
//     listed keys in `m`. Missing listed keys are silently skipped (the
//     renderer must not panic on absent names because `disable_none` /
//     `query_paths` may filter them).
//   - children: the keys of `m` to emit as child elements, in lexicographic
//     order. The marker key is excluded.
//
// When `m` is nil (the value was not a map) or has no marker, attrs is nil
// and children is empty; the caller falls back to scalar rendering.
func splitAttrs(m map[string]any, isMap bool) (map[string]string, []string) {
	if !isMap {
		return nil, nil
	}
	raw, ok := m[InfoMetaAttr]
	if !ok {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		return nil, keys
	}
	names := toStringList(raw)
	attrs := make(map[string]string, len(names))
	for _, name := range names {
		if v, ok := m[name]; ok && v != nil {
			attrs[name] = fmt.Sprintf("%v", v)
		}
	}
	childSet := make(map[string]struct{}, len(m))
	children := make([]string, 0, len(m))
	for k := range m {
		if k == InfoMetaAttr {
			continue
		}
		if _, isAttr := attrs[k]; isAttr {
			continue
		}
		childSet[k] = struct{}{}
		children = append(children, k)
	}
	_ = childSet
	return attrs, children
}

// toStringList coerces a YAML/JSON-decoded value into []string. The marker
// value comes from the KCL runtime as a list of strings; YAML decoding
// produces either []any or []string depending on input shape, so we accept
// both. Invalid marker values are treated as empty so a malformed marker
// degrades to "render as plain child element" rather than failing the
// whole document.
func toStringList(v any) []string {
	switch list := v.(type) {
	case []string:
		return list
	case []any:
		out := make([]string, 0, len(list))
		for _, item := range list {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

// asStringMap normalises both map[string]any and map[any]any into a
// map[string]any and reports whether `data` is a map at all. When `data`
// is not a map the returned map is nil and `isMap` is false.
func asStringMap(data any) (map[string]any, bool) {
	switch m := data.(type) {
	case map[string]any:
		return m, true
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, v := range m {
			out[fmt.Sprintf("%v", k)] = v
		}
		return out, true
	default:
		return nil, false
	}
}

// sortedKeys returns the keys of m in lexicographic order. Attribute order is
// not semantically meaningful in XML, but a stable order makes tests and
// golden-file comparisons deterministic.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j-1] > keys[j]; j-- {
			keys[j-1], keys[j] = keys[j], keys[j-1]
		}
	}
	return keys
}

// escapeString escapes special XML characters in a string.
func escapeString(s string) string {
	var buf bytes.Buffer
	xml.Escape(&buf, []byte(s))
	return buf.String()
}