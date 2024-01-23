package render

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func TmplRender(opts TmplOpts) error {
	content, err := os.ReadFile(opts.TmplFile)
	if err != nil {
		return fmt.Errorf("<os.ReadFile()> err: %s", err)
	}

	// 创建模板对象
	tmpl, err := template.New(filepath.Base(opts.TmplFile)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("<template.New()> err: %s", err)
	}

	// 创建输出的 Go 文件
	goFile, err := os.Create(filepath.Join(opts.OutputDir, opts.OutputGoName))
	if err != nil {
		return fmt.Errorf("<os.Create()> err: %s", err)
	}
	defer goFile.Close()

	// 渲染模板并写入到 Go 文件中
	err = tmpl.Execute(goFile, opts.Data)
	if err != nil {
		return fmt.Errorf("<render.Execute()> err: %s", err)
	}

	return nil
}
