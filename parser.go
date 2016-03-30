package gotag

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"time"
)

type FilterNode func(*ast.File)

type nodeFilter struct {
	Stack []FilterNode
}

func (n *nodeFilter) Run(f *ast.File) {
	for _, filter := range n.Stack {
		filter(f)
	}
}

func (n *nodeFilter) UseFilter(f FilterNode) {
	if n.Stack == nil {
		n.Stack = make([]FilterNode, 0)
	}
	n.Stack = append(n.Stack, f)
}

func (n *nodeFilter) Clear() {
	n.Stack = make([]FilterNode, 0)
}

type Parser struct {
	*nodeFilter
}

func NewParser() *Parser {
	p := &Parser{nodeFilter: new(nodeFilter)}
	p.nodeFilter = new(nodeFilter)
	return p
}

func (p *Parser) ParseFile(fileName string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	p.Run(f)
	writeFile(fset, f, fileName)
}

func (p *Parser) ParsePkg(path string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	counter := len(pkgs)
	if counter == 0 {
		return
	}
	chn := make(chan int, 0)
	for _, pkg := range pkgs {
		go func(fset *token.FileSet, files map[string]*ast.File) {
			for path, file := range files {
				go func(fset *token.FileSet, f *ast.File, path string) {
					fmt.Println(path)

					p.Run(f)
					writeFile(fset, f, path)
				}(fset, file, path)
			}
			time.Sleep(500 * time.Millisecond)
			chn <- 1
		}(fset, pkg.Files)
	}
	for counter > 0 {
		<-chn
		counter -= 1
	}
}

var parserSet = NewParser()

func ParseFile(fileName string) {
	parserSet.ParseFile(fileName)
}

func ParsePkg(path string) {
	parserSet.ParsePkg(path)
}

func UseFilter(f FilterNode) {
	parserSet.UseFilter(f)
}

func ResetFilter() {
	parserSet.Clear()
}

func writeFile(fset *token.FileSet, f *ast.File, fileName string) {
	write, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer write.Close()
	w := bufio.NewWriter(write)
	err = format.Node(w, fset, f)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Flush()
}
