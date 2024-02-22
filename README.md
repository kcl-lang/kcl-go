# KCL Go SDK

[![GoDoc](https://godoc.org/github.com/kcl-lang/kcl-go?status.svg)](https://godoc.org/github.com/kcl-lang/kcl-go)
[![Coverage Status](https://coveralls.io/repos/github/kcl-lang/kcl-go/badge.svg)](https://coveralls.io/github/kcl-lang/kcl-go)
[![license](https://img.shields.io/github/license/kcl-lang/kcl-go.svg)](https://github.com/kcl-lang/kcl-go/blob/master/LICENSE)

[KCL](https://github.com/kcl-lang/kcl) is an open-source, constraint-based record and functional language that enhances the writing of complex configurations, including those for cloud-native scenarios. With its advanced programming language technology and practices, KCL is dedicated to promoting better modularity, scalability, and stability for configurations. It enables simpler logic writing and offers ease of automation APIs and integration with homegrown systems.

## What is it for?

You can use KCL to

+ [Generate low-level static configuration data](https://kcl-lang.io/docs/user_docs/guides/configuration) such as JSON, YAML, etc., or [integrate with existing data](https://kcl-lang.io/docs/user_docs/guides/data-integration).
+ Reduce boilerplate in configuration data with the [schema modeling](https://kcl-lang.io/docs/user_docs/guides/schema-definition).
+ Define schemas with [rule constraints for configuration data and validate](https://kcl-lang.io/docs/user_docs/guides/validation) them automatically.
+ Organize, simplify, unify and manage large configurations without side effects through [gradient automation schemes and GitOps](https://kcl-lang.io/docs/user_docs/guides/automation).
+ Manage large configurations in a scalable way with [isolated configuration blocks](https://kcl-lang.io/docs/reference/lang/tour#config-operations).
+ Mutating or validating Kubernetes resources with [cloud-native configuration tool plugins](https://kcl-lang.io/docs/user_docs/guides/working-with-k8s/).
+ Used as a platform engineering programming language to deliver modern applications with [Kusion Stack](https://kusionstack.io).

## Building & Testing

- [Install Go 1.21+](https://go.dev/dl/)

```bash
go test ./...
```

## Run KCL Code with Go SDK

```go
package main

import (
	"fmt"

	kcl "kcl-lang.io/kcl-go"
)

func main() {
	yaml := kcl.MustRun("kubernetes.k", kcl.WithCode(code)).GetRawYamlResult()
	fmt.Println(yaml)
}

const code = `
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
```

Run the command:

```bash
go run ./examples/kubernetes/main.go
```

Output:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

## Run KCL Code with Go Plugin

```go
package main

import (
	"fmt"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/native"                // Import the native API
	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin" // Import the hello plugin
)

func main() {
	// Note we use `native.MustRun` here instead of `kcl.MustRun`, because it needs the cgo feature.
	yaml := native.MustRun("main.k", kcl.WithCode(code)).GetRawYamlResult()
	fmt.Println(yaml)
}

const code = `
import kcl_plugin.hello

name = "kcl"
three = hello.add(1,2)  # hello.add is written by Go
`
```

## Documents

See the [KCL website](https://kcl-lang.io)

## License

Apache License Version 2.0
