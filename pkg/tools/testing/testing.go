package testing

import (
	"fmt"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type TestOptions struct {
	PkgList   []string
	RunRegRxp string
	FailFast  bool
	// NoRun     bool
}

type TestResult struct {
	Info []TestCaseInfo
}

type TestCaseInfo struct {
	Name       string
	Duration   uint64
	LogMessage string
	ErrMessage string
}

func (t *TestCaseInfo) Pass() bool {
	return t.ErrMessage == ""
}

func (t *TestCaseInfo) Fail() bool {
	return t.ErrMessage != ""
}

// TODO
func (t *TestCaseInfo) Skip() bool {
	return false
}

func (t *TestCaseInfo) Format() string {
	status := fmt.Sprintf("%s: %s (%dms)", t.Name, t.StatusString(), t.Duration/1000)
	if t.LogMessage != "" {
		return fmt.Sprintf("%s\n%s", status, t.LogMessage)
	}
	if t.Fail() {
		return fmt.Sprintf("%s\n%s", status, t.ErrMessage)
	}
	return status
}

func (t *TestCaseInfo) StatusString() string {
	if t.Pass() {
		return "PASS"
	} else if t.Fail() {
		return "FAIL"
	} else if t.Skip() {
		return "SKIPPED"
	}
	return "ERROR"
}

func Test(testOpts *TestOptions, opts ...kcl.Option) (TestResult, error) {
	if testOpts == nil {
		testOpts = &TestOptions{}
	}
	args := kcl.NewOption().Merge(opts...)
	if err := args.Err; err != nil {
		return TestResult{}, err
	}

	svc := kcl.Service()
	resp, err := svc.Test(&gpyrpc.Test_Args{
		ExecArgs:  args.ExecProgram_Args,
		PkgList:   testOpts.PkgList,
		RunRegexp: testOpts.RunRegRxp,
		FailFast:  testOpts.FailFast,
	})
	if err != nil {
		return TestResult{}, err
	}
	var info []TestCaseInfo
	for _, i := range resp.GetInfo() {
		info = append(info, TestCaseInfo{
			Name:       i.Name,
			Duration:   i.Duration,
			LogMessage: i.LogMessage,
			ErrMessage: i.Error,
		})
	}
	return TestResult{
		Info: info,
	}, nil
}
