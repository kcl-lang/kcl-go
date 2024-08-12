package gen

import (
	"errors"
	"io"
	"path/filepath"
	"strings"

	"kcl-lang.io/kcl-go/pkg/source"
)

type GenKclOptions struct {
	Mode                  Mode
	CastingOption         CastingOption
	UseIntegersForNumbers bool
}

// Mode is the mode of kcl schema code generation.
type Mode int

const (
	ModeAuto Mode = iota
	ModeGoStruct
	ModeJsonSchema
	ModeTerraformSchema
	ModeJson
	ModeYaml
	ModeToml
	ModeProto
	ModeTextProto
)

type kclGenerator struct {
	opts *GenKclOptions
}

// GenKcl translate other formats to kcl schema code. Now support go struct and json schema.
func GenKcl(w io.Writer, filename string, src interface{}, opts *GenKclOptions) error {
	return newKclGenerator(opts).GenSchema(w, filename, src)
}

func newKclGenerator(opts *GenKclOptions) *kclGenerator {
	if opts == nil {
		opts = new(GenKclOptions)
	}
	return &kclGenerator{
		opts: opts,
	}
}

func (k *kclGenerator) GenSchema(w io.Writer, filename string, src interface{}) error {
	if k.opts.Mode == ModeAuto {
		code, err := source.ReadSource(filename, src)
		if err != nil {
			return err
		}
		codeStr := string(code)
		switch filepath.Ext(filename) {
		case ".json":
			switch {
			case strings.Contains(codeStr, "$schema"):
				k.opts.Mode = ModeJsonSchema
			case strings.Contains(codeStr, "\"provider_schemas\""):
				k.opts.Mode = ModeTerraformSchema
			default:
				k.opts.Mode = ModeJson
			}
		case ".yaml", "yml":
			k.opts.Mode = ModeYaml
		case ".toml":
			k.opts.Mode = ModeToml
		case ".go":
			k.opts.Mode = ModeGoStruct
		case ".proto":
			k.opts.Mode = ModeProto
		case ".textproto":
			k.opts.Mode = ModeTextProto
		default:
			return errors.New("failed to detect mode")
		}
	}

	switch k.opts.Mode {
	case ModeGoStruct:
		return k.genSchemaFromGoStruct(w, filename, src)
	case ModeJsonSchema:
		return k.genSchemaFromJsonSchema(w, filename, src)
	case ModeTerraformSchema:
		return k.genSchemaFromTerraformSchema(w, filename, src)
	case ModeJson:
		return k.genKclFromJsonData(w, filename, src)
	case ModeYaml:
		return k.genKclFromYaml(w, filename, src)
	case ModeToml:
		return k.genKclFromToml(w, filename, src)
	case ModeProto:
		return k.genKclFromProto(w, filename, src)
	default:
		return errors.New("unknown mode")
	}
}
