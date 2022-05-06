package langserver

import (
	"testing"
)

func TestWordAt(t *testing.T) {
	testData := map[string]struct {
		f      File
		pos    Position
		expect string
	}{
		"empty file": {
			f: File{
				LanguageID: "KCL",
				Text:       "",
				Version:    0,
			},
			pos: Position{
				Line:      0,
				Character: 0,
			},
			expect: "",
		},
		"line out of range": {
			f: File{
				LanguageID: "KCL",
				Text:       "a = 1",
				Version:    0,
			},
			pos: Position{
				Line:      1,
				Character: 0,
			},
			expect: "",
		},
		"column out of range": {
			f: File{
				LanguageID: "KCL",
				Text:       "a = 1",
				Version:    0,
			},
			pos: Position{
				Line:      0,
				Character: 10,
			},
			expect: "",
		},
		"keyword": {
			f: File{
				LanguageID: "KCL",
				Text:       "schema Person:",
				Version:    0,
			},
			pos: Position{
				Line:      0,
				Character: 1,
			},
			expect: "schema",
		},
		"name": {
			f: File{
				LanguageID: "KCL",
				Text:       "a = b",
				Version:    0,
			},
			pos: Position{
				Line:      0,
				Character: 4,
			},
			expect: "b",
		},
		"blank": {
			f: File{
				LanguageID: "KCL",
				Text:       "a =  b",
				Version:    0,
			},
			pos: Position{
				Line:      0,
				Character: 4,
			},
			expect: "  ",
		},
		"punctuation": {
			f: File{
				LanguageID: "KCL",
				Text:       "a+=b",
				Version:    0,
			},
			pos: Position{
				Line:      0,
				Character: 2,
			},
			expect: "+=",
		},
	}
	for name, data := range testData {
		t.Run(name, func(t *testing.T) {
			got := data.f.WordAt(data.pos)
			if data.expect != got {
				t.Fatalf("expect: %s, got: %s", data.expect, got)
			}
		})
	}
}
