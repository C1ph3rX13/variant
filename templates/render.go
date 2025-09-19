package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"variant/compilers"
)

func (t *TmplCtx) TmplRender() (string, error) {

	tmpl := template.Must(template.ParseGlob("templates/tmpls/*.tmpl"))

	p, err := InitOutputPath()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(p, compilers.RandomGoFile())
	goFile, fErr := os.Create(filename)
	if fErr != nil {
		return "", fmt.Errorf("create go file failed: %w", fErr)
	}
	defer goFile.Close()

	eErr := tmpl.ExecuteTemplate(goFile, t.Mode(), t)
	if eErr != nil {
		return "", fmt.Errorf("execute template failed: %w", eErr)
	}

	return filename, nil
}

// Mode 根据构建的结构体选择对应的模板渲染
func (t *TmplCtx) Mode() string {
	switch {
	case t.Local != nil:
		return "Local"
	case t.Dynamic != nil:
		return "Dynamic"
	case t.Pokemon != nil:
		return "Pokemon"
	default:
		return ""
	}
}

// InitOutputPath 初始化并返回输出目录路径。
func InitOutputPath() (string, error) {
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return "", fmt.Errorf("failed to get working directory: %w", wdErr)
	}

	outPath := filepath.Join(wd, "output")
	if err := os.MkdirAll(outPath, 0750); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	return outPath, nil
}
