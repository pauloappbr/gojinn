package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gojinn",
	Short: "Gojinn: The Sovereign Serverless Cloud",
	Long: `Gojinn is a high-performance, secure, and sovereign serverless platform.
It replaces the complexity of AWS Lambda + K8s with a single binary.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
