// Copyright The KCL Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package path

import (
	"os"
	"path/filepath"

	"kcl-lang.io/kcl-go/pkg/env"
)

// lazypath is an lazy-loaded path buffer for the XDG base directory specification.
type lazypath string

func (l lazypath) path(envVar string, defaultFn func() string, elem ...string) string {

	// There is an order to checking for a path.
	// 1. See if a KCL specific environment variable has been set.
	// 2. Fall back to a default
	base := os.Getenv(envVar)
	if base != "" {
		return filepath.Join(base, filepath.Join(elem...))
	}
	if base == "" {
		base = defaultFn()
	}
	return filepath.Join(base, string(l), filepath.Join(elem...))
}

// libPath defines the base directory relative to which user specific non-essential data files
// should be stored.
func (l lazypath) libPath(elem ...string) string {
	return l.path(env.GetLibHome(), libHome, filepath.Join(elem...))
}
