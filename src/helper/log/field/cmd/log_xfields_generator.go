package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"learn/mywebcrawler/helper/log/base"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const field_decl_template = `
type {{.}}Field struct {
	name string
	fieldType FieldType
	value {{if eq . "object"}}interface{}{{else}}{{.}}{{end}}
}

func (field *{{.}}Field) Name() string {
	return field.name
}

func (field *{{.}}Field) Type() FieldType {
	return field.fieldType
}

func (field *{{.}}Field) Value() interface{} {
	return field.value
}

func {{title .}}(name string, value {{if eq . "object"}}interface{}{{else}}{{.}}{{end}}) Field{
	return &{{.}}Field{name: name, fieldType: {{title .}}Type, value: value}
}

`

var (
	inputPath  string
	outputPath string
)

func init() {
	flag.StringVar(&inputPath, "input", "", "The path that contains the target go source files. ")
	flag.StringVar(&outputPath, "output", "", "The path for output the go source file.")
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tlog_xfields_generator [flags] \n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("log_xfields_generator: ")
	flag.Usage = Usage
	flag.Parse()
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Getwd err:%v", err)
	}
	if len(inputPath) == 0 {
		inputPath = currentPath
		log.Printf("WARNING: Not specified the flag named input, use current path '%s'.",
			currentPath)
	} else {
		if !isDir(inputPath) {
			log.Fatalf("ERROR: The input path '%s' is not a ")
		}
	}

	targetFilePath := filepath.Join(inputPath, "field.go")
	prefixes, err := findFieldTypePrefixes(targetFilePath)
	if err != nil {
		log.Fatalf("ERROR: Parse error: %v\n", err)
	}
	var gen Generator
	content, err := gen.generate("field", prefixes...)
	if err != nil {
		log.Fatalf("ERROR: Generate error:%v\n", err)
	}
	outputFilePath := filepath.Join(outputPath, "xfields.go")
	err = ioutil.WriteFile(outputFilePath, content, 0644)
	if err != nil {
		log.Fatalf("ERROR:%v\n", err)
	}
	log.Printf("It has successfully generated a Go source file: %v\n", outputFilePath)
}

func isDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func findFieldTypePrefixes(filePath string) ([]string, error) {
	astFile, err := parser.ParseFile(
		token.NewFileSet(), filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	prefixes := []string{}
	for _, decl := range astFile.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					name := valueSpec.Names[0].Name
					if valueSpec.Type.(*ast.Ident).Name == "FieldType" && name != "UnknownType" {
						prefix := name[:strings.LastIndex(name, "Type")]
						prefixes = append(prefixes, strings.ToLower(prefix))
					}
				}
			}
		}
	}
	return prefixes, nil
}

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) reset() {
	g.buf.Reset()
}

func (g *Generator) generate(pkgName string, prefixes ...string) ([]byte, error) {
	var content []byte
	g.genHeader(pkgName)
	err := g.genFieldDecls(prefixes...)
	if err == nil {
		defer g.buf.Reset()
		content, err = g.format()
	}
	return content, err
}

func (g *Generator) genHeader(pkgName string) {
	g.buf.WriteString("// generated by log_xfields_generator")
	flag.VisitAll(func(fg *flag.Flag) {
		g.buf.WriteString(" -")
		g.buf.WriteString(fg.Name)
		g.buf.WriteString(" ")
		g.buf.WriteString(fg.Value.String())
	})
	g.buf.WriteString("\n// generation time: ")
	g.buf.WriteString(time.Now().Format(base.TIMESTAMP_FORMAT))
	g.buf.WriteString("\n// DO NOT EDIT!!\n")
	g.buf.WriteString("package ")
	g.buf.WriteString(pkgName)
	g.buf.WriteString("\n")
}

func (g *Generator) genFieldDecls(prefixes ...string) error {
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	t := template.Must(template.New("xfield").Funcs(funcMap).Parse(field_decl_template))
	for _, prefix := range prefixes {
		err := t.Execute(&g.buf, prefix)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) format() ([]byte, error) {
	originalSrc := g.buf.Bytes()
	formatedSrc, err := format.Source(originalSrc)
	if err != nil {
		log.Printf("WARNING: Cannot format the generated Go source,"+
			"please build it for the detail, err:%v\n", err)
		return originalSrc, nil
	}
	return formatedSrc, nil
}
