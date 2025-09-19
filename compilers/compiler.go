package compilers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"variant/log"
	"variant/strip"
)

func (b *BuildCtx) BuildArgs() []string {
	var args []string
	ldFlags := []string{"-s", "-w"}

	// 处理通用参数
	if b.HideConsole {
		ldFlags = append(ldFlags, "-H", "windowsgui")
		log.Infof("HideConsole: %v", b.HideConsole)
	}
	if b.BuildMode != "" {
		ldFlags = append(ldFlags, "-buildmode", b.BuildMode)
		log.Infof("BuildMode: %v", b.BuildMode)
	}

	// 构建基础命令
	if b.Obfuscate {
		args = append(args, "garble")
		if b.Debug {
			args = append(args, "-debug")
		}
		if b.Tiny {
			args = append(args, "-tiny")
		}
		if b.Literals {
			args = append(args, "-literals")
		}
		if b.Seed {
			args = append(args, "-seed=random")
		}
		args = append(args, "build", "-ldflags="+strings.Join(ldFlags, " "))
	} else {
		args = append(args, "go", "build", "-ldflags", strings.Join(ldFlags, " "))
	}

	// 添加公共参数
	args = append(args,
		"-o", b.OutExeFile,
		"-trimpath", b.SrcGoFile,
	)

	return args
}

func (b *BuildCtx) Compile() error {
	if err := b.formatCode(); err != nil {
		return err
	}

	args := b.BuildArgs()

	log.Infof("PatchDir: %v", b.WorkingDir)
	log.Infof("Compiling: %v", args)
	if err := b.execCmd(args); err != nil {
		return err
	}
	
	log.Infof("Compile Succeeded: %s", b.OutExeFile)
	return nil
}

func (b *BuildCtx) Strip() {
	outName := fmt.Sprintf("striped_%s", b.OutExeFile)
	path := filepath.Join(b.WorkingDir, b.OutExeFile)
	outPath := filepath.Join(b.WorkingDir, outName)
	strip.GoStrip(path, outPath)
}

func (b *BuildCtx) formatCode() error {
	return b.execCmd([]string{"goimports", "-w", b.SrcGoFile})
}

func (b *BuildCtx) execCmd(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if b.WorkingDir != "" {
		cmd.Dir = b.WorkingDir
	}

	return cmd.Run()
}
