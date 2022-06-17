// Copyright 2021 The KCL Authors. All rights reserved.

package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
	"kusionstack.io/kclvm-go/scripts"
)

const (
	tEnvNumCpu      = "KCLVM_GO_API_TEST_NUM_CPU"
	tRestServerAddr = "127.0.0.1:7001"
	tKclGoPkg       = "kusionstack.io/kclvm-go/cmds/kcl-go"
)

func TestMain(m *testing.M) {
	if s := os.Getenv(tEnvNumCpu); s != "" {
		if x, err := strconv.Atoi(s); err == nil {
			fmt.Println("TestMain: nWorker =", x)
			kclvm_runtime.InitRuntime(x)
		}
	}

	// go install kcl-go
	if _, err := scripts.RunGoInstall(
		filepath.Join(kclvm_runtime.GetKclvmRoot(), "bin"),
		tKclGoPkg,
	); err != nil {
		panic(err)
	}

	go func() {
		if err := RunRestServer(tRestServerAddr); err != nil {
			log.Fatal(err)
		}
	}()

	// wait for http server ready
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		if _, err := http.Get("http://" + tRestServerAddr); err == nil {
			break
		}
	}

	os.Exit(m.Run())
}
