package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"variant/log"
)

func (upx UpxOpts) UpxPacker() error {
	if err := upx.checkUpx(); err != nil {
		log.Fatal(err)
	}

	if upx.Level == "" {
		return fmt.Errorf(upx.levelHelp())
	}

	args := []string{
		upx.Level,
		"-q",
		"-v",
	}

	if upx.Keep {
		args = append(args, "-k")
	}

	if upx.Force {
		args = append(args, "-f")
	}

	if upx.SrcExe != "" {
		args = append(args, filepath.Join(upx.SrcPath, upx.SrcExe))
	}

	log.Infof("Upx Compress: %v", args)
	cmd := exec.Command(filepath.Join(upx.UpxPath, "upx.exe"), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("<cmd.Run()> err: %w", err)
	}

	return nil
}

func (upx UpxOpts) checkUpx() error {
	_, err := os.Stat(filepath.Join(upx.UpxPath, "upx.exe"))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("can't find upx.exe in %s", upx.UpxPath)
		}
	}

	return nil
}

func (upx UpxOpts) levelHelp() string {
	help := `
	Compression tuning options:
	-1     compress faster [-123456789]                  
	-9     compress better [-123456789]
	--lzma              try LZMA [slower but tighter than NRV]
	--brute             try all available compression methods & filters [slow]
	--ultra-brute       try even more compression variants [very slow]`

	return help
}
