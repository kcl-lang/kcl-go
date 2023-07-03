// Copyright 2023 The KCL Authors. All rights reserved.

// Run CGO Mode:
// KCLVM_SERVICE_CLIENT_HANDLER=native KCLVM_PLUGIN_DEBUG=1 go run -tags=kclvm_service_capi .

package main

import (
	"fmt"

	"kcl-lang.io/kcl-go"
	_ "kcl-lang.io/kcl-go/pkg/kcl_plugin/hello_plugin"
)

func main() {
	yaml := kclvm.MustRun("kubernetes.k", kclvm.WithCode(k_code)).GetRawYamlResult()
	fmt.Println(yaml)
}

const k_code = `
apiVersion = "apps/v1"
kind = "Deployment"
metadata = {
    name = "nginx"
    labels.app = "nginx"
}
spec = {
    replicas = 3
    selector.matchLabels = metadata.labels
    template.metadata.labels = metadata.labels
    template.spec.containers = [
        {
            name = metadata.name
            image = "${metadata.name}:1.14.2"
            ports = [{ containerPort = 80 }]
        }
    ]
}
`
