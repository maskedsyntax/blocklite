package main

import (
	"blocklite/utils"
	"fmt"
)

func main() {
	fmt.Printf("%x\n", utils.SHA256("Hello, World!"))
}
