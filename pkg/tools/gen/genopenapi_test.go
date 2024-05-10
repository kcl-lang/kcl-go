package gen

import (
	"os"
	"path/filepath"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestExportOpenAPIV3Spec(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("get work directory failed")
	}
	pkgPath := filepath.Join(cwd, "testdata", "openapi", "app")
	spec, err := ExportOpenAPIV3Spec(pkgPath)
	if err != nil {
		t.Fatal(err)
	}
	json, err := spec.Components.Schemas["models.schema.v1.AppConfiguration"].MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	got := string(json)

	expect := `{"default":"","description":"AppConfiguration is a developer-centric definition that describes how to run an Application. This application model builds upon a decade of experience at AntGroup running super large scale internal developer platform, combined with best-of-breed ideas and practices from the community.","example":{"Default example":{"value":"# Instantiate an App with a long-running service and its image is \"nginx:v1\"\n\nimport models.schema.v1 as ac\nimport models.schema.v1.workload as wl\nimport models.schema.v1.workload.container as c\n\nappConfiguration = ac.AppConfiguration {\n    workload: wl.Service {\n        containers: {\n            \"nginx\": c.Container {\n                image: \"nginx:v1\"\n            }\n        }\n    }\n}"}},"properties":{"annotations":{"additionalProperties":{"default":"","type":"string"},"default":"","description":"Annotations are key/value pairs that attach arbitrary non-identifying metadata to resources.","type":"object","x-kcl-decorators":[{"name":"info","keywords":{"hidden":"True"}}],"x-kcl-dict-key-type":{"type":"string"}},"database":{"$ref":"#/definitions/models.schema.v1.accessories.Database"},"labels":{"additionalProperties":{"default":"","type":"string"},"default":"","description":"Labels can be used to attach arbitrary metadata as key-value pairs to resources.","type":"object","x-kcl-decorators":[{"name":"info","keywords":{"hidden":"True"}}],"x-kcl-dict-key-type":{"type":"string"}},"monitoring":{"$ref":"#/definitions/models.schema.v1.monitoring.Prometheus"},"opsRule":{"$ref":"#/definitions/models.schema.v1.trait.OpsRule"},"workload":{"default":"","description":"Workload defines how to run your application code. Currently supported workload profile\nincludes Service and Job.","type":"object","x-kcl-union-types":[{"description":"Service is a kind of workload profile that describes how to run your application code. This is typically used for long-running web applications that should \"never\" go down, and handle short-lived latency-sensitive web requests, or events.","ref":"#/definitions/models.schema.v1.workload.Service","referencedBy":["models.schema.v1.workload.Service"]},{"description":"Job is a kind of workload profile that describes how to run your application code. This is typically used for tasks that take from a few seconds to a few days to complete.","ref":"#/definitions/models.schema.v1.workload.Job","referencedBy":["models.schema.v1.workload.Job"]}]}},"required":["workload"],"type":"object","x-kcl-type":{"type":"AppConfiguration","import":{"package":"models.schema.v1","alias":"app_configuration.k"}}}`
	assert2.Equal(t, expect, got)
}

func TestExportSwaggerV2Spec(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("get work directory failed")
	}
	pkgPath := filepath.Join(cwd, "testdata", "doc", "k8s")
	got, err := ExportSwaggerV2SpecString(pkgPath)
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
