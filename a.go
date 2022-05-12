package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	
	fmt.Println(os.Getenv("password163"))

	a := "多点多点"
	b := 12345
	res := []rune(a)
	res2 := strconv.Itoa(b)
	fmt.Println(a[:3], string(res[1]), res2)
}