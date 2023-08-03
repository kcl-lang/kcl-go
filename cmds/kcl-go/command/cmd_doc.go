// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"kcl-lang.io/kcl-go/pkg/tools/doc"
)

const version = "v0.0.1"

func NewDocCmd() *cli.Command {
	return &cli.Command{
		Hidden: true,
		Name:   "doc",
		Usage:  "show documentation for package or symbol",
		UsageText: `# Generate document for current package
kcl-go doc generate

# Start a local KCL document server
kcl-go doc start`,
		Subcommands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "generates documents from code and examples",
				UsageText: `# Generate Markdown document for current package
kcl-go doc generate

# Generate Html document for current package
kcl-go doc generate --format html

# Generate Markdown document for specific package
kcl-go doc generate --file-path <package path>

# Generate Markdown document for specific package to a <target directory>
kcl-go doc generate --file-path <package path> --target <target directory>`,
				Flags: []cli.Flag{
					// todo: look for packages recursive
					// todo: package path list
					&cli.StringFlag{
						Name: "file-path",
						Usage: `Relative or absolute path to the KCL package root when running kcl-doc command from
	outside of the KCL package root directory.
	If not specified, docs of all the KCL models under the work directory will be generated.`,
					},
					&cli.BoolFlag{
						Name:  "ignore-deprecated",
						Usage: "do not generate documentation for deprecated schemas",
						Value: false,
					},
					&cli.StringFlag{
						Name:  "format",
						Usage: "The document format to generate. Supported values: markdown, html, openapi",
						Value: string(doc.Markdown),
					},
					&cli.StringFlag{
						Name:  "target",
						Usage: "If not specified, the current work directory will be used. A docs/ folder will be created under the target directory",
					},
				},
				Action: func(context *cli.Context) error {
					opts := doc.GenOpts{
						Path:             context.String("file-path"),
						IgnoreDeprecated: context.Bool("ignore-deprecated"),
						Format:           context.String("format"),
						Target:           context.String("target"),
					}

					genContext, err := opts.ValidateComplete()
					if err != nil {
						fmt.Println(fmt.Errorf("generate failed: %s", err))
					}

					err = genContext.GenDoc()
					if err != nil {
						fmt.Println(fmt.Errorf("generate failed: %s", err))
						return err
					} else {
						fmt.Println(fmt.Sprintf("Generate Complete! Check generated docs in %s", genContext.Target))
						return nil
					}
				},
			},
			{
				Name:  "start",
				Usage: "starts a document website locally",
				Action: func(context *cli.Context) error {
					fmt.Println("not implemented")
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "version",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				_ = cli.ShowCommandHelp(c, c.Command.Name)
				return nil
			}
			arg := c.Args().First()
			if arg == "version" {
				fmt.Println(version)
			}
			return nil
		},
	}
}
