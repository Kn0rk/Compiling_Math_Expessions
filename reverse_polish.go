package main

import "fmt"
import "knork.org/compiler/pkg/reverse_polish"

func main() {
	reverse_polish.ParseFile("./inputs/input1.txt")
	fmt.Println("ELF file 'minimal_elf' created.")

}
