package scripts

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func IsMd5Text(s string) bool {
	s = strings.TrimSpace(s)
	matched, err := regexp.MatchString(`^[a-f0-9]{32}$`, s)
	if err != nil {
		panic(err)
	}
	return matched
}

func MD5File(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
