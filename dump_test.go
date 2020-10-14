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

	if !strings.Contains(stdout, "[DEBUG] ") {
		t.Fatalf("invalid stdout: '''\n%v'''", stdout)
	}
	if !strings.Contains(stdout, "dump_test.go:17: dump this bad boy. idx: `1`; strVar: `some data`") {
		t.Fatalf("invalid stdout: '''\n%v'''", stdout)
	}
	if !strings.Contains(stdout, "dump_test.go:21: kv: `map[x:5.6 y:4.5]`; sli: `[true false false]`") {
		t.Fatalf("invalid stdout: '''\n%v'''", stdout)
	}
	if !strings.Contains(stdout, "structVal: `{Data:data string privateValue:[map[k:v] map[a:b]]}`; structVal.privateValue[0][\"k\"]: `v`; structVal.Data: `data string`") {
		t.Fatalf("invalid stdout: '''\n%v'''", stdout)
	}
	if !strings.Contains(stdout, "\"other structure, goes here\", idx: `other structure, goes here`; `1`") {
		t.Fatalf("invalid stdout: '''\n%v'''", stdout)
	}
	if !strings.Contains(stdout, "[DEBUG] target line is invalid. Dump should start with `Dump(` and end with `)`") {
		t.Fatalf("invalid stdout: '''\n%v'''", stdout)
	}
	if !(strings.Contains(stdout, "repeated. i: `0`") && strings.Contains(stdout, "repeated. i: `1`") && strings.Contains(stdout, "repeated. i: `2`")) {
		t.Fatalf("invalid stdout: '''\n%v'''", stdout)
	}
	fmt.Println(stdout)
}

func captureStdout(t *testing.T, f func()) string {
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
