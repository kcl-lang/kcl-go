// Copyright The KCL Authors. All rights reserved.

package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"kcl-lang.io/kcl-go/pkg/kclvm_runtime"
)

const tEnvNumCpu = "KCLVM_GO_API_TEST_NUM_CPU"
const tRestServerAddr = "127.0.0.1:7001"

func TestMain(m *testing.M) {
	if s := os.Getenv(tEnvNumCpu); s != "" {
		if x, err := strconv.Atoi(s); err == nil {
			fmt.Println("TestMain: nWorker =", x)
			kclvm_runtime.InitRuntime(x)
		}
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
