package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/JackBee2912/godocs/internal/gpt"
	"github.com/JackBee2912/godocs/internal/markdown"
	"github.com/JackBee2912/godocs/internal/parser"
	"github.com/JackBee2912/godocs/internal/router"

	"github.com/spf13/cobra"
)

var (
	sourceRoot string
	apiKey     string
	routerFile string
)

var initCmd = &cobra.Command{
	Use:   "i",
	Short: "Initialize documentation generation",
	Run: func(cmd *cobra.Command, args []string) {
		if sourceRoot == "" || apiKey == "" || routerFile == "" {
			log.Fatal("❌ Error: --sourceRoot, --apiKey and --routerFile must be provided.")
		}

		// Parse router file
		apiMappings, err := router.ParseRouterFile(routerFile)
		if err != nil {
			log.Fatalf("❌ Error parsing router file: %v", err)
		}

		subDirs := []string{
			"api", // ✨ bạn sau chỉnh thành param cũng được
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

				// lấy api info nếu có
				apiInfo, ok := apiMappings[fn.Name]

				contextBuilder := &strings.Builder{}
				contextBuilder.WriteString(fmt.Sprintf("Function: %s\n", fn.Name))
				contextBuilder.WriteString(fmt.Sprintf("Comment: %s\n", fn.Comment))

				if ok {
					contextBuilder.WriteString(fmt.Sprintf("HTTP Method: %s\n", apiInfo.Method))
					contextBuilder.WriteString(fmt.Sprintf("Path: %s\n", apiInfo.Path))
					contextBuilder.WriteString(fmt.Sprintf("Authentication Required: %v\n", apiInfo.AuthRequired))
				}

				if fn.RequestModel != "" {
					contextBuilder.WriteString(fmt.Sprintf("Request Model: %s\n", fn.RequestModel))
				}

				if len(fn.ErrorCodes) > 0 {
					contextBuilder.WriteString(fmt.Sprintf("Possible Error Codes: %v\n", fn.ErrorCodes))
				}

				doc, err := gpt.GenerateMarkdownDocumentation(contextBuilder.String(), apiKey)
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
	initCmd.Flags().StringVar(&routerFile, "routerFile", "", "Path to router.go file (required)")
}
