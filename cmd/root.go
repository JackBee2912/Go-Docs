package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godocs",
	Short: "godocs - a CLI tool to generate API documentation from Go source code",
	Long:  `godocs helps you scan Go source files and generate markdown API documentation easily.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
