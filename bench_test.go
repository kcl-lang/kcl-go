// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm

// KCLVM_GO_API_TEST_NUM_CPU=1 go test -bench=.

import (
	"os"
	"sync"
	"testing"
)

func BenchmarkValidateCode_hello_1(b *testing.B) {
	const N = 1
	tBenchValidateCode(b, "./testdata/vet/hello.k.json", "./testdata/vet/hello.k", N)
}
func BenchmarkValidateCode_hello_4(b *testing.B) {
	const N = 4
	tBenchValidateCode(b, "./testdata/vet/hello.k.json", "./testdata/vet/hello.k", N)
}
func BenchmarkValidateCode_hello_8(b *testing.B) {
	const N = 8
	tBenchValidateCode(b, "./testdata/vet/hello.k.json", "./testdata/vet/hello.k", N)
}

func BenchmarkValidateCode_sample_1(b *testing.B) {
	const N = 1
	tBenchValidateCode(b, "./testdata/vet/sample.k.json", "./testdata/vet/sample.k", N)
}
func BenchmarkValidateCode_sample_4(b *testing.B) {
	const N = 4
	tBenchValidateCode(b, "./testdata/vet/sample.k.json", "./testdata/vet/sample.k", N)
}
func BenchmarkValidateCode_sample_8(b *testing.B) {
	const N = 8
	tBenchValidateCode(b, "./testdata/vet/sample.k.json", "./testdata/vet/sample.k", N)
}

func tBenchValidateCode(b *testing.B, datafile, codefile string, N int) {
	var (
		data = tLoadFile(b, datafile)
		code = tLoadFile(b, codefile)
	)

	b.ResetTimer()

	var limit = make(chan struct{}, N)
	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			limit <- struct{}{}
			defer func() { <-limit }()

			if _, err := ValidateCode(data, code, nil); err != nil {
				_ = err // ignore error
			}
		}()
	}
	wg.Wait()
}

func tLoadFile(tb testing.TB, path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		tb.Fatal(err)
	}
	return string(data)
}
