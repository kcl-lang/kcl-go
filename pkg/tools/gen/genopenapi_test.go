package gen

import (
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestExportSwaggerV2Spec(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("get work directory failed")
	}
	pkgPath := filepath.Join(cwd, "testdata", "doc", "k8s")
	got, err := ExportSwaggerV2Spec(pkgPath)
	if err != nil {
		t.Fatal(err)
	}

	expect := `{
    "definitions": {
        "apps.Deployment": {
            "type": "object",
            "properties": {
                "metadata": {
                    "type": "string"
                },
                "podSpec": {
                    "type": "object"
                }
            },
            "required": [
                "metadata",
                "podSpec"
            ],
            "x-kcl-type": {
                "type": "Deployment",
                "import": {
                    "package": "apps",
                    "alias": "deployment.k"
                }
            }
        },
        "core.PodSpec": {
            "type": "object",
            "properties": {
                "image": {
                    "type": "string"
                }
            },
            "required": [
                "image"
            ],
            "x-kcl-type": {
                "type": "PodSpec",
                "import": {
                    "package": "core",
                    "alias": "podSpec.k"
                }
            }
        }
    },
    "paths": {},
    "swagger": "2.0",
    "info": {
        "title": "",
        "version": ""
    }
}`

	assert2.Equal(t, expect, got)
}
