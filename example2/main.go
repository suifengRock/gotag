package main

import (
	"fmt"
	tag "gotag"
)

type User struct {
	// @Tag: json bson
	UserName string
	Age      int
}

type Book struct {
	Name   string
	Author string
}

func main() {
	tag.UseFilter(tag.FilterStruct)

	tag.ParseFile("src/gotag/example2/main.go")

	// you can also parae package
	//tag.ParsePkg("src/gotag/example2/")

}
