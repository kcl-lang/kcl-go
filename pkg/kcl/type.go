package kcl

import (
	"kcl-lang.io/kcl-go/pkg/source"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// GetSchemaType returns the schema types associated with the given schema name
// from a KCL file or source code provided.
//
// Parameters:
//   - filename: The name of the KCL file.
//   - src: The source content of the KCL file, which can be one of string, []byte or io.Reader types.
//   - schemaName: The name of the schema to get the types for.
//
// Returns:
//   - A slice of pointers to KclType representing the schema types.
//   - An error if there is any failure in the process.
func GetSchemaType(filename string, src any, schemaName string) ([]*gpyrpc.KclType, error) {
	mapping, err := GetSchemaTypeMapping(filename, src, schemaName)
	if err != nil {
		return nil, err
	}
	return getValues(mapping), nil
}

// GetFullSchemaType returns the full schema types for the given schema name
// from a list of KCL file paths.
//
// Parameters:
//   - pathList: A list of KCL file paths.
//   - schemaName: The name of the schema to get the types for.
//   - opts: Additional options to configure the processing.
//
// Returns:
//   - A slice of pointers to KclType representing the full schema types.
//   - An error if there is any failure in the process.
func GetFullSchemaType(pathList []string, schemaName string, opts ...Option) ([]*gpyrpc.KclType, error) {
	mapping, err := GetFullSchemaTypeMapping(pathList, schemaName, opts...)
	if err != nil {
		return nil, err
	}
	return getValues(mapping), nil
}

// GetFullSchemaTypeMapping returns the full schema type mapping for the given
// schema name from a list of KCL file paths.
//
// Parameters:
//   - pathList: A list of KCL file paths.
//   - schemaName: The name of the schema to get the type mapping for.
//   - opts: Additional options to configure the processing.
//
// Returns:
//   - A map where the key is the schema name and the value is a pointer to KclType representing the schema type.
//   - An error if there is any failure in the process.
func GetFullSchemaTypeMapping(pathList []string, schemaName string, opts ...Option) (map[string]*gpyrpc.KclType, error) {
	opts = append(opts, *NewOption().Merge(WithKFilenames(pathList...)))
	args, err := ParseArgs(pathList, opts...)
	if err != nil {
		return nil, err
	}

	svc := Service()
	resp, err := svc.GetSchemaTypeMapping(&gpyrpc.GetSchemaTypeMappingArgs{
		ExecArgs:   args.ExecProgramArgs,
		SchemaName: schemaName,
	})

	if err != nil {
		return nil, err
	}

	return resp.SchemaTypeMapping, nil
}

// GetSchemaTypeMapping returns the schema type mapping for the given schema name
// from a KCL file or source code provided.
//
// Parameters:
//   - filename: The name of the KCL file.
//   - src: The source content of the KCL file, which can be one of string, []byte or io.Reader types.
//   - schemaName: The name of the schema to get the type mapping for.
//
// Returns:
//   - A map where the key is the schema name and the value is a pointer to KclType representing the schema type.
//   - An error if there is any failure in the process.
func GetSchemaTypeMapping(filename string, src any, schemaName string) (map[string]*gpyrpc.KclType, error) {
	source, err := source.ReadSource(filename, src)
	if err != nil {
		return nil, err
	}
	svc := Service()
	resp, err := svc.GetSchemaTypeMapping(&gpyrpc.GetSchemaTypeMappingArgs{
		ExecArgs: &gpyrpc.ExecProgramArgs{
			KFilenameList: []string{filename},
			KCodeList:     []string{string(source)},
		},
		SchemaName: schemaName,
	})
	if err != nil {
		return nil, err
	}
	return resp.SchemaTypeMapping, nil
}
