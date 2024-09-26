package gores

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"variant/log"
)

func (gwr GoWinRes) HandleWinRes() error {
	if gwr.shouldInitRes() {
		if err := gwr.initRes(); err != nil {
			return fmt.Errorf("initRes() err: %w", err)
		}
	}

	if err := gwr.makeRes(); err != nil {
		return fmt.Errorf("makeRes() err: %w", err)
	}

	if err := gwr.patchRes(); err != nil {
		return fmt.Errorf("patchRes() err: %w", err)
	}

	log.Infof("Patch Succeeded: %v", gwr.PatchFile)
	return nil
}

func (gwr GoWinRes) shouldInitRes() bool {
	// 检查编译目录是否存在 winres 文件夹
	resPath := filepath.Join(gwr.CompilePath, "winres")
	folderInfo, err := os.Stat(resPath)
	return os.IsNotExist(err) || !folderInfo.IsDir()
}

func (gwr GoWinRes) initRes() error {
	return gwr.execResCmd([]string{"go-winres", "init"})
}

func (gwr GoWinRes) makeRes() error {
	return gwr.execResCmd([]string{"go-winres", "make"})
}

func (gwr GoWinRes) patchRes() error {
	patchCmd := []string{
		"go-winres",
		"patch",
		"--in",
		filepath.Join("winres", "winres.json"),
		gwr.PatchFile,
	}
	log.Infof("Patching: %v", patchCmd)
	return gwr.execResCmd(patchCmd)
}

func (gwr GoWinRes) Extract() error {
	patchCmd := []string{
		"go-winres",
		"extract",
		gwr.ExtractFile,
	}

	if gwr.ExtractDir != "" {
		patchCmd = append(patchCmd, "--dir", gwr.ExtractDir)
		log.Infof("Output Directory (default: winres): %v", gwr.ExtractDir)
	}

	log.Infof("Extract: %v", patchCmd)
	return gwr.execResCmd(patchCmd)
}

func (gwr GoWinRes) execResCmd(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if gwr.CompilePath != "" {
		cmd.Dir = gwr.CompilePath
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s failed: %v", cmd, err)
	}

	return nil
}
