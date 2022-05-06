// Copyright 2021 The KCL Authors. All rights reserved.

package ktest

import (
	"bytes"
	_ "embed"
	"log"
	"text/template"
)

//go:embed __kcl_test_main.tmpl.k
var __kcl_test_main_k string

func genTestMainFile(testSchemaNames []string) string {
	t := template.Must(template.New("").Parse(__kcl_test_main_k))

	var buf bytes.Buffer
	if err := t.Execute(&buf, testSchemaNames); err != nil {
		log.Panic(err)
	}

	return buf.String()
}
