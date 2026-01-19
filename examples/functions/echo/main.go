package main

import (
	"io"
	"os"
)

func main() {
	data, _ := io.ReadAll(os.Stdin)

	os.Stdout.WriteString("👋 Hello from Gojinn! You sent: ")
	os.Stdout.Write(data)
}
