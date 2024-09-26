package gores

type GoWinRes struct {
	CompilePath string // 编译目录
	ExtractFile string // 需要提取资源文件的对象
	ExtractDir  string // 提取对象资源存放路径
	PatchFile   string // 需要添加资源文件的对象
}

type ResTmpl struct {
	ResPath   string
	OutputDir string
}

type ResDate struct {
	ICOName                           string // 图标资源的名称
	Name                              string // 应用程序的名称
	Version                           string // 应用程序的版本号
	Description                       string // 应用程序的描述信息
	MinimumOs                         string // 指定应用程序所需的最低操作系统版本
	ExecutionLevel                    string // 指定应用程序的执行权限级别 ("asInvoker;requireAdministrator;highestAvailable")
	UIAccess                          bool   // 是否允许用户界面访问
	AutoElevate                       bool   // 是否自动提升权限
	DpiAwareness                      string // DPI（像素每英寸），设置为"system"表示依赖于系统的DPI
	DisableTheming                    bool   // 是否禁用主题支持
	DisableWindowFiltering            bool   // 是否禁用窗口过滤
	HighResolutionScrollingAware      bool   // 是否支持高分辨率滚动
	UltraHighResolutionScrollingAware bool   // 是否支持超高分辨率滚动
	LongPathAware                     bool   // 是否支持长路径
	PrinterDriverIsolation            bool   // 打印机驱动程序隔离
	GDIScaling                        bool   // GDI缩放支持
	SegmentHeap                       bool   // 段堆支持
	UseCommonControlsV6               bool   // 是否使用公共控件v6
	FixedFileVersion                  string // 文件的版本号
	FixedProductVersion               string // 产品的版本号
	Comments                          string // 备注信息
	CompanyName                       string // 公司名称
	FileDescription                   string // 文件描述
	FileVersion                       string // 文件的版本号
	InternalName                      string // 内部名称
	LegalCopyright                    string // 法律版权信息
	LegalTrademarks                   string // 法律商标信息
	OriginalFilename                  string // 原始文件名
	PrivateBuild                      string // 私有构建信息
	ProductName                       string // 产品名称
	ProductVersion                    string // 产品版本
	SpecialBuild                      string // 特殊构建信息
}
