package sandbox

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"variant/crypto"
)

func WMICCheckVirtual() (bool, error) {
	// 执行命令获取计算机系统模型信息
	cmd := exec.Command("wmic", "path", "Win32_ComputerSystem", "get", "Model")
	stdout, err := cmd.Output()
	if err != nil {
		return false, err
	}

	// 将模型信息转换为小写字母方便匹配
	model := strings.ToLower(string(stdout))

	// 检查模型信息中是否包含虚拟机的关键词
	virtualKeywords := []string{"virtualbox", "virtual", "vmware", "kvm", "bochs", "hvm domu", "parallels"}
	for _, keyword := range virtualKeywords {
		if strings.Contains(model, keyword) {
			return true, nil
		}
	}

	return false, nil
}

// CheckVirtualFiles 调用 PathExists 检查沙箱或者虚拟机关键文件是否存在，如果存在则退出当前进程
func CheckVirtualFiles() {
	//files := []string{
	//	"C:\\windows\\System32\\Drivers\\Vmmouse.sys",
	//	"C:\\windows\\System32\\Drivers\\vmtray.dll",
	//	"C:\\windows\\System32\\Drivers\\VMToolsHook.dll",
	//	"C:\\windows\\System32\\Drivers\\vmmousever.dll",
	//	"C:\\windows\\System32\\Drivers\\vmhgfs.dll",
	//	"C:\\windows\\System32\\Drivers\\vmGuestLib.dll",
	//	"C:\\windows\\System32\\Drivers\\VBoxMouse.sys",
	//	"C:\\windows\\System32\\Drivers\\VBoxGuest.sys",
	//	"C:\\windows\\System32\\Drivers\\VBoxSF.sys",
	//	"C:\\windows\\System32\\Drivers\\VBoxVideo.sys",
	//	"C:\\windows\\System32\\vboxdisp.dll",
	//	"C:\\windows\\System32\\vboxhook.dll",
	//	"C:\\windows\\System32\\vboxoglerrorspu.dll",
	//	"C:\\windows\\System32\\vboxoglpassthroughspu.dll",
	//	"C:\\windows\\System32\\vboxservice.exe",
	//	"C:\\windows\\System32\\vboxtray.exe",
	//	"C:\\windows\\System32\\VBoxControl.exe",
	//}

	files, _ := crypto.AesBase32Decrypt("ZJHLTRY5AHDK7LYBXRYUSXNRXZ3PBXPDQJILPV3VHAVSD4EBS6DNZ2NFVW3KX7PBUJIA2WANNP2ZA4GVFTIGN7MZUGTVPYABSA2GXJX6AVFUTB2RCG6EKUTHA3MCCVFKUEHX7BMBJBMXUXX3KRFFMPDV3MXX4OTPCIFSXINTSL6NTPFXJCOE65KKNSCLYMMV4LCBQGECUDT2IW3PB576B62LMAJSN6YKKYNQBH5JLO7KSJ4H6KOXYMZBVNNMCJ7FTPZSXWZ2NUUZUC3LSOIGZZVIDTORVC2SFE64AWYR2OLAQFCWZND6YG4Q2DQUEHYKQRAHFEDX72EM4Z6GEEESUXDNC5ZWWBVA4BMECB7YIFX2HZSFZAMFGVG775DNS2ORC34FK2E7WYWBYZN2XACVK32FU3OGTKJZIDHRSEUBHNXN32JKF2OGZG5VAX26VNKPRPMVRYHHUV5ODV6XGJH25K7PEADDQB6FMMNF632EW2CHVFVZMX4JPUEGIJQRQ2E7ISB4UPCPUPYWTTCG6B52ALOUFLFFMSO66BWU3JOSBAZASTM4JOPSDSAJTKTPFD4QBY37DT3NMKPP6YZBREH4QPDBJYVFPWRP5HBRFXZI6DTCGP4KVZQHLBT4LQPSKC4XB3BWAQ5BYFHFNVO3TYUYYMCEWEFBEGVNF4APOZU2ZYMGNOQ4MAICY5Z7WHTJCAYNONU2PGQ2JBH7UXBGSKJIQWMTOIE7OCROVB46X362VYC5FFQIGORCAORJ46GNT43XYCGIQZJCAYWNBCLY4UT43XIVTFA47Z3IC5YNRGANVQEKMJS6KUO5SM57TR2NBAZK3QBEHUKTWIITWLMOBGKZ3CUYU33PAOEXKDYDZARPIMHYDATT6RPBIFJPLJI7KPWIXJITKVUP7BZEJLYY2S7LCB5YF7CCC5D5UW6ZQG45N4LXUB4EP5ZYPHVJPE63Q43SZA4NJPKMPD655B6FLAKQUQAKZBMJQHTPX76QSOEJYYP32EAI3EHHGXOW7JBBEUFFIABQDQ6XG7L4742CB5QDQ7JKCY", []byte("0123456789aaaaaa"), []byte("0123456789aaaaaa"))

	for _, file := range files {
		exists, _ := PathExists(strconv.Itoa(int(file)))
		if exists {
			os.Exit(0)
		}
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
