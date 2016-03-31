# gotag
golang auto generate struct tag

## example
run this main.go

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

then you can see the user struct

	type User struct {
		// @Tag: json bson
		UserName string `json:"user_name" bson:"user_name"`
		Age      int 	`json:"age" bson:"age"`
	}

## notice
1. you need write ```@Tag:``` comment on the first field of the struct
2. after the ```@Tag:```, write you need the tag
3. see the struct user in example 

## extend
1. add or overloaded tag's translate handle: [see this](https://github.com/suifengRock/gotag/blob/master/tag.go)
2. add or overloaded node's parse handle: [see this](https://github.com/suifengRock/gotag/blob/master/parser.go)
3. use middleware to filter and parse node 

## support
If you do have a contribution for the package feel free to put up a Pull Request or open Issue.
