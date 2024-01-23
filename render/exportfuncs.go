package render

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"variant/log"
)

type FunctionType int

const (
	Encrypt FunctionType = iota
	Decrypt
)

// GetExportedFuncsFromFolder 获取指定文件夹中所有Go文件，并提取每个文件中的导出函数
func GetExportedFuncsFromFolder(folderPath string, functionType FunctionType) ([]string, error) {
	fset := token.NewFileSet() // 创建一个新的文件集合，用于保存解析产生的信息
	pkgs, err := parser.ParseDir(fset, folderPath, func(info os.FileInfo) bool {
		// 只解析Go源代码文件
		return !info.IsDir() && filepath.Ext(info.Name()) == ".go"
	}, parser.ParseComments) // 解析指定文件夹中的Go代码，并保留注释信息
	if err != nil {
		return nil, err
	}

	funcPrefix, err := getFuncPrefix(functionType)
	if err != nil {
		log.Fatal(err)
	}

	exportedFunctions := make([]string, 0)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			funcs := getExportedFuncsFromFile(file, funcPrefix)     // 获取当前文件中的导出函数
			exportedFunctions = append(exportedFunctions, funcs...) // 将导出函数添加到结果列表中
		}
	}

	return exportedFunctions, nil
}

// getExportedFuncsFromFile 获取指定Go源文件中的导出函数
func getExportedFuncsFromFile(file *ast.File, funcPrefix string) []string {
	functions := make([]string, 0)

	for _, decl := range file.Decls {
		// 只提取函数声明
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// 只提取导出函数
		if funcDecl.Name.IsExported() {
			if strings.Contains(funcDecl.Name.Name, funcPrefix) {
				functions = append(functions, funcDecl.Name.Name)
			}
			continue
		}
	}

	return functions
}

// getFuncPrefix 根据函数类型获取函数名前缀
func getFuncPrefix(functionType FunctionType) (string, error) {
	switch functionType {
	case Encrypt:
		return "Encrypt", nil
	case Decrypt:
		return "Decrypt", nil
	default:
		return "", errors.New("please select encryption or decryption functions")
	}
}
