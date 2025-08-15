package daikinac

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func decodeName(name string) (s string) {
	out := make([]byte, 0, len(name))

	for {
		next := strings.IndexByte(name, '%')
		if next == -1 {
			break
		}
		out = append(out, []byte(name[:next])...)
		name = name[next:]

		if len(name) < 3 {
			break
		}

		enc, err := strconv.ParseInt(name[1:3], 16, 0)
		if err != nil {
			out = append(out, '%')
			name = name[1:]
			continue
		}

		// two-digit escapes can generate more than a single char of UTF-8
		var buf [4]byte
		length := utf8.EncodeRune(buf[:], rune(enc))

		out = append(out, buf[:length]...)
		name = name[3:]
	}

	out = append(out, []byte(name)...)
	return string(out)
}
