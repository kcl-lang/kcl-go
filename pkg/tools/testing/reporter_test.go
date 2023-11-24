package testing

import (
	"bytes"
	"testing"
)

func TestPrettyReporter(t *testing.T) {
	var buf bytes.Buffer
	result := TestResult{
		Info: []TestCaseInfo{
			{
				Name:     "test_foo",
				Duration: 1024,
			},
			{
				Name:       "test_bar",
				Duration:   2048,
				ErrMessage: "Error: assert failed",
			},
		},
	}

	r := PrettyReporter{
		Output:  &buf,
		Verbose: false,
	}
	if err := r.Report(&result); err != nil {
		t.Fatal(err)
	}

	exp := `test_foo: PASS (1ms)
test_bar: FAIL (2ms)
Error: assert failed
--------------------------------------------------------------------------------
PASS: 1/2
FAIL: 1/2
`

	if exp != buf.String() {
		t.Fatalf("Expected:\n\n%v\n\nGot:\n\n%v", exp, buf.String())
	}
}
