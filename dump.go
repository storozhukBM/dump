package dump

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func Dump(values ...interface{}) {
	if len(values) == 0 {
		return
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("[DEBUG] can't capture stack for dump\n")
		return
	}
	targetFile, openErr := os.Open(file)
	if openErr != nil {
		fmt.Printf("[DEBUG] can't open file: %v\n", file)
		return
	}
	defer func() {
		err := targetFile.Close()
		if err != nil {
			fmt.Printf("[DEBUG] can't close file: %v; err: %v\n", file, err)
		}
	}()

	scanner := bufio.NewScanner(targetFile)
	lineCnt := 0
	targetLine := ""
	for scanner.Scan() {
		lineCnt++
		if lineCnt != line {
			continue
		}
		fileLine := strings.Trim(scanner.Text(), ` 	`)
		dumpStartIdx := strings.Index(fileLine, "Dump(")
		dumpEndIdx := strings.LastIndex(fileLine, ")")
		if dumpStartIdx < 0 || dumpEndIdx < 0 {
			fmt.Printf("[DEBUG] target line is invalid. Dump should start with `Dump(` and end with `)`: %v\n", fileLine)
			return
		}
		targetLine = fileLine[dumpStartIdx+5 : dumpEndIdx]
		break
	}
	dumpVariables := strings.Split(targetLine, ", ")
	if len(dumpVariables) != len(values) {
		buff := &bytes.Buffer{}
		_, _ = fmt.Fprintf(buff, "[DEBUG] %v:%v: ", file, line)
		_, _ = fmt.Fprintf(buff, "%v: ", targetLine)
		for idx, val := range values {
			_, _ = fmt.Fprintf(buff, "`%+v`", val)
			if idx < len(values)-1 {
				_, _ = fmt.Fprintf(buff, "; ")
			}
		}
		_, _ = fmt.Fprintf(buff, "\n")
		fmt.Print(buff.String())
		return
	}

	buff := &bytes.Buffer{}
	_, _ = fmt.Fprintf(buff, "[DEBUG] %v:%v: ", file, line)
	for idx, variable := range dumpVariables {
		isStringLiteral := strings.HasPrefix(variable, `"`) && strings.HasSuffix(variable, `"`)
		isStringLiteral = isStringLiteral || strings.HasPrefix(variable, "`") && strings.HasSuffix(variable, "`")
		if isStringLiteral {
			_, _ = fmt.Fprintf(buff, "%v", variable[1:len(variable)-1])
		} else {
			_, _ = fmt.Fprintf(buff, "%v: `%+v`", variable, values[idx])
			if idx < len(values)-1 {
				_, _ = fmt.Fprintf(buff, "; ")
			}
		}
	}
	_, _ = fmt.Fprintf(buff, "\n")
	fmt.Print(buff.String())
}
