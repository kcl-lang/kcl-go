package langserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf16"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentReference(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params ReferenceParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return h.findRefs(params.TextDocument.URI, &params.TextDocumentPositionParams)
}

func (h *langHandler) findRefs(uri DocumentURI, params *TextDocumentPositionParams) ([]Location, error) {
	f, ok := h.files[uri]
	if !ok {
		return nil, fmt.Errorf("document not found: %v", uri)
	}
	word, err := h.wordAtPos(f, params)
	if err != nil {
		return nil, err
	}
	if isBuiltinName(word) {
		return nil, nil
	}
	fname, err := fromURI(uri)
	if err != nil {
		return nil, err
	}
	fname = filepath.ToSlash(fname)
	if runtime.GOOS == "windows" {
		fname = strings.ToLower(fname)
	}

	base := h.rootPath
	if base == "" {
		return nil, nil
	}
	locations, err := h.findTag(base, word)
	if err != nil {
		return locations, err
	}
	var result []Location
	for _, location := range locations {
		if location.Range.Start.Line == params.Position.Line && location.Range.Start.Character <= params.Position.Character && location.Range.End.Character >= params.Position.Character {
			// ignore the ref itself
			continue
		}
		result = append(result, location)
	}
	return result, nil
}

func (h *langHandler) findTag(pathName string, tag string) ([]Location, error) {
	if strings.TrimSpace(tag) == "" {
		return nil, nil
	}
	fileOrDir, err := os.Stat(pathName)
	if err != nil {
		return nil, err
	}
	var files []string
	var locations []Location
	switch mode := fileOrDir.Mode(); {
	case mode.IsDir():
		err := filepath.Walk(pathName,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() && strings.HasSuffix(path, ".k") {
					files = append(files, path)
				}
				return nil
			})
		if err != nil {
			return nil, err
		}
	case mode.IsRegular():
		files = append(files, pathName)
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		fullPath := filepath.Clean(file)
		b, err := ioutil.ReadFile(fullPath)
		if err != nil {
			continue
		}
		lines := strings.Split(string(b), "\n")
		for i, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, "#") {
				continue
			}
			if strings.Contains(trimmedLine, tag) {
				startIndex := strings.Index(trimmedLine, tag)
				endIndex := startIndex + len(tag) - 1
				if startIndex != 0 {
					prev := []rune(trimmedLine)[startIndex-1]
					if unicode.IsLetter(prev) || unicode.IsDigit(prev) || string(prev) == "_" {
						continue
					}
				}
				if endIndex < len(trimmedLine)-1 {
					next := []rune(trimmedLine)[endIndex+1]
					if unicode.IsLetter(next) || unicode.IsDigit(next) || string(next) == "_" {
						continue
					}
				}

				locations = append(locations, Location{
					URI: toURI(fullPath),
					Range: Range{
						Start: Position{Line: i, Character: strings.Index(line, tag)},
						End:   Position{Line: i, Character: strings.Index(line, tag) + len(tag)},
					},
				})
			}
		}
		f.Close()
	}

	return locations, nil
}

func isBuiltinName(word string) bool {
	var builtinNames = []string{
		"import", "as", "schema", "mixin", "protocol", "relaxed", "check", "for",
		"assert", "if", "elif", "else", "or", "and", "not", "in", "is", "final", "lambda",
		"all", "any", "filter", "map",
		"any", "str", "int", "float", "bool",
		"True", "False", "None", "Undefined",
	}
	for _, name := range builtinNames {
		if name == word {
			return true
		}
	}
	return false
}

func (h *langHandler) wordAtPos(f *File, params *TextDocumentPositionParams) (string, error) {
	lines := strings.Split(f.Text, "\n")
	if params.Position.Line < 0 || params.Position.Line > len(lines) {
		return "", fmt.Errorf("invalid position: %v", params.Position)
	}
	chars := utf16.Encode([]rune(lines[params.Position.Line]))
	if params.Position.Character < 0 || params.Position.Character > len(chars) {
		return "", fmt.Errorf("invalid position: %v", params.Position)
	}
	if chars[0] == '#' {
		return "", nil
	}
	return wordAtIndex(params.Position.Character, chars), nil
}

func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r)
}

func isWordLeadingChar(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func wordAtIndex(index int, chars []uint16) string {
	if !isWordChar(rune(chars[index])) {
		return ""
	}
	prevWord := false
	startPos := index
	endPos := len(chars)

	for i, char := range chars {
		currentWord := isWordChar(rune(char))
		currentLeading := isWordLeadingChar(rune(char))
		if i <= index {
			if currentLeading && !prevWord {
				startPos = i
			}
			if !currentWord {
				startPos = index
			}
		}
		if i > index && !currentWord {
			endPos = i
			break
		}
		prevWord = currentWord
	}
	if unicode.IsDigit(rune(chars[startPos])) {
		return ""
	}
	return string(utf16.Decode(chars[startPos:endPos]))
}
