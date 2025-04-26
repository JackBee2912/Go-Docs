package cmd

import (
	"fmt"
	"godocs/internal/gpt"
	"godocs/internal/markdown"
	"godocs/internal/parser"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	sourceRoot string
	apiKey     string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize documentation generation",
	Run: func(cmd *cobra.Command, args []string) {
		if sourceRoot == "" || apiKey == "" {
			log.Fatal("❌ Error: Both --sourceRoot and --apiKey must be provided.")
		}

		subDirs := []string{
			"api", // ✨ bạn muốn có thể cho param nhiều subDirs cũng ok sau
		}

		for _, subDir := range subDirs {
			sourceDir := filepath.Join(sourceRoot, subDir)

			funcs, err := parser.ParseFunctionsFromDir(sourceDir)
			if err != nil {
				log.Printf("❌ Parse error in %s: %v", sourceDir, err)
				continue
			}

			for _, fn := range funcs {
				fnNameLower := strings.ToLower(fn.Name)

				if strings.HasPrefix(fnNameLower, "new") {
					fmt.Printf("⏩ Skipping function (prefix new): %s\n", fn.Name)
					continue
				}

				fmt.Printf("🚀 Processing %s in %s...\n", fn.Name, sourceDir)
				context := fmt.Sprintf("Function: %s\nComment: %s\n", fn.Name, fn.Comment)

				doc, err := gpt.GenerateMarkdownDocumentation(context, apiKey)
				if err != nil {
					log.Printf("⚠️ GPT Error (%s): %v\n", fn.Name, err)
					continue
				}

				err = markdown.SaveMarkdownFile(sourceDir, fn.Name, doc)
				if err != nil {
					log.Printf("⚠️ File Write Error (%s): %v\n", fn.Name, err)
					continue
				}

				fmt.Printf("✅ Done: %s in %s\n", fn.Name, sourceDir)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&sourceRoot, "sourceRoot", "", "Root directory containing source code (required)")
	initCmd.Flags().StringVar(&apiKey, "apiKey", "", "API key for GPT generation (required)")
}
