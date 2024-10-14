package env

import (
	"os"
	"strconv"
	"sync"
)

var (
	once sync.Once
	mu   sync.Mutex
	env  *EnvSettings
	// Can be inject value with -ldflags "-X kcl-lang.io/kcl-go/pkg/env.libHome=/path/to/lib"
	libHome string
	// Can be inject value with -ldflags "-X kcl-lang.io/kcl-go/pkg/env.disableInstallArtifact=true"
	disableInstallArtifact string = "false"
	// Can be inject value with -ldflags "-X kcl-lang.io/kcl-go/pkg/env.disableArtifactInPath=true"
	disableArtifactInPath string = "true"
)

// EnvSettings represents environment settings for the KCL Go SDK.
type EnvSettings struct {
	LibHome                  string
	DisableInstallArtifact   bool
	DisableUseArtifactInPath bool
}

// Instance returns a singleton instance of EnvSettings.
func instance() *EnvSettings {
	once.Do(func() {
		env = &EnvSettings{
			LibHome:                  envOr(os.Getenv("KCL_LIB_HOME"), libHome),
			DisableInstallArtifact:   envBoolOrString("KCL_GO_DISABLE_INSTALL_ARTIFACT", disableInstallArtifact),
			DisableUseArtifactInPath: envBoolOrString("KCL_GO_DISABLE_ARTIFACT_IN_PATH", disableArtifactInPath),
		}
	})
	return env
}

// GetLibHome returns the LibHome value from the singleton instance of EnvSettings.
// Deprecated: use KCL_LIB_HOME env to set the kcl lib home
func GetLibHome() string {
	return instance().LibHome
}

// SetLibHome sets the LibHome value in the singleton instance of EnvSettings.
// Deprecated: use KCL_LIB_HOME env to get the kcl lib home
func SetLibHome(value string) {
	mu.Lock()
	defer mu.Unlock()
	instance().LibHome = value
}

// SetDisableInstallArtifact sets the DisableInstallArtifact value in the singleton instance of EnvSettings.
func SetDisableInstallArtifact(value bool) {
	mu.Lock()
	defer mu.Unlock()
	instance().DisableInstallArtifact = value
}

// GetDisableInstallArtifact returns the DisableInstallArtifact value from the singleton instance of EnvSettings.
func GetDisableInstallArtifact() bool {
	return instance().DisableInstallArtifact
}

// GetDisableUseArtifactInPath returns the DisableUseArtifactInPath value from the singleton instance of EnvSettings.
func GetDisableUseArtifactInPath() bool {
	return instance().DisableUseArtifactInPath
}

// SetDisableUseArtifactInPath sets the DisableUseArtifactInPath value in the singleton instance of EnvSettings.
func SetDisableUseArtifactInPath(value bool) {
	mu.Lock()
	defer mu.Unlock()
	instance().DisableUseArtifactInPath = value
}

// Enable the fast eval mode.
func EnableFastEvalMode() {
	// Set the fast eval mode for KCL
	os.Setenv("KCL_FAST_EVAL", "1")
}

// envOr returns the value of the specified environment variable, or the
// default value if it does not exist.
func envOr(name, def string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return def
}

// envBoolOr returns the boolean value of the specified environment variable,
// or the default value if it does not exist.
func envBoolOr(name string, def bool) bool {
	if name == "" {
		return def
	}
	envVal := envOr(name, strconv.FormatBool(def))
	ret, err := strconv.ParseBool(envVal)
	if err != nil {
		return def
	}
	return ret
}

// envBoolOrString returns the boolean value of the specified environment variable,
// or the default string value if it does not exist.
func envBoolOrString(name string, def string) bool {
	var ret bool
	ret, _ = strconv.ParseBool(def)
	return envBoolOr(name, ret)
}
