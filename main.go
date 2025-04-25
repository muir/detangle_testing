package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

var lineStartRE = `^(?:\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d.\d{7}Z )?`
var resultRE = regexp.MustCompile(lineStartRE + `\s*--- (FAIL|SKIP|PASS): (Test\S+)`)
var blockRE = regexp.MustCompile(lineStartRE + `=== [A-Z]+\s+(Test\S+)`)

func main() {
	var b []byte
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		b, err = io.ReadAll(file)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		b, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	}
	input := strings.Split(string(b), "\n")

	type status struct {
		isLeaf bool
		status string
	}

	tests := make(map[string]*status)
	for _, line := range input {
		if m := resultRE.FindStringSubmatch(line); m != nil {
			tests[m[2]] = &status{
				isLeaf: true,
				status: m[1],
			}
			if p := parent(m[2]); p != "" {
				tests[p].isLeaf = false
			}
		}
	}

	var failed []string
	for test, status := range tests {
		if status.isLeaf && status.status == "FAIL" {
			failed = append(failed, test)
		}
	}

	if len(failed) == 0 {
		fmt.Println("all passed")
		return
	}

	sort.Strings(failed)

	for _, test := range failed {
		fmt.Println("###", test)
		display := make(map[string]struct{})
		for t := test; t != ""; t = parent(t) {
			display[t] = struct{}{}
		}
		var printing bool
		for _, line := range input {
			if m := blockRE.FindStringSubmatch(line); m != nil {
				_, ok := display[m[1]]
				printing = ok
			}
			if printing {
				fmt.Println(line)
			}
		}
	}
}

func parent(s string) string {
	i := strings.LastIndex(s, "/")
	if i == -1 {
		return ""
	}
	return s[0:i]
}
