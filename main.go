package main

import (
	"knork.org/compiler/pkg/reverse_polish"
)

func main() {
	reverse_polish.Compile("./inputs/input1.txt", "go_elf")

}
