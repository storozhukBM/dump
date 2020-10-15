package dump_test

import (
	"bytes"
	"fmt"
	"github.com/storozhukBM/dump"
	"io"
	"os"
	"strings"
	"testing"
)

func TestDump(t *testing.T) {
	stdout := captureStdout(t, func() {
		idx := 1
		strVar := "some data"
		dump.Dump("dump this bad boy. ", idx, strVar)

		kv := map[string]float64{"x": 5.6, "y": 4.5}
		sli := []bool{true, false, false}
		dump.Dump(kv, sli)

		type Some struct {
			Data         string
			privateValue []map[string]string
		}
		structVal := Some{
			Data:         "data string",
			privateValue: []map[string]string{{"k": "v"}, {"a": "b"}},
		}
		dump.Dump(structVal, structVal.privateValue[0]["k"], structVal.Data)

		dump.Dump("other structure, goes here", idx)

		dump.Dump()

		dump.Dump(
			idx,
		)

		for i := 0; i < 3; i++ {
			dump.Dump("repeated. ", i)
		}
	})

	mustContain(t, stdout, "[DEBUG] ")
	mustContain(t, stdout, "dump_test.go:17: dump this bad boy. idx: `1`; strVar: `some data`")
	mustContain(t, stdout, "dump_test.go:21: kv: `map[x:5.6 y:4.5]`; sli: `[true false false]`")
	mustContain(t, stdout, "structVal: `{Data:data string privateValue:[map[k:v] map[a:b]]}`; structVal.privateValue[0][\"k\"]: `v`; structVal.Data: `data string`")
	mustContain(t, stdout, "\"other structure, goes here\", idx: `other structure, goes here`; `1`")
	mustContain(t, stdout, "[DEBUG] target line is invalid. Dump should start with `Dump(` and end with `)`")
	mustContain(t, stdout, "repeated. i: `0`")
	mustContain(t, stdout, "repeated. i: `1`")
	mustContain(t, stdout, "repeated. i: `2`")
	fmt.Println(stdout)
}

func mustContain(t *testing.T, target string, expectedPart string) {
	t.Helper()
	if !strings.Contains(target, expectedPart) {
		t.Fatalf("invalid target: '''\n%v'''", target)
	}
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("can't capture stdout: %v", pipeErr)
	}
	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	f()

	closeErr := w.Close()
	if closeErr != nil {
		t.Fatalf("can't close pipe: %v", closeErr)
	}

	var buf bytes.Buffer
	_, copyErr := io.Copy(&buf, r)
	if copyErr != nil {
		t.Fatalf("can't copy captures output: %v", closeErr)
	}
	return buf.String()
}
