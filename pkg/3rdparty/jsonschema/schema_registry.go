package jsonschema

import (
	"context"
	"fmt"
	"strings"
)

var sr *SchemaRegistry

// SchemaRegistry maintains a lookup table between schema string references
// and actual schemas
type SchemaRegistry struct {
	schemaLookup  map[string]*Schema
	contextLookup map[string]*Schema
}

// GetSchemaRegistry provides an accessor to a globally available schema registry
func GetSchemaRegistry() *SchemaRegistry {
	if sr == nil {
		sr = &SchemaRegistry{
			schemaLookup:  map[string]*Schema{},
			contextLookup: map[string]*Schema{},
		}
	}
	return sr
}

// ResetSchemaRegistry resets the main SchemaRegistry
func ResetSchemaRegistry() {
	sr = nil
}

// Get fetches a schema from the top level context registry or fetches it from a remote
func (sr *SchemaRegistry) Get(ctx context.Context, uri string) *Schema {
	uri = strings.TrimRight(uri, "#")
	schema := sr.schemaLookup[uri]
	if schema == nil {
		fetchedSchema := &Schema{}
		err := FetchSchema(ctx, uri, fetchedSchema)
		if err != nil {
			schemaDebug(fmt.Sprintf("[SchemaRegistry] Fetch error: %s", err.Error()))
			return nil
		}

		fetchedSchema.DocPath = uri
		// TODO(arqu): meta validate schema
		schema = fetchedSchema
		sr.schemaLookup[uri] = schema
	}
	return schema
}

// GetKnown fetches a schema from the top level context registry
func (sr *SchemaRegistry) GetKnown(uri string) *Schema {
	uri = strings.TrimRight(uri, "#")
	return sr.schemaLookup[uri]
}

// GetLocal fetches a schema from the local context registry
func (sr *SchemaRegistry) GetLocal(uri string) *Schema {
	uri = strings.TrimRight(uri, "#")
	return sr.contextLookup[uri]
}

// Register registers a schema to the top level context
func (sr *SchemaRegistry) Register(sch *Schema) {
	if sch.DocPath == "" {
		return
	}
	sr.schemaLookup[sch.DocPath] = sch
}

// RegisterLocal registers a schema to a local context
func (sr *SchemaRegistry) RegisterLocal(sch *Schema) {
	if sch.Id != "" && IsLocalSchemaID(sch.Id) {
		sr.contextLookup[sch.Id] = sch
	}

	if sch.HasKeyword("$anchor") {
		anchorKeyword := sch.Keywords["$anchor"].(*Anchor)
		anchorURI := sch.DocPath + "#" + string(*anchorKeyword)
		if sr.contextLookup == nil {
			sr.contextLookup = map[string]*Schema{}
		}
		sr.contextLookup[anchorURI] = sch
	}
}
