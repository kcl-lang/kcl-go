# KCL Go SDK

[![GoDoc](https://godoc.org/github.com/kcl-lang/kcl-go?status.svg)](https://godoc.org/github.com/kcl-lang/kcl-go)
[![Coverage Status](https://coveralls.io/repos/github/kcl-lang/kcl-go/badge.svg)](https://coveralls.io/github/kcl-lang/kcl-go)
[![license](https://img.shields.io/github/license/kcl-lang/kcl-go.svg)](https://github.com/kcl-lang/kcl-go/blob/master/LICENSE)

## Building

- [Install Go 1.19+](https://go.dev/dl/)

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

## Run KCL Code with Go SDK

```go
package main

import (
	"fmt"

	kcl "kcl-lang.io/kcl-go"
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
