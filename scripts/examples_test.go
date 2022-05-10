package scripts_test

import (
	"log"

	"kusionstack.io/kclvm-go/scripts"
)

func Example() {
	if err := scripts.SetupKclvm("./kclvm_root"); err != nil {
		log.Fatal(err)
	}
}
