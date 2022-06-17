package list

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/scripts"
)

const (
	tRestfulAddr = "127.0.0.1:7001"
	tKclGoPkg    = "kusionstack.io/kclvm-go/cmds/kcl-go"
)

func TestMain(m *testing.M) {
	// go install kcl-go
	if _, err := scripts.RunGoInstall(
		filepath.Join(kclvm_runtime.GetKclvmRoot(), "bin"),
		tKclGoPkg,
	); err != nil {
		panic(err)
	}

	// start restful server
	{
		go func() {
			if err := service.RunRestServer(tRestfulAddr); err != nil {
				log.Fatal(err)
			}
		}()

		// wait for http server ready
		client := service.NewRestClient(tRestfulAddr)
		for i := 1; i <= 3; i++ {
			if _, err := client.Ping(context.Background(), nil); err == nil {
				break
			}
			time.Sleep(time.Second << i)
		}
	}

	os.Exit(m.Run())
}
