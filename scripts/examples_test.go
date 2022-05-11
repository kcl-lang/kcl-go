package scripts_test

import (
	"log"

	"kusionstack.io/kclvm-go/scripts"
)

func Example_all() {
	// scripts.KclvmDownloadUrlBase_mirrors = []string{ ... }

	if err := scripts.SetupKclvmAll("./_build"); err != nil {
		log.Fatal(err)
	}
}

func Example_setupKclvm() {
	// scripts.DefaultKclvmVersion = "...dev-version..."
	// scripts.KclvmDownloadUrlBase_mirrors = []string{ ... }

	scripts.DefaultKclvmTriple = "kclvm-centos"
	if err := scripts.SetupKclvm("./kclvm_root_centos"); err != nil {
		log.Fatal(err)
	}

	scripts.DefaultKclvmTriple = "kclvm-Darwin"
	if err := scripts.SetupKclvm("./kclvm_root_Darwin"); err != nil {
		log.Fatal(err)
	}

	scripts.DefaultKclvmTriple = "kclvm-Darwin-arm64"
	if err := scripts.SetupKclvm("./kclvm_root_Darwin_arm64"); err != nil {
		log.Fatal(err)
	}

	scripts.DefaultKclvmTriple = "kclvm-ubuntu"
	if err := scripts.SetupKclvm("./kclvm_root_ubuntu"); err != nil {
		log.Fatal(err)
	}
}
