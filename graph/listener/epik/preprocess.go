package epik

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func preprocess(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		a := strings.SplitN(sc.Text(), " - ", 4)
		if len(a) != 3 {
			return nil, fmt.Errorf("incorrect column num of line: %s", sc.Text())
		}
		if strings.Contains(a[2], " ") {
			if (a[2][0] == '<' && a[2][len(a[2])-1] == '>') ||
				(a[2][0] == '"' && a[2][len(a[2])-1] == '"') {
				// do nothing
			} else {
				a[2] = `"` + a[2] + `"`
			}
		}
		_, err := fmt.Fprintln(&buf, strings.Join(append(a, "."), " "))
		if err != nil {
			return nil, fmt.Errorf("write preprocess(%s) result error: %v", sc.Text(), err)
		}
	}
	return buf.Bytes(), nil
}
