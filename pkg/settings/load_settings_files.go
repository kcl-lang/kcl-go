// Copyright The KCL Authors. All rights reserved.

package settings

import (
	"kcl-lang.io/kcl-go/pkg/native"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type LoadSettingsFilesResult = gpyrpc.LoadSettingsFilesResult

// LoadSettingsFiles loads the KCL settings files and returns the merged
// CLI configs and options.
func LoadSettingsFiles(workDir string, files []string) (*LoadSettingsFilesResult, error) {
	svc := native.NewNativeServiceClient()
	return svc.LoadSettingsFiles(&gpyrpc.LoadSettingsFilesArgs{
		WorkDir: workDir,
		Files:   files,
	})
}
