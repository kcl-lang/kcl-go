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
	disableInstallArtifact bool
	// Can be inject value with -ldflags "-X kcl-lang.io/kcl-go/pkg/env.disableArtifactInPath=true"
	disableArtifactInPath bool
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
			DisableInstallArtifact:   envBoolOr("KCL_GO_DISABLE_INSTALL_ARTIFACT", disableInstallArtifact),
			DisableUseArtifactInPath: envBoolOr("KCL_GO_DISABLE_ARTIFACT_IN_PATH", disableArtifactInPath),
		}
	})
	return env
}

// GetLibHome returns the LibHome value from the singleton instance of EnvSettings.
func GetLibHome() string {
	return instance().LibHome
}

// SetLibHome sets the LibHome value in the singleton instance of EnvSettings.
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
