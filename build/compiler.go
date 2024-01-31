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

	log.Infof("CompilePath: %v", cOpts.CompilePath)
	log.Infof("HideConsole: %v", cOpts.HideConsole)
	log.Infof("Go Build: %v", cmdArgs)

	if err := cOpts.fmtCode(); err != nil {
		return err
	}

	if err := cOpts.execCmd(cmdArgs); err != nil {
		return err
	}
	log.Infof("Compile Done: %v", cOpts.ExeFileName)

	return nil
}

func (cOpts CompileOpts) fmtCode() error {
	fmtCmd := []string{"goimports", "-w", cOpts.GoFileName}
	log.Infof("Gofmt: %v", fmtCmd)

	if err := cOpts.execCmd(fmtCmd); err != nil {
		return err
	}

	return nil
}

func (cOpts CompileOpts) execCmd(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if cOpts.CompilePath != "" {
		cmd.Dir = cOpts.CompilePath
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s failed: %v", cmd, err)
	}

	return nil
}
