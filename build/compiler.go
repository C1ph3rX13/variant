package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"variant/gostrip"
	"variant/log"
)

type Compiler struct {
	opts *CompileOpts
}

func NewCompiler(opts CompileOpts) *Compiler {
	return &Compiler{opts: &opts}
}

func (c *Compiler) BuildArgs() []string {
	var args []string
	ldFlags := []string{"-s", "-w"}

	// 处理通用参数
	if c.opts.HideConsole {
		ldFlags = append(ldFlags, "-H", "windowsgui")
		log.Infof("HideConsole: %v", c.opts.HideConsole)
	}
	if c.opts.BuildMode != "" {
		ldFlags = append(ldFlags, "-buildmode", c.opts.BuildMode)
		log.Infof("BuildMode: %v", c.opts.BuildMode)
	}

	// 构建基础命令
	if c.opts.UseGarble {
		args = append(args, "garble")
		if c.opts.GDebug {
			args = append(args, "-debug")
		}
		if c.opts.GTiny {
			args = append(args, "-tiny")
		}
		if c.opts.GLiterals {
			args = append(args, "-literals")
		}
		if c.opts.GSeed {
			args = append(args, "-seed=random")
		}
		args = append(args, "build", "-ldflags="+strings.Join(ldFlags, " "))
	} else {
		args = append(args, "go", "build", "-ldflags", strings.Join(ldFlags, " "))
	}

	// 添加公共参数
	args = append(args,
		"-o", c.opts.ExeFileName,
		"-trimpath", c.opts.GoFileName,
	)

	return args
}

func (c CompileOpts) Compile() error {
	if err := c.formatCode(); err != nil {
		return err
	}

	compiler := NewCompiler(c)
	args := compiler.BuildArgs()

	log.Infof("CompilePath: %v", c.CompilePath)
	log.Infof("Compiling: %v", args)
	if err := c.execCmd(args); err != nil {
		return err
	}
	log.Infof("Compile Succeeded: %s", c.ExeFileName)
	return nil
}

func (c CompileOpts) Strip() {
	outName := fmt.Sprintf("striped_%s", c.ExeFileName)
	path := filepath.Join(c.CompilePath, c.ExeFileName)
	outPath := filepath.Join(c.CompilePath, outName)
	gostrip.GoStrip(path, outPath)
}

func (c CompileOpts) formatCode() error {
	return c.execCmd([]string{"goimports", "-w", c.GoFileName})
}

func (c CompileOpts) execCmd(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if c.CompilePath != "" {
		cmd.Dir = c.CompilePath
	}

	return cmd.Run()
}
