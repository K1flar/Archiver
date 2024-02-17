package main

import (
	"flag"
	"fmt"
)

func main() {
	t := flag.String("oot", "asd", "string")
	flag.Parse()
	fmt.Println(*t)
}
