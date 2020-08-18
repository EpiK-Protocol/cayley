package epik

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/epik-protocol/epik-gateway-backend/clog"
)

var escaper = strings.NewReplacer(
	"\\", "\\\\",
	"\"", "\\\"",
	"\n", "\\n",
	"\r", "\\r",
	"\t", "\\t",
)

func preprocess(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		a := strings.SplitN(sc.Text(), " - ", 4)
		if len(a) != 3 && len(a) != 4 {
			clog.Errorf("illegal line: %s", sc.Text())
			continue
		}
		for i := range a {
			if notRaw(a[i]) {
				continue
			}
			a[i] = `"` + escaper.Replace(a[i]) + `"`
		}
		_, err := fmt.Fprintln(&buf, strings.Join(append(a, "."), " "))
		if err != nil {
			return nil, fmt.Errorf("write preprocess(%s) result error: %v", sc.Text(), err)
		}
	}
	return buf.Bytes(), nil
}

func notRaw(v string) bool {
	if len(v) > 2 {
		if (v[0] == '<' && v[len(v)-1] == '>') || (v[0] == '"' && v[len(v)-1] == '"') {
			return true
		} else if v[:2] == "_:" {
			return true
		} else if i := strings.Index(v, `"^^<`); i > 0 && v[0] == '"' && v[len(v)-1] == '>' {
			return true
		} else if i := strings.Index(v, `"@`); i > 0 && v[0] == '"' && v[len(v)-1] != '"' {
			return true
		}
	}
	return false
}
