package main

import (
	"kusionstack.io/kclvm-go/pkg/kpm"
	"os"
)

func main() {
	err := kpm.CLI(os.Args...)
	if err != nil {
		println(err.Error())
	}

}
