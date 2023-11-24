package testing

import (
	"fmt"
	"io"
	"strings"
)

func DefaultReporter(output io.Writer) Reporter {
	return &PrettyReporter{
		Output: output,
	}
}

// Reporter defines the interface for reporting test results.
type Reporter interface {
	// Report is called with a channel that will contain test results.
	Report(result *TestResult) error
}

// PrettyReporter reports test results in a simple human readable format.
type PrettyReporter struct {
	Output  io.Writer
	Verbose bool
}

// Report prints the test report to the reporter's output.
func (r PrettyReporter) Report(result *TestResult) error {
	if result == nil {
		return nil
	}
	dirty := false
	var pass, fail, skip int

	var failures []*TestCaseInfo
	for _, info := range result.Info {
		if info.Pass() {
			pass++
		} else if info.Skip() {
			skip++
		} else if info.ErrMessage != "" {
			fail++
			failures = append(failures, &info)
		}
	}

	for _, info := range result.Info {
		fmt.Fprintln(r.Output, info.Format())
	}

	if dirty {
		r.hl()
	}

	total := pass + fail + skip

	r.hl()

	if pass != 0 {
		fmt.Fprintln(r.Output, "PASS:", fmt.Sprintf("%d/%d", pass, total))
	}

	if fail != 0 {
		fmt.Fprintln(r.Output, "FAIL:", fmt.Sprintf("%d/%d", fail, total))
	}

	if skip != 0 {
		fmt.Fprintln(r.Output, "SKIPPED:", fmt.Sprintf("%d/%d", skip, total))
	}

	return nil
}

func (r PrettyReporter) hl() {
	fmt.Fprintln(r.Output, strings.Repeat("-", 80))
}
