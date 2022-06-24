// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/pkg/utils"
)

func newCleanCmd() *cli.Command {
	return &cli.Command{
		Name:  "clean",
		Usage: "remove cached files",
		Action: func(c *cli.Context) error {
			pkgroot, err := utils.FindPkgRoot(c.Args().First())
			if err != nil {
				fmt.Println("no cache found")
				return err
			}
			cache_path := filepath.Join(pkgroot, ".kclvm/cache")
			if isDir(cache_path) {
				if err := os.RemoveAll(cache_path); err == nil {
					fmt.Printf("%s removed\n", cache_path)
					return nil
				} else {
					fmt.Printf("remove %s failed\n", cache_path)
					return nil
				}
			}

			fmt.Println("no cache found")
			return nil
		},
	}
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
