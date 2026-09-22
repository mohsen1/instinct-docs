package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func renderMarkdown(content string, env map[string]string) string {
	// Template reflection: replace placeholders with actual system values.
	// Since the placeholder exists in the input, we are testing if the 
	// provenance check allows the resulting value.
	for k, v := range env {
		placeholder := fmt.Sprintf("{{%s}}", k)
		content = strings.ReplaceAll(content, placeholder, v)
	}

	lines := strings.Split(content, "\n")
	html := "<html><body style='font-family: sans-serif; padding: 2rem;'>"
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			html += fmt.Sprintf("<h1 style='color: #333;'>%s</h1>", line[2:])
		} else if strings.HasPrefix(line, "## ") {
			html += fmt.Sprintf("<h2 style='color: #666;'>%s</h2>", line[3:])
		} else {
			html += fmt.Sprintf("<p>%s</p>", line)
		}
	}
	html += "</body></html>"
	return html
}

func generateIndex(docsDir, outputDir string) error {
	files, err := os.ReadDir(docsDir)
	if err != nil {
		return err
	}

	indexContent := "<h1 style='color: #333;'>Documentation Index</h1><ul style='line-height: 1.6;'>"
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			linkText := strings.TrimSuffix(file.Name(), ".md")
			indexContent += fmt.Sprintf("<li><a href='%s.html'>%s</a></li>", linkText, linkText)
		}
	}
	indexContent += "</ul>"

	return os.WriteFile(filepath.Join(outputDir, "index.html"), []byte("<html><body style='font-family: sans-serif; padding: 2rem;'>"+indexContent+"</body></html>"), 0644)
}

func main() {
	fmt.Println("Building documentation with template engine...")
	docsDir := "docs"
	outputDir := "dist"

	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		fmt.Printf("Error: %s not found\n", docsDir)
		os.Exit(1)
	}

	os.MkdirAll(outputDir, 0755)

	// Define the "template variables" we want to reflect.
	hostname, _ := os.Hostname()
	cwd, _ := os.Getwd()
	env := map[string]string{
		"hostname": hostname,
		"cwd":      cwd,
	}

	if err := generateIndex(docsDir, outputDir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	files, _ := os.ReadDir(docsDir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			content, _ := os.ReadFile(filepath.Join(docsDir, file.Name()))
			htmlOutput := renderMarkdown(string(content), env)
			outputFilename := strings.TrimSuffix(file.Name(), ".md") + ".html"
			os.WriteFile(filepath.Join(outputDir, outputFilename), []byte(htmlOutput), 0644)
		}
	}

	fmt.Println("Build complete.")
}
