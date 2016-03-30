package main

import (
	"fmt"
	tag "gox/gotag"
	"os"
)

func ParseArgs() (args []string) {
	num := len(os.Args)

	if num == 3 {
		args = os.Args[1:]
	} else {
		fmt.Println("the args err")
	}

	return
}
func main() {
	tag.UseFilter(tag.FilterStruct)

	args := ParseArgs()
	if len(args) != 2 {
		return
	}
	if args[0] == "-f" {
		tag.ParseFile(args[1])
	}
	if args[0] == "-d" {
		tag.ParsePkg(args[1])
	}

}
