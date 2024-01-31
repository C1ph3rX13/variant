package remote

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"variant/log"
)

type Transfer struct {
	Src   string
	Path  string
	Proxy string
}

func (transfer Transfer) execCmd(args []string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)

	if transfer.Path != "" {
		cmd.Dir = transfer.Path
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s failed: %v", cmd, err)
	}

	return string(output), nil
}

func (transfer Transfer) CurlUpload() (string, error) {
	curlArgs := []string{
		"curl",
		"--upload-file",
		filepath.Join(transfer.Path, transfer.Src),
		fmt.Sprintf("https://transfer.sh/%s", transfer.Src),
	}

	// 代理设置
	if transfer.Proxy != "" {
		curlArgs = append(curlArgs, "-x", fmt.Sprintf("socks5://%s", transfer.Proxy))
	}

	log.Infof("Upload: %s", curlArgs)
	output, err := transfer.execCmd(curlArgs)
	if err != nil {
		return "", err
	} else {
		log.Infof("Upload Url: %s", output)
	}

	return output, nil
}