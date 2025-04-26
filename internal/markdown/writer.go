package markdown

import (
	"os"
	"path/filepath"
	"strings"
)

// SaveMarkdownFile saves the given content to a markdown file under docs/api/<projectName>/<endpoint>.md
func SaveMarkdownFile(sourceDir, endpointName, content string) error {
	// Lấy tên thư mục cuối cùng làm project name
	projectName := filepath.Base(sourceDir)

	// Làm sạch tên file
	fileName := strings.ToLower(strings.ReplaceAll(endpointName, " ", "_")) + ".md"
	dirPath := filepath.Join("docs", projectName)
	filePath := filepath.Join(dirPath, fileName)

	// Tạo thư mục nếu chưa có
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	// Ghi file
	return os.WriteFile(filePath, []byte(content), 0644)
}
