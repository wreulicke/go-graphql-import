package main

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/wreulicke/go-graphql-import/imports"
)

func processImport(schema string) ([]*imports.ParseResult, error) {
	buf := bytes.NewBufferString(schema)
	sc := bufio.NewScanner(buf)

	var results []*imports.ParseResult
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "# import") {
			p := imports.NewParser(line)
			r, errors := p.Parse()
			if len(errors) > 0 {
				return nil, errors[0]
			}
			results = append(results, r)
		}
	}
	return results, nil
}
