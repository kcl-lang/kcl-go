package main

import (
	"fmt"

	kcl "kcl-lang.io/kcl-go"
	_ "kcl-lang.io/kcl-go/pkg/kcl_plugin/hello_plugin"
)

func main() {
	yaml := kcl.MustRun("kubernetes.k", kcl.WithCode(k_code)).GetRawYamlResult()
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
