# variant

Golang Malware Framework

## Description

Project Configuration

## Msys2

[Download | Msys2](https://www.msys2.org/)

### Update Msys2 Repo

[Msys2 | Tsinghua Open Source Mirror](https://mirrors.tuna.tsinghua.edu.cn/help/msys2/)

```cmd
# 更新本地软件包数据库，但不安装或更新任何软件包
pacman -Sy

# 更新软件包数据库，并且更新系统中所有已安装的包，确保系统保持最新状态
pacman -Syu 

# 只更新已安装的包，而不更新软件包数据库，通常在数据库已经同步的情况下使用
pacman -Su 
```

### Compilation Environment

```cmd
pacman -S mingw-w64-x86_64-gcc

# pass
pacman -S  mingw-w64-x86_64-toolchain

# pass
pacman -S mingw-w64-x86_64-gtk3

# pass
pacman -S  mingw-w64-x86_64-glade
```

### System Environment Variables

MINGW_HOME

```
D:\msys64\mingw64
```

C_INCLUDE_PATH

```
%MINGW_HOME%\include
```

LIBRARY_PATH

```
%MINGW_HOME%\lib
```

Path

```
%MINGW_HOME%\bin;
```

## xwindows

下载并置于项目根目录

xwindows | https://github.com/C1ph3rX13/xwindows

## Initialize

[Run initialize.bat](https://github.com/C1ph3rX13/variant/blob/main/initialize.bat) to initialize
