package build

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"variant/log"
)

// Compiler 接口定义
type Compiler interface {
	BuildCompilerArgs() []string
}

// GoCompiler 具体编译器实现
type GoCompiler struct {
	HideConsole bool
	BuildMode   string
	ExeFileName string
	GoFileName  string
}

// NewGoCompiler 创建GoCompiler实例
func NewGoCompiler(c CompileOpts) *GoCompiler {
	return &GoCompiler{
		HideConsole: c.HideConsole,
		BuildMode:   c.BuildMode,
		ExeFileName: c.ExeFileName,
		GoFileName:  c.GoFileName,
	}
}

// BuildCompilerArgs 构建Go编译命令参数
func (gc *GoCompiler) BuildCompilerArgs() []string {
	var compilerArgs []string
	compilerArgs = append(compilerArgs, "-s", "-w")

	if gc.HideConsole {
		compilerArgs = append(compilerArgs, "-H", "windowsgui")
		log.Infof("HideConsole: %v", gc.HideConsole)
	}
	if gc.BuildMode != "" {
		compilerArgs = append(compilerArgs, "-buildmode", gc.BuildMode)
		log.Infof("BuildMode: %v", gc.BuildMode)
	}

	return []string{
		"go",
		"build",
		"-ldflags", strings.Join(compilerArgs, " "),
		"-o", gc.ExeFileName,
		"-trimpath", gc.GoFileName,
	}
}

// GarbleCompiler 具体编译器实现
type GarbleCompiler struct {
	GDebug      bool
	Tiny        bool
	Literals    bool
	GSeed       bool
	HideConsole bool
	BuildMode   string
	ExeFileName string
	GoFileName  string
}

// NewGarbleCompiler 创建GarbleCompiler实例
func NewGarbleCompiler(c CompileOpts) *GarbleCompiler {
	return &GarbleCompiler{
		GDebug:      c.GDebug,
		Tiny:        c.Tiny,
		Literals:    c.Literals,
		GSeed:       c.GSeed,
		HideConsole: c.HideConsole,
		BuildMode:   c.BuildMode,
		ExeFileName: c.ExeFileName,
		GoFileName:  c.GoFileName,
	}
}

// BuildCompilerArgs 构建Garble编译命令参数
func (gc *GarbleCompiler) BuildCompilerArgs() []string {
	var garbleCompilerArgs []string
	garbleCompilerArgs = append(garbleCompilerArgs, "garble")

	if gc.GDebug {
		garbleCompilerArgs = append(garbleCompilerArgs, "-debug")
	}
	if gc.Tiny {
		garbleCompilerArgs = append(garbleCompilerArgs, "-tiny")
	}
	if gc.Literals {
		garbleCompilerArgs = append(garbleCompilerArgs, "-literals")
	}
	if gc.GSeed {
		garbleCompilerArgs = append(garbleCompilerArgs, "-seed=random")
	}

	var goCompilerArgs []string
	if gc.HideConsole {
		goCompilerArgs = append(goCompilerArgs, "-H", "windowsgui")
		log.Infof("HideConsole: %v", gc.HideConsole)
	}
	if gc.BuildMode != "" {
		goCompilerArgs = append(goCompilerArgs, "-buildmode", gc.BuildMode)
		log.Infof("BuildMode: %v", gc.BuildMode)
	}

	var compilerArgs []string
	compilerArgs = append(garbleCompilerArgs, "build", "-ldflags=-s -w")
	compilerArgs = append(compilerArgs, goCompilerArgs...)
	compilerArgs = append(compilerArgs, "-o", gc.ExeFileName, "-trimpath", gc.GoFileName)

	return compilerArgs
}

// compile 使用Compiler接口进行编译
func (c CompileOpts) compile(comp Compiler) error {
	if err := c.formatCode(); err != nil {
		return err
	}

	compilerArgs := comp.BuildCompilerArgs()

	log.Infof("CompilePath: %v", c.CompilePath)
	log.Infof("Compiling: %v", compilerArgs)
	if err := c.execCmd(compilerArgs); err != nil {
		return err
	}
	log.Infof("Compile Succeeded: %s", c.ExeFileName)

	return nil
}

// GoCompile 使用GoCompiler进行编译
func (c CompileOpts) GoCompile() error {
	comp := NewGoCompiler(c)
	return c.compile(comp)
}

// GarbleCompile 使用GarbleCompiler进行编译
func (c CompileOpts) GarbleCompile() error {
	comp := NewGarbleCompiler(c)
	return c.compile(comp)
}

func (c CompileOpts) formatCode() error {
	fmtCmd := []string{"goimports", "-w", c.GoFileName}
	log.Infof("Formatting Code: %v", fmtCmd)
	return c.execCmd(fmtCmd)
}

func (c CompileOpts) execCmd(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if c.CompilePath != "" {
		cmd.Dir = c.CompilePath
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command %q failed: %v", args[0], err)
	}
	return nil
}
