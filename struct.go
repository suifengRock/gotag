package gotag

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

const (
	CommentToken = "//"
	StructToken  = "@Tag:"
)

var (
	genAfterDelToken = false
)

func SetGenAfterDelToken(isDel bool) {
	genAfterDelToken = isDel
}

func FilterStruct(f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.StructType:
			parseStruct(t)
			return false
		}
		return true
	})
}

func parseStruct(n *ast.StructType) {
	if n.Fields.NumFields() == 0 {
		return
	}
	first := n.Fields.List[0]
	if first.Doc == nil {
		return
	}
	tagMap, tagIndex := getTags(first.Doc.List)
	if len(tagMap) == 0 {
		return
	}
	if genAfterDelToken {
		newList := first.Doc.List[:tagIndex]
		newList = append(newList, first.Doc.List[tagIndex+1:]...)
		first.Doc.List = newList
	}
	for _, field := range n.Fields.List {
		if field.Tag == nil {
			name := field.Names[0].String()
			pos := field.Type.Pos() + 1
			field.Tag = newTag(tagMap, pos, name)
		} else {
			for name, _ := range tagMap {
				handle := GetTagHandle(name)
				if handle == nil {
					continue
				}
				if strings.Contains(field.Tag.Value, fmt.Sprintf("%s:\",", name)) {
					fieldName := field.Names[0].String()
					repStr := fmt.Sprintf("%s:\"%s", name, handle(fieldName))
					oldStr := fmt.Sprintf("%s:\"", name)
					field.Tag.Value = strings.Replace(field.Tag.Value, oldStr, repStr, 1)
				} else if !strings.Contains(field.Tag.Value, fmt.Sprintf("%s:\"", name)) {
					fieldName := field.Names[0].String()
					repStr := fmt.Sprintf("`%s:\"%s\"", name, handle(fieldName))
					oldStr := "`"
					field.Tag.Value = strings.Replace(field.Tag.Value, oldStr, repStr, 1)
				}
			}
		}
	}
}

func newTag(tagMap map[string]bool, pos token.Pos, fieldName string) *ast.BasicLit {
	tag := new(ast.BasicLit)
	tag.ValuePos = pos
	tag.Kind = token.STRING
	comment := ""
	for name, _ := range tagMap {
		handle := GetTagHandle(name)
		if handle == nil {
			continue
		}
		comment = fmt.Sprintf("%s%s:\"%s\" ", comment, name, handle(fieldName))
	}
	tag.Value = fmt.Sprintf("`%s`", comment)
	return tag
}

func getTags(list []*ast.Comment) (map[string]bool, int) {
	tagMap := map[string]bool{}
	for i, c := range list {
		if strings.Contains(c.Text, StructToken) {
			tagComment := strings.Replace(c.Text, CommentToken, "", 1)
			tagComment = strings.Replace(tagComment, StructToken, "", 1)
			tagComment = strings.TrimSpace(tagComment)
			tags := strings.Split(tagComment, " ")
			for _, tag := range tags {
				tagMap[tag] = true
			}
			return tagMap, i
		}
	}
	return tagMap, 0
}
