package build

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"variant/log"
)

func (c CompileOpts) GoCompile() error {
	return c.compile("go")
}

func (c CompileOpts) GarbleCompile() error {
	return c.compile("garble")
}

func (c CompileOpts) compile(cmd string) error {
	// 预先格式化代码
	if err := c.formatCode(); err != nil {
		return err
	}

	// 根据条件构建编译命令
	cmdArgs := c.buildCmdArgs(cmd)

	// 执行编译
	log.Infof("CompilePath: %v", c.CompilePath)
	log.Infof("Compiling: %v", cmdArgs)
	if err := c.execCmd(cmdArgs); err != nil {
		return err
	}
	log.Infof("Compile Succeeded: %s", c.ExeFileName)

	return nil
}

func (c CompileOpts) formatCode() error {
	fmtCmd := []string{"goimports", "-w", c.GoFileName}
	log.Infof("Formatting Code: %v", fmtCmd)

	return c.execCmd(fmtCmd)
}

// buildCmdArgs 根据 CompileOpts 参数构建编译命令的参数列表
func (c CompileOpts) buildCmdArgs(cmd string) []string {
	var (
		ldflagsBuilder strings.Builder
		cmdArgs        []string
	)

	ldflagsBuilder.WriteString("-s -w")
	if c.HideConsole {
		ldflagsBuilder.WriteString(" -H windowsgui")
		log.Infof("HideConsole: %v", c.HideConsole)
	}
	// 编译模式, go help buildmode
	if c.BuildMode != "" {
		ldflagsBuilder.WriteString(fmt.Sprintf(" -buildmode=%s", c.BuildMode))
		log.Infof("BuildMode: %v", c.BuildMode)
	}

	cmdArgs = []string{cmd}
	if cmd == "garble" {
		if c.GDebug {
			cmdArgs = append(cmdArgs, "-debug")
		}
		if c.Tiny {
			cmdArgs = append(cmdArgs, "-tiny")
		}
		if c.Literals {
			cmdArgs = append(cmdArgs, "-literals")
		}
		if c.GSeed {
			cmdArgs = append(cmdArgs, "-seed=random")
		}
	}

	cmdArgs = append(cmdArgs,
		"build",
		"-ldflags="+ldflagsBuilder.String(),
		"-o", c.ExeFileName,
		"-trimpath", c.GoFileName,
	)

	return cmdArgs
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
		return fmt.Errorf("%s failed: %v", cmd, err)
	}

	return nil
}
