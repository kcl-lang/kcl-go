# KCL Go SDK

[![GoDoc](https://godoc.org/github.com/KusionStack/kclvm-go?status.svg)](https://godoc.org/github.com/KusionStack/kclvm-go)
[![Coverage Status](https://coveralls.io/repos/github/KusionStack/kclvm-go/badge.svg)](https://coveralls.io/github/KusionStack/kclvm-go)
[![license](https://img.shields.io/github/license/KusionStack/kclvm-go.svg)](https://github.com/KusionStack/kclvm-go/blob/master/LICENSE)

## Building

- [Install Go 1.19+](https://go.dev/dl/)
- [Install KCLVM](https://kcl-lang.io/docs/user_docs/getting-started/install)

```bash
$ go run ./cmds/kcl-go
$ go run ./cmds/kcl-go run hello.k
name: kcl
age: 1
two: 2
x0:
  name: kcl
  age: 1
x1:
  name: kcl
  age: 101
```

## Testing

```bash
go test ./...
```

## Run KCL Code with KCLVM Go SDK

```go
package main

import (
	"fmt"

	"kusionstack.io/kclvm-go"
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

## License

Apache License Version 2.0
