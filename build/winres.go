package build

import (
	"fmt"
	"os"
	"path/filepath"
	"variant/log"
)

func (c CompileOpts) HandleWinRes() error {
	if c.shouldInitRes() {
		if err := c.initRes(); err != nil {
			return fmt.Errorf("initRes() err: %w", err)
		}
	}

	if err := c.makeRes(); err != nil {
		return fmt.Errorf("makeRes() err: %w", err)
	}

	if err := c.patchRes(); err != nil {
		return fmt.Errorf("patchRes() err: %w", err)
	}

	log.Infof("Patch Succeeded: %v", c.ExeFileName)
	return nil
}

func (c CompileOpts) shouldInitRes() bool {
	// 检查编译目录是否存在 winres 文件夹
	resPath := filepath.Join(c.CompilePath, "winres")
	folderInfo, err := os.Stat(resPath)
	return os.IsNotExist(err) || !folderInfo.IsDir()
}

func (c CompileOpts) initRes() error {
	return c.execCmd([]string{"go-winres", "init"})
}

func (c CompileOpts) makeRes() error {
	return c.execCmd([]string{"go-winres", "make"})
}

func (c CompileOpts) patchRes() error {
	patchCmd := []string{
		"go-winres",
		"patch",
		"--in",
		filepath.Join("winres", "winres.json"),
		c.ExeFileName,
	}
	log.Infof("Patching: %v", patchCmd)
	return c.execCmd(patchCmd)
}
