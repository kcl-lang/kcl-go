// Copyright The KCL Authors. All rights reserved.

package validate

import (
	"errors"
	"os"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// ExternalPackage is a single mapping from a KCL package name to its
// on-disk path, in the same format accepted by the `kcl run --external`
// flag. Each entry is forwarded to the KCL core service as an
// `ExternalPkg` argument so the validator can resolve external schemas
// referenced by the file being validated.
type ExternalPackage struct {
	// PkgName is the imported KCL package name (e.g. "k8s" or "ext").
	PkgName string
	// PkgPath is the path on disk where the package source lives.
	PkgPath string
}

// ValidateOptions represents the options for the Validate function.
type ValidateOptions struct {
	Schema        string // The schema to validate against.
	AttributeName string // The attribute name to validate.
	Format        string // The format of the data.
	// ExternalPackages is an optional list of external KCL packages the
	// validated file may depend on. Each entry maps a package name (as
	// used in `import <name>` statements) to its location on disk, so
	// the validator can resolve those imports during type checking.
	// When this list is empty the validator falls back to the
	// historical single-file behaviour and any unresolved external
	// import will surface as a `CannotFindModule` error, matching
	// pre-existing semantics.
	ExternalPackages []ExternalPackage
}

// toExternalPkgs converts the high-level ExternalPackage list into the
// gRPC representation expected by the KCL service. The result is nil
// when no packages were provided so that the gRPC payload stays as
// compact as before for callers that do not need this feature.
func (o *ValidateOptions) toExternalPkgs() []*gpyrpc.ExternalPkg {
	if o == nil || len(o.ExternalPackages) == 0 {
		return nil
	}
	out := make([]*gpyrpc.ExternalPkg, 0, len(o.ExternalPackages))
	for _, pkg := range o.ExternalPackages {
		if pkg.PkgName == "" || pkg.PkgPath == "" {
			continue
		}
		out = append(out, &gpyrpc.ExternalPkg{
			PkgName: pkg.PkgName,
			PkgPath: pkg.PkgPath,
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// Validate validates the given data file against the specified
// schema file with the provided options.
func Validate(dataFile, schemaFile string, opts *ValidateOptions) (ok bool, err error) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return false, err
	}
	if opts == nil {
		opts = &ValidateOptions{}
	}
	svc := kcl.Service()
	resp, err := svc.ValidateCode(&gpyrpc.ValidateCodeArgs{
		File:          schemaFile,
		Data:          string(data),
		Schema:        opts.Schema,
		AttributeName: opts.AttributeName,
		Format:        opts.Format,
		ExternalPkgs:  opts.toExternalPkgs(),
	})
	if err != nil {
		return false, err
	}
	var e error = nil
	if resp.ErrMessage != "" {
		e = errors.New(resp.ErrMessage)
	}
	return resp.Success, e
}

// ValidateCode validates the given in-memory data string against the
// provided in-memory schema code, with the given options. This is
// intended for callers that already hold the schema and data as
// strings (for example, a code editor performing live validation);
// callers passing files on disk should prefer [Validate].
func ValidateCode(data, code string, opts *ValidateOptions) (ok bool, err error) {
	if opts == nil {
		opts = &ValidateOptions{}
	}
	svc := kcl.Service()
	resp, err := svc.ValidateCode(&gpyrpc.ValidateCodeArgs{
		Data:          data,
		Code:          code,
		Schema:        opts.Schema,
		AttributeName: opts.AttributeName,
		Format:        opts.Format,
		ExternalPkgs:  opts.toExternalPkgs(),
	})
	if err != nil {
		return false, err
	}
	var e error = nil
	if resp.ErrMessage != "" {
		e = errors.New(resp.ErrMessage)
	}
	return resp.Success, e
}

// ValidateCodeFile validates the data read from `dataFile` against the
// given schema source. Both the data and code can also be supplied
// directly via `data`/`code`; when `dataFile` is non-empty it takes
// precedence over `data` and the contents of the file are read.
func ValidateCodeFile(dataFile, data, code string, opts *ValidateOptions) (ok bool, err error) {
	if opts == nil {
		opts = &ValidateOptions{}
	}
	svc := kcl.Service()
	resp, err := svc.ValidateCode(&gpyrpc.ValidateCodeArgs{
		Datafile:      dataFile,
		Data:          data,
		Code:          code,
		Schema:        opts.Schema,
		AttributeName: opts.AttributeName,
		Format:        opts.Format,
		ExternalPkgs:  opts.toExternalPkgs(),
	})
	if err != nil {
		return false, err
	}
	var e error = nil
	if resp.ErrMessage != "" {
		e = errors.New(resp.ErrMessage)
	}
	return resp.Success, e
}
