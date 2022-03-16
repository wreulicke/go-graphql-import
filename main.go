package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

func LoadSchema(filePath string, visited map[string]*ast.Schema) (*ast.Schema, error) {
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	source := ast.Source{
		Input: string(bs),
	}
	schema, gqlErr := gqlparser.LoadSchema(&source)
	if gqlErr != nil {
		return nil, fmt.Errorf("schema parsing is failed. %w", gqlErr)
	}
	visited[filePath] = schema

	results, err := processImport(string(bs))
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		rel := filepath.Join(filepath.Dir(filePath), r.From)
		if _, ok := visited[rel]; ok {
			continue
		}
		schema, err := LoadSchema(rel, visited)
		if err != nil {
			return nil, err
		}
		visited[rel] = schema
	}

	return schema, nil
}

func mainInternal() error {
	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		return err
	}
	schema, err := LoadSchema(path, make(map[string]*ast.Schema))
	if err != nil {
		return err
	}
	fmt.Println(schema)
	return nil
}

func main() {
	err := mainInternal()
	if err != nil {
		log.Fatal(err)
	}
}
