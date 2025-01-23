package build

import (
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

// ReadCertTable 从PE文件中读取证书表
func (ct *CertThief) ReadCertTable(filePath string) ([]byte, error) {
	f, err := os.Open(filePath) // 修改为流式读取
	if err != nil {
		return nil, fmt.Errorf("file open failed: %w", err)
	}
	defer f.Close()

	peFile, err := pe.NewFile(f)
	if err != nil {
		return nil, fmt.Errorf("PE parse failed: %w", err)
	}
	return peFile.CertificateTable, nil
}

// SaveSignedPE 保存带有证书表的PE文件
func (ct *CertThief) SaveSignedPE(certTable []byte, outputFileName string) error {
	p := filepath.Join(ct.SignDir, ct.SrcFile)
	peFile, errOpen := pe.Open(p)
	if errOpen != nil {
		return fmt.Errorf("failed to open PE file: %v", errOpen)
	}
	defer peFile.Close()

	peFile.CertificateTable = certTable

	outputFilePath := filepath.Join(ct.SignDir, outputFileName)
	if errWrite := peFile.WriteFile(outputFilePath); errWrite != nil {
		return fmt.Errorf("failed to write PE file: %v", errWrite)
	}

	log.Infof("Signed file saved to: %s", outputFilePath)
	return nil
}

// SaveCertificate 将证书表保存到指定文件
func (ct *CertThief) SaveCertificate() error {
	p := filepath.Join(ct.SignDir, ct.DstFile)
	certTable, err := ct.ReadCertTable(p)
	if err != nil {
		return err
	}

	certFilePath := filepath.Join(ct.SignDir, ct.CertFile)
	if errWrite := os.WriteFile(certFilePath, certTable, os.ModePerm); errWrite != nil {
		return fmt.Errorf("failed to save certificate: %v", errWrite)
	}

	log.Infof("Certificate saved to %s", certFilePath)
	return nil
}

// SignWithStolenCert 使用窃取的证书对源文件进行签名
func (ct *CertThief) SignWithStolenCert() error {
	pCert := filepath.Join(ct.SignDir, ct.CertFile) // 证书路径
	bCert, err := os.ReadFile(pCert)
	if err != nil {
		return fmt.Errorf("sign failed: %w (cert: %s)", err, pCert)
	}

	return ct.SaveSignedPE(bCert, ct.SignedPE)
}

// SignExecutable 使用目标文件的证书对源文件进行签名
func (ct *CertThief) SignExecutable() error {
	pFile := filepath.Join(ct.SignDir, ct.DstFile) // 待签名文件路径
	certTable, err := ct.ReadCertTable(pFile)
	if err != nil {
		return err
	}
	return ct.SaveSignedPE(certTable, ct.SignedPE)
}
