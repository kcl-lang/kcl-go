// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/pkg/tools/list"
)

func newLispAppCmd() *cli.Command {
	return &cli.Command{
		Hidden:    false,
		Name:      "list-app",
		Usage:     "list app files/packages ",
		ArgsUsage: "[pkgpath]",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "info",
				Usage: "show root and pkgpath",
				Value: true,
			},

			&cli.BoolFlag{
				Name:  "file",
				Usage: "list files",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "pkg",
				Usage: "list packages",
			},

			&cli.StringSliceFlag{
				Name:  "changed-files",
				Usage: "set changed files",
			},

			&cli.BoolFlag{
				Name:  "show-abs",
				Usage: "use abs path",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "show-index",
				Usage: "show index",
				Value: true,
			},

			&cli.BoolFlag{
				Name:  "include-all",
				Usage: "include all elems",
				Value: false,
			},

			&cli.StringFlag{
				Name:  "kcl-yaml-file",
				Usage: "set custom kcl.yaml file",
				Value: "kcl.yaml",
			},
			&cli.StringFlag{
				Name:  "project-yaml-file",
				Usage: "set custom project.yaml file",
				Value: "project.yaml",
			},
			&cli.BoolFlag{
				Name:  "use-fast-parser",
				Usage: "use fast parser",
			},
		},

		Action: func(c *cli.Context) error {
			var (
				flagShowInfo = c.Bool("info")

				flagListFile    = c.Bool("file")
				flagListPackage = c.Bool("pkg")

				flagHasAbsPath = c.Bool("show-abs")
				flagHasIndex   = c.Bool("show-index")

				flagAll = c.Bool("include-all")

				flagUseFastParser = c.Bool("use-fast-parser")
			)

			pkgroot, pkgpath, err := list.FindPkgInfo(c.Args().First())
			if err != nil {
				fmt.Println("ERR:", err)
				os.Exit(1)
			}

			var goodPath = func(i int, s string) string {
				if flagHasAbsPath {
					if flagHasIndex {
						return fmt.Sprintf("%d: %s", i, pkgroot+"/"+s)
					} else {
						return pkgroot + "/" + s
					}
				} else {
					if flagHasIndex {
						return fmt.Sprintf("%d: %s", i, s)
					} else {
						return s
					}
				}
			}

			if flagShowInfo {
				fmt.Println("pkgroot:", pkgroot)
				if pkgpath != "" {
					if flagHasAbsPath {
						fmt.Println("pkgpath:", pkgroot+"/"+pkgpath)
					} else {
						fmt.Println("pkgpath:", pkgpath)
					}
				}
			}

			if flagUseFastParser {
				depParser := list.NewSingleAppDepParser(pkgroot, list.Option{
					KclYaml:     c.String("kcl-yaml-file"),
					ProjectYaml: c.String("project-yaml-file"),
				})

				if flagListFile {
					if pkgpath != "" {
						appFiles, err := depParser.GetAppFiles(pkgpath, flagAll)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						for i, s := range appFiles {
							fmt.Println(goodPath(i, s))
						}
					}
				}
				if flagListPackage {
					if pkgpath != "" {
						appFiles, err := depParser.GetAppPkgs(pkgpath, flagAll)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						for i, s := range appFiles {
							fmt.Println(goodPath(i, s))
						}
					}
				}

				return nil
			}

			depParser := list.NewDepParser(pkgroot, list.Option{
				KclYaml:     c.String("kcl-yaml-file"),
				ProjectYaml: c.String("project-yaml-file"),
			})
			if err := depParser.GetError(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if flagListFile {
				if pkgpath != "" {
					for i, s := range depParser.GetAppFiles(pkgpath, flagAll) {
						fmt.Println(goodPath(i, s))
					}
				}
			}
			if flagListPackage {
				if pkgpath != "" {
					for i, s := range depParser.GetAppPkgs(pkgpath, flagAll) {
						fmt.Println(goodPath(i, s))
					}
				}
			}

			return nil
		},
	}
}
