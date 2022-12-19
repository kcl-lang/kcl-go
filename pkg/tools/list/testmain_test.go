package list

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
)

const (
	tRestServerAddr = "127.0.0.1:7001"
	tKclGoPkg       = "kusionstack.io/kclvm-go/cmds/kcl-go"
)

func TestMain(m *testing.M) {
	// go install kcl-go
	if err := tInstallKclGo(); err != nil {
		log.Fatal("run go install failed: ", err)
	}

	go func() {
		if err := tStartPyFastHttpServer(tRestServerAddr); err != nil {
			log.Fatal("start py fast http server failed: ", err)
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

func tInstallKclGo() error {
	var gobin string
	if runtime.GOOS == "windows" {
		gobin = kclvm_runtime.MustGetKclvmRoot()
	} else {
		gobin = filepath.Join(kclvm_runtime.MustGetKclvmRoot(), "bin")
	}

	if out, err := tRunGoInstall(gobin, tKclGoPkg); err != nil {

		return fmt.Errorf("%s: %s", err.Error(), string(out))
	}
	return nil
}

func tRunGoInstall(gobin, pkg string) (output []byte, err error) {
	cmd := exec.Command("go", "install", pkg)
	cmd.Env = []string{fmt.Sprintf("GOBIN=%s", gobin)}
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(kv), "GOBIN=") {
			cmd.Env = append(cmd.Env, kv)
		}
	}
	return cmd.CombinedOutput()
}

func tStartPyFastHttpServer(address string) error {
	args := []string{"-m", "kclvm.program.rpc-server", fmt.Sprintf("-http=%s", address)}
	cmd := exec.Command(kclvm_runtime.MustGetKclvmPath(), args...)
	return cmd.Start()
}
