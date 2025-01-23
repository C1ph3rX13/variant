package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (upx UpxOpts) Pack() error {
	if err := upx.verify(); err != nil {
		return fmt.Errorf("precheck failed: %w", err)
	}

	args := []string{upx.Level, "-q", "-v"}
	args = append(args, upx.buildFlags()...)
	args = append(args, upx.exePath())

	cmd := exec.Command(filepath.Join(upx.UpxPath, "upx.exe"), args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}
	return nil
}

func (upx UpxOpts) verify() error {
	if _, err := os.Stat(filepath.Join(upx.UpxPath, "upx.exe")); err != nil {
		return fmt.Errorf("UPX validation: %w", err)
	}
	if upx.Level == "" {
		return fmt.Errorf("compression level required\n%s", upx.help())
	}
	return nil
}

func (upx UpxOpts) buildFlags() (flags []string) {
	if upx.Keep {
		flags = append(flags, "-k")
	}
	if upx.Force {
		flags = append(flags, "-f")
	}
	return
}

func (upx UpxOpts) exePath() string {
	return filepath.Join(upx.SrcPath, upx.SrcExe)
}

func (upx UpxOpts) help() string {
	cmd := exec.Command(filepath.Join(upx.UpxPath, "upx.exe"), "-h")
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
