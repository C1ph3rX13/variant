package render

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func (tOpts TmplOpts) TmplRender() error {
	content, err := os.ReadFile(tOpts.TmplFile)
	if err != nil {
		return fmt.Errorf("<os.ReadFile()> err: %s", err)
	}

	// 创建模板对象
	tmpl, err := template.New(filepath.Base(tOpts.TmplFile)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("<template.New()> err: %s", err)
	}

	// 创建输出的 Go 文件
	goFile, err := os.Create(filepath.Join(tOpts.OutputDir, tOpts.OutputGoName))
	if err != nil {
		return fmt.Errorf("<os.Create()> err: %s", err)
	}
	defer goFile.Close()

	// 渲染模板并写入到 Go 文件中
	err = tmpl.Execute(goFile, tOpts.Data)
	if err != nil {
		return fmt.Errorf("<render.Execute()> err: %s", err)
	}

	return nil
}
