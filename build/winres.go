package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"variant/log"
)

func checkInit(path string) bool {
	resPath := filepath.Join(path, "winres")

	folderInfo, err := os.Stat(resPath)
	if os.IsNotExist(err) || !folderInfo.IsDir() {
		return true
	}

	return false
}

func (cOpts CompileOpts) execCmd(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cOpts.CompilePath

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s failed: %v", cmd, err)
	}

	return nil
}

func (cOpts CompileOpts) resInit() error {
	resCmd := []string{"go-winres", "init"}
	if err := cOpts.execCmd(resCmd); err != nil {
		return err
	}

	return nil
}

func (cOpts CompileOpts) makeInit() error {
	resCmd := []string{"go-winres", "make"}
	if err := cOpts.execCmd(resCmd); err != nil {
		return err
	}

	return nil
}

func (cOpts CompileOpts) Winres() error {
	if checkInit(cOpts.CompilePath) {
		err := cOpts.resInit()
		if err != nil {
			return fmt.Errorf("resInit() err: %w", err)
		}
	}

	err := cOpts.makeInit()
	if err != nil {
		return fmt.Errorf("makeInit() err: %w", err)
	}

	patchCmd := []string{
		"go-winres",
		"patch",
		"--in",
		".\\winres\\winres.json",
		cOpts.ExeFileName,
	}

	if err = cOpts.execCmd(patchCmd); err != nil {
		return fmt.Errorf("patch err: %w", err)
	}

	log.Infof("Patch Done: %v", cOpts.ExeFileName)

	return nil
}
