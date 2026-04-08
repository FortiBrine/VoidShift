//go:build !linux

package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Println("❌ Unsupported OS:", runtime.GOOS)
	fmt.Println("This program supports only Linux")
	os.Exit(1)
}
