package build

import (
	"fmt"
	"github.com/Binject/debug/pe"
	"os"
	"path/filepath"
	"variant/log"
)

// SaveCert 保存证书到指定文件
func (sOpts SignOpts) SaveCert() error {
	// 指定保存EXE文件证书
	target := filepath.Join(sOpts.SignPath, sOpts.Thief)

	cert, err := getCert(target)
	if err != nil {
		return fmt.Errorf("<SaveCert getCert() err: %v>", err)
	}

	err = os.WriteFile(filepath.Join(sOpts.SignPath, sOpts.DstCert), cert, os.ModePerm)
	if err != nil {
		return fmt.Errorf("<SaveCert os.WriteFile() err: %v>", err)
	} else {
		log.Infof("Save %s Cert Done: %s", sOpts.Thief, sOpts.DstCert)
	}

	return nil
}

// getCert 获取签名证书
func getCert(singed string) ([]byte, error) {
	peFile, err := pe.Open(singed)
	if err != nil {
		return nil, fmt.Errorf("<getCert pe.Open() err: %v>", err)
	}
	defer peFile.Close()

	if len(peFile.CertificateTable) == 0 {
		return nil, fmt.Errorf("<CertificateTable is empty>")
	}

	return peFile.CertificateTable, nil
}

// signThief 将文件签名为指定证书
func signThief(unSign string, signedOutput string, certTable []byte) error {
	peFile, err := pe.Open(unSign)
	if err != nil {
		return fmt.Errorf("<SignThief pe.Open() err: %v>", err)
	}
	defer peFile.Close()

	peFile.CertificateTable = certTable

	err = peFile.WriteFile(signedOutput)
	if err != nil {
		return fmt.Errorf("<SignThief pe.WriteFile() err: %v>", err)
	}

	return nil
}

// CertThief 使用指定证书签名文件
func (sOpts SignOpts) CertThief() error {
	// 指定签名所需的证书文件
	targetCert := filepath.Join(sOpts.SignPath, sOpts.Cert)

	thiefTable, err := os.ReadFile(targetCert)
	if err != nil {
		return fmt.Errorf("<CertThief os.ReadFile() err: %v>", err)
	}

	// 编译后未签名
	unSignExe := filepath.Join(sOpts.SignPath, sOpts.UnSign)
	// 添加签名后输出的文件
	signedOutput := filepath.Join(sOpts.SignPath, sOpts.Signed)

	err = signThief(unSignExe, signedOutput, thiefTable)
	if err != nil {
		return fmt.Errorf("<CertThief signThief() err: %v>", err)
	} else {
		log.Infof("CertThief Done: %s", sOpts.Signed)
	}

	return nil
}

// ExeThief 使用指定可执行文件签名另一个文件
func (sOpts SignOpts) ExeThief() error {
	// 指定签名所需的 EXE 文件
	targetExe := filepath.Join(sOpts.SignPath, sOpts.Thief)

	peFile, err := pe.Open(targetExe)
	if err != nil {
		return fmt.Errorf("<ExeThief pe.Open() err: %v>", err)
	}
	defer peFile.Close()

	// 编译后未签名
	unSignExe := filepath.Join(sOpts.SignPath, sOpts.UnSign)
	// 添加签名后输出的文件
	signedOutput := filepath.Join(sOpts.SignPath, sOpts.Signed)

	err = signThief(unSignExe, signedOutput, peFile.CertificateTable)
	if err != nil {
		return fmt.Errorf("<ExeThief signThief() err: %v>", err)
	} else {
		log.Infof("ExeThief Done: %s", sOpts.Signed)
	}

	return nil
}
