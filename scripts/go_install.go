package scripts

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunGoInstall(gobin, pkg string) (output []byte, err error) {
	cmd := exec.Command("go", "install", pkg)
	cmd.Env = []string{fmt.Sprintf("GOBIN=%s", gobin)}
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(kv), "GOBIN=") {
			cmd.Env = append(cmd.Env, kv)
		}
	}
	return cmd.CombinedOutput()
}
