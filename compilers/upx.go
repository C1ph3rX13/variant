package compilers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"variant/log"
)

func (u *UpxCtx) Pack() error {
	if err := u.verify(); err != nil {
		return fmt.Errorf("precheck failed: %w", err)
	}

	args := []string{u.Level, "-q", "-v"}
	args = append(args, u.buildFlags()...)
	args = append(args, u.exePath())
	
	log.Infof("Running UPX: %s %v\n", filepath.Join(u.UpxPath, "upx.exe"), args)

	cmd := exec.Command(filepath.Join(u.UpxPath, "upx.exe"), args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}
	return nil
}

func (u *UpxCtx) verify() error {
	upxPath := filepath.Join(u.UpxPath, "upx.exe")
	if _, err := os.Stat(upxPath); err != nil {
		return fmt.Errorf("UPX not found at %s: %w", upxPath, err)
	}

	exePath := u.exePath()
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return fmt.Errorf("target file not found: %s", exePath)
	}

	if u.Level == "" {
		return fmt.Errorf("compression level required\n%s", u.help())
	}
	return nil
}

func (u *UpxCtx) buildFlags() (flags []string) {
	if u.Keep {
		flags = append(flags, "-k")
	}
	if u.Force {
		flags = append(flags, "-f")
	}
	return
}

func (u *UpxCtx) exePath() string {
	path := filepath.Join(u.WorkingDir, u.SrcExe)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("File not found: %s", path))
	}
	return path
}
func (u *UpxCtx) help() string {
	cmd := exec.Command(filepath.Join(u.UpxPath, "upx.exe"), "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fallbackHelp() // 失败时返回硬编码帮助
	}
	return string(output)
}

func fallbackHelp() string {
	return `Compression tuning options:
			  -1     compress faster [-123456789]
			  -9     compress better [-123456789]
			  --lzma    try LZMA [slower but tighter]
			  --brute   try all methods [slow]
			  --ultra-brute  try more variants [very slow]`
}
