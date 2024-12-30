# variant

Golang Malware Framework

## Description

本项目会不断添加各种免杀的技术，但是不适合直接不做任何修改的编译和使用，即使是有随机特征的编译

Docs：[Install | variant](https://github.com/C1ph3rX13/variant/tree/main/docs/Install.md)，[Usage | variant](https://github.com/C1ph3rX13/variant/tree/main/docs/Usage.md)

## Update

## 2024.12.30

1. 说明文档更新：安装和代码使用

### 2024.12.19

1. `hook`，新增`AMSIByPass()`，进程选择`powershell.exe`即可
2. `crypto`，新增中文加密`Buddha()`
3. `sandbox`，新增时区判断`IsBeijingTimezone()`
4. 新增`inject`注入模块，`AddressOfEntryPointInject()`、`CreatRemoteThreadInject()`
5. TODO：删除`wdll`模块，使用衍生库 [C1ph3rX13 | xwindows](https://github.com/C1ph3rX13/xwindows)

### 2024.10.18

1. 更新`SigTheif`，简化代码逻辑，更换读取文件的方法
2. 更新渲染模板 v6，更换调用逻辑
3. `demo`更新适配 v6 渲染模板

### 2024.9.26

1. 新增`gores`模块，支持`自定义`文件的资源信息和`复制`其他对象资源信息（todo：ICON随机Hash）

```go
// 自定义资源信息，修改 variant/gores/render.go 中的 func NewResDate()
// 推荐使用下面 Extract() 方法直接复制指定对象的资源信息和文件
func NewResDate() ResDate {
	goRes := ResDate{
		ICOName:                           "icon.png",
		Name:                              "WPS Office",
		Version:                           "12.1.0.16399",
		Description:                       "WPS Office",
		MinimumOs:                         "win7",
		ExecutionLevel:                    "requireAdministrator",
		UIAccess:                          false,
		AutoElevate:                       true,
		DpiAwareness:                      "system",
		DisableTheming:                    false,
		DisableWindowFiltering:            false,
		HighResolutionScrollingAware:      false,
		UltraHighResolutionScrollingAware: false,
		LongPathAware:                     false,
		PrinterDriverIsolation:            false,
		GDIScaling:                        false,
		SegmentHeap:                       false,
		UseCommonControlsV6:               false,
		FixedFileVersion:                  "12.1.0.16399",
		FixedProductVersion:               "WPS Office",
		Comments:                          "",
		CompanyName:                       "",
		FileDescription:                   "WPS Office",
		FileVersion:                       "12.1.0.16399",
		InternalName:                      "",
		LegalCopyright:                    "Copyright©2024 Kingsoft Corporation. All rights reserved.",
		LegalTrademarks:                   "",
		OriginalFilename:                  "wps_host.exe",
		PrivateBuild:                      "",
		ProductName:                       "WPS Office",
		ProductVersion:                    "WPS Office",
		SpecialBuild:                      "",
	}

	return goRes
}

// 渲染输出 winres.json 文件
// 将 ICON 和 winres.json 放置在编译目录的 winres 文件夹中即可
winres := gores.ResTmpl{
		ResPath:   "gores/gores.tmpl",
		OutputDir: "output",
	}

	err := winres.ResRender()
	if err != nil {
		panic(err)
	}
```

2. 更新`gores`编译方法，从`build`模块中分离，新增资源提取方法`Extract()`

```go
    // 添加图标和文件信息
	winres := gores.GoWinRes{
		CompilePath: "output",          // 指定编译目录
		ExtractFile: "Code.exe",        // 指定提取资源文件的对象
		ExtractDir:  "",                // 指定提取资源文件后输出的路径
		PatchFile:   cOpts.ExeFileName, // 指定使用 Patch 添加资源文件的对象
	}

	// 提取 vscode 所有的资源文件
	err = winres.Extract()
	if err != nil {
		log.Fatal(err)
	}

	// 使用 Patch 添加资源文件到编译后的程序
	err = winres.HandleWinRes()
	if err != nil {
		log.Fatal(err)
	}
```

3. 资源提取方法`Extract()`实现的效果

![winres](https://raw.githubusercontent.com/C1ph3rX13/variant/main/images/winres.png)

### 2024.8.15

1. 更新`render`模块，支持新增的`cloader`，模板渲染调用结构体优化
2. 更新渲染模板`v5`，支持新增的`cloader`
3. `demo`更新，适配其他模块的更新
4. `Hide Cmd`隐藏执行窗口的函数移动到`sandbox`模块
5. `DLL`渲染模板和调用方式更新

### 2024.8.14

1. `cloader`模块新增 23 个`CGO`类型 loader
2. 跟随编译需求更新`initialize.bat`，新增依赖检查
3. `xwindows`模块文档完善，新增 20+ API （详情见：[xwindows](https://github.com/C1ph3rX13/xwindows) 仓库）

### 2024.8.13

1. 新增`cloader`模块，使用`CGO`调用`C`，代码更简洁

### 2024.7.20

1. `enc`模块改名`encoder`
2. `encoder`模块新增降熵方法`ReduceEntropy()`、恢复方法`ReverseEntropy()`
3. `render`模块`TmplRender()`方法更新，符合整体设计逻辑

```go
func (tOpts TmplOpts) TmplRender() error {}
```

4. `wdll`模块更新，添加更多API

### 2024.6.26

1. `hook`模块重构，现支持`etwpatch`，`selfdelete`
2. `dynamic`模块重构，现支持`获取任意URL资源 SHA256 的指定切片`，`新增对AES和DES Key和IV生成位数限制`
3. `network`模块优化，所有的请求方式均支持代理（可选），优化客户端的仿真设置

### 2024.5.13

1. 新增`DLL`编译：指定`DllBase.tmpl`模板，使用`BuildMode: "c-shared"`
2. 新增`DLL`渲染模板：`DllBase.tmpl`
3. 新增`DLL`编译demo，查看`demo`下`GoDLL`文件夹

### 2024.4.1

1. 新增计划任务隐藏

   > + 修改 Index 为 0 隐藏：IndexToZero()
   > + 修改/删除 SD 项：ChangeSD() & DeleteSD()
   > + 删除注册表文件夹 SD 项：DeleteSD()
   > + 删除注册表中的计划任务文件&文件夹：RegDeleteTaskDir()
   > + 删除计划任务XML文件：DeleteTaskDir() & DeleteTaskFile()

![RegTasks](https://raw.githubusercontent.com/C1ph3rX13/variant/main/images/RegTask.gif)

### 2024.3.25

1. 优化`alive`自动维权模块：`WinTaskXML(), SetWinTask() `

### 2024.3.21

1. 新增`alive`自动维权模块，现支持注册表启动项，计划任务`COM API`
2. 文档更新

### 2024.3.19

1. 更新自删除，将`hook.SelfDelete()`置于`loader,inject`方法之前即可

![selfdelete](https://raw.githubusercontent.com/C1ph3rX13/variant/main/images/SelfDelete.gif)

2. 新增`hashdump`功能
2. 更新`demo`

### 2024.3.12

1. 新增`initialize.bat`，自动配置项目依赖项，移除`go.mod, go.sum`
2. 新增反沙箱：`Beep`，利用该方法达到`Sleep`的效果
3. `wdll`模块同步更新

### 2024.2.28

1. `remote`模块改名为`network`
2. 新增多种加密：`elliptic_curve`，`morse`，`pokemon`，`rot13`，`rot47`
3. 模板更新`V4 Base.tmpl`：兼容全部渲染方式（参数加载暂时除外）
4. 加载模块更新：`ADsMemLoad`
5. `remder`模块同步更新
6. `enc`模块新增：`PokemonStrings`，`pokemon`加密专用方法；修改其他加密的兼容性

### 2024.2.20

1. 新增`loader`：`ipv4`, `macaddress`, `enumsystemlocales + Hell's Gate + Halo's Gate technique`
2. `Dll`模块重构调用方式，特征更新
3. 新增`hook`模块，`Hook函数检测`，`ETW Patch`，`权限检测/提权`

### 2024.2.18

1. 新增`loader`：`EnumerateLoadedModulesLoad`, `EnumChildWindowsLoad`, `EnumPageFilesWLoad`
2. `Dll`模块根据`loader`同步更新
3. 新增`garble`编译参数：`-tiny`
4. 新增构建模式：`-buildmode`，详情执行`go help buildmode`查看
5. 新增`FileAnalyzer`方法，计算文件的特征：`entropy`, `md5`, `sha1`, `sha256`, `sha512`

### 2024.2.7

1. 新增`garble(需安装)`编译，支持`-seed, -literals, -debug`，但会导致编译的文件体积增大和熵值增加
2. 重构编译模块，现在支持`原生go编译, garble编译`
3. 调整编译的流程和日志输出逻辑
4. 重构`Upx`模块，支持自定义`Upx.exe`的路径
5. 重构`Winres`模块，现在使用`HandleWinRes`方法可以直接添加图标
6. 预计新增Go编译器： [llvm](https://github.com/llvm/llvm-project/tree/llvmorg-17.0.6), [tinygo](https://github.com/tinygo-org/tinygo) - 实验性
7. 新增`install.bat`，初始化工具

### 2024.2.6

1. 新增两种远程加载的方式`UsersCloud, file.io (web)`
2. 新增`github.com/imroc/req/v3`的请求客户端
3. 优化远程加载模块的函数描述
4. 简化远程模块上传加密文件的结构体
5. 模板更新，兼容新的远程加载方式

### 2024.2.5

1. 新增参数加载模块，可以自定义(随机)密钥
2. 模板简化，兼容所有模块的渲染
3. 渲染模块优化，兼容所有模块的调用
4. 新增远程加载模块，远程加密数据的上传会在渲染阶段完成，上传(`curl:已完成, web:开发中...`)支持代理
5. 新增动态数据模块可以和任意加载方式联动
6. 新增上传加密`Payload`随机化
7. 新增`loader:earlybird`

### 2024.2.4

1. 新增动态获取解密数据的模块 - `Dynamic: payload, key, iv`
2. 按照本地、远程、参数加载的方式重构，减少模板渲染的复杂度
3. 新增分离加载模块
4. 渲染模板优化

### 2024.1.30

1. `lzw`压缩导致熵值上升到`7.0+`；测试`fmt.Printf("Hello World")`编译后熵值在`6.0-6.1`之间，压缩模块下次一定
2. 远程加载：模块、模板、配置更新
3. 重构加密模块，利用反射根据传入的方法签名判断加密
4. 使用`windows package`重写了`DLL`调用模块，[`syscall`&`windows`的总结](https://c1ph3rx13.github.io/posts/2024-01-30/The-low-level-Operating-System-Primitives/)
5. 模板新增根据自定义导入对应库设置
6. 新增多个多重加密/编码方法，`Base32/62`编码测试熵值比较低

### 2024.1.29

1. 新增渲染判断，模板可根据结构体来渲染
2. 更改为验证为沙箱之后程序正常退出
3. 新增代码检查：检查通过再编译；初始化：`go install golang.org/x/tools/cmd/goimports@latest`
4. 细化结构体：根据功能区分
5. 新增压缩算法：`lzw`, `zstd`；熵值模块
6. 新增动态方法：`GetSelfSHA256Nth`
7. 优化代码逻辑

### 2024.1.27

1. 模块化配置
2. 模板渲染
3. 编译控制
4. UPX压缩：需要检查`upx.exe`是否正确放置在`build`文件夹，才能正确初始化
5. 签名伪造
6. 添加图标文件信息：初始化：`go install github.com/tc-hib/go-winres@latest`
7. 模板更新
8. 反沙箱模块

## Thanks

https://github.com/Ne0nd0g/go-shellcode

https://github.com/safe6Sec/GolangBypassAV

https://github.com/afwu/GoBypass

https://github.com/piiperxyz/AniYa

https://github.com/wumansgy/goEncrypt

https://github.com/TideSec/GoBypassAV

https://github.com/Pizz33/GobypassAV-shellcode

https://github.com/timwhitez/Doge-Gabh

https://github.com/MrTuxx/OffensiveGolang
