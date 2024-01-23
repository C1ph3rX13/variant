package build

import (
	"fmt"
	"os"
	"os/exec"
	"variant/log"
)

func (cOpts CompileOpts) Compile() error {
	ldflags := "-s -w"
	if cOpts.HideConsole {
		ldflags = "-H windowsgui " + ldflags
	}

	cmdArgs := []string{
		"go", "build",
		fmt.Sprintf("-ldflags=%s", ldflags),
		"-o", cOpts.ExeFileName,
		"-trimpath", cOpts.GoFileName,
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Infof("CompilePath: %v", cOpts.CompilePath)
	log.Infof("HideConsole: %v", cOpts.HideConsole)
	log.Infof("Go Build: %v", cmdArgs)

	if cOpts.CompilePath != "" { // 如果指定了编译目录，则设置 Cmd.Dir 属性
		cmd.Dir = cOpts.CompilePath
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("<cmd.Run()> err: %w", err)
	} else {
		log.Infof("Compile Done: %v", cOpts.ExeFileName)
	}

	return nil
}
