## frp Android AAR 构建说明

本文件说明如何在本仓库中生成 `frp.aar`（Android AAR），以及相关工具和环境变量的配置方法。

---

## 1. 前置条件

- 支持的平台：
  - Linux / macOS（推荐直接在 shell 中构建）
  - Windows：使用 **Git Bash / WSL** 来运行 `make` 和 `.sh` 脚本
- 已安装 Go（推荐 ≥ 1.20，当前为 Go 1.25.6）
- 本仓库当前模块名为 `github.com/fatedier/frp`，本文描述的是本分支现有结构下的构建方式，后续升级 frp 版本时只要目录结构与目标包名不变，本说明同样适用。
 - 已安装 `zip` / `unzip` 命令行工具（用于解压/重新打包 AAR）

确认 `GOPATH`：

```bash
go env GOPATH
```

在你当前环境下一般为（Windows 示例）：`C:\Users\admin\go`，并确保：

- `GOPATH/bin` 已加入系统 `PATH`（例如 `C:\Users\admin\go\bin`）

### 1.1 安装 zip / unzip

`make_frp_dex_aar.sh` 需要系统中有 `zip` / `unzip` 命令：

- **Debian / Ubuntu：**

  ```bash
  sudo apt-get install zip unzip
  ```

- **CentOS / RHEL：**

  ```bash
  sudo yum install zip unzip
  ```

- **macOS（Homebrew）：**

  ```bash
  brew install zip unzip
  ```

- **Windows + Git Bash（已安装 Chocolatey）：**

  在 **管理员 PowerShell 或 CMD** 中执行：

  ```powershell
  choco install zip unzip -y
  ```

  然后重新打开 Git Bash。

---

## 2. 安装 gomobile / gobind

在 **Linux / macOS shell** 或 **Windows 的 Git Bash / WSL** 中执行：

```bash
# 建议先设置 Go 国内代理（如已设置可跳过）
go env -w GOPROXY=https://goproxy.cn,direct

go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest
```

安装完成后，检查：

```bash
ls "$GOPATH/bin"
# 确认有 gomobile(.exe)、gobind(.exe)
```

### 2.1 初始化 gomobile（只需一次）

```bash
gomobile init
```

---

## 3. 安装 Android SDK 与 build-tools

Android SDK 可以通过 **Android Studio** 或独立的 **sdkmanager** 安装。下面给出典型步骤和路径示例。

### 3.1 Windows + Git Bash 环境

1. 安装 Android Studio（或者单独安装 Android SDK）。
2. 打开 Android Studio → `More Actions` → `SDK Manager`：
   - 在 **SDK Location** 中记下 `Android SDK location`，例如：
     - `C:\Users\admin\AppData\Local\Android\Sdk`
   - 在 **SDK Tools** 里勾选并安装：
     - `Android SDK Build-Tools`（建议选择一个稳定版本，例如 `34.0.0`）
3. 打开 Git Bash 或 CMD，确认 build-tools 目录下有该版本：

   ```bash
   ls "C:/Users/admin/AppData/Local/Android/Sdk/build-tools"
   # 看到 34.0.0、33.0.2 等版本目录
   ```

4. 将你打算使用的版本写入 `make_frp_dex_aar.sh` 中：

   ```bash
   BUILD_TOOLS_VERSION="34.0.0"
   ```

> 提示：如果你的 SDK 安装在其它盘（比如 `D:\Android\Sdk`），只需要把上面的路径和脚本里的默认 `SDK_ROOT` 改成对应路径即可。

### 3.2 Linux / macOS 环境

1. 安装 Android Studio 或使用命令行工具安装 SDK，例如：

   ```bash
   # 示例：使用 sdkmanager 安装（路径仅为示例）
   sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0"
   ```

2. 记下 SDK 安装路径，例如：

   ```text
   /home/yourname/Android/Sdk
   ```

3. 设置环境变量（见下节 `ANDROID_SDK_ROOT`）：

   ```bash
   export ANDROID_SDK_ROOT="$HOME/Android/Sdk"
   ```

4. 同样在 `make_frp_dex_aar.sh` 中将 `BUILD_TOOLS_VERSION` 设置为你安装的 build-tools 版本。

---

## 4. 配置 ANDROID_SDK_ROOT / ANDROID_HOME

本项目中使用的 `make_frp_dex_aar.sh` 依赖 **Android SDK 路径**。

有两种配置方式：**系统级** 或 **每次临时设置（推荐直接在构建 shell 中设置）**。

### 4.1 系统级（Windows 示例）

1. 打开：此电脑 → 右键 → 属性 → 高级系统设置 → 环境变量。
2. 在“系统变量”中新建：
   - **变量名**：`ANDROID_SDK_ROOT`
   - **变量值**：`C:\Users\admin\AppData\Local\Android\Sdk`
3. 可选：同样方式设置 `ANDROID_HOME` 为相同路径。
4. 确认后，**重新打开 Git Bash / 终端**。

### 4.2 仅当前 shell 会话（Linux / macOS / Git Bash）

```bash
export ANDROID_SDK_ROOT="/c/Users/admin/AppData/Local/Android/Sdk"  # Windows Git Bash 示例
# 或 Linux / macOS 示例：
# export ANDROID_SDK_ROOT="$HOME/Android/Sdk"

# 如需：
# export ANDROID_HOME="$ANDROID_SDK_ROOT"
```

> 备注：`make_frp_dex_aar.sh` 已内置一个 Windows Git Bash 默认路径 `/c/Users/admin/AppData/Local/Android/Sdk`，如果 SDK 的真实路径不同，请修改脚本顶部的 `SDK_ROOT` 或设置环境变量覆盖。

---

## 5. 项目内与 AAR 相关的文件

- `Makefile.cross-compiles`
  - 目标 `android-aar`：
    - 调用 `gomobile bind` 生成 `frp_raw.aar`；
    - 调用 `make_frp_dex_aar.sh` 将其转换为带 `classes.dex` 的最终 `frp.aar`。

- `make_frp_dex_aar.sh`
  - 路径：仓库根目录。
  - 关键变量：
    - `SDK_ROOT`：优先用 `ANDROID_SDK_ROOT` / `ANDROID_HOME`，否则尝试默认路径。
    - `BUILD_TOOLS_VERSION`：需要与 `build-tools` 实际版本一致，例如：

      ```bash
      BUILD_TOOLS_VERSION="34.0.0"   # 如果你只安装了 33.0.2，请改为 33.0.2
      ```

  - 关键步骤：
    1. 解压 `frp_raw.aar` → 取出 `classes.jar`；
    2. 使用 `d8` 把 `classes.jar` 转成 `classes.dex`；
    3. 将 `classes.dex` 与原始 `jni/`、`res/` 等重新打包为 `frp.aar`。

---

## 6. 生成 AAR 的命令

### 6.1 仅生成 Android AAR

在 **Linux / macOS shell** 或 **Windows 的 Git Bash** 中执行（推荐）：

```bash
cd ~/go/src/github.com/jahen/frp

# 如未配置系统级 ANDROID_SDK_ROOT，需要先在当前 shell 设置一次
# export ANDROID_SDK_ROOT="/c/Users/admin/AppData/Local/Android/Sdk"

make -f Makefile.cross-compiles android-aar
```

成功后你应该在 `release/` 目录看到：

- `frp_raw.aar`：纯 gomobile 产物；
- `frp.aar`：已经包含 `classes.dex` 的最终 AAR，可直接给 Android 项目使用。

### 6.2 同时生成多平台二进制 + AAR

```bash
cd ~/go/src/github.com/jahen/frp
make -f Makefile.cross-compiles
```

- `app` 目标：构建各平台的 `frpc_xxx` / `frps_xxx`；
- `android-aar` 目标：按上文步骤生成 `frp.aar`。

---

## 7. AAR 导出 API 概览

`gomobile bind` 绑定的是以下两个 Go 包：

- `github.com/fatedier/frp/cmd/frp/frpc`
- `github.com/fatedier/frp/cmd/frp/frps`

对应导出的大致 API（Go 侧函数名，gomobile 会映射为 Java/Kotlin 方法）：

### 7.1 `cmd/frp/frpc`

- `Run(cfgFilePath string)`
- `RunDir(cfgDir string)`
- `RunContent(uid, cfgContent string) (err string)`
- `RunFile(uid, cfgFilePath string) (err string)`
- `Close(uid string) bool`
- `GetUids() string`
- `IsRunning(uid string) bool`

### 7.2 `cmd/frp/frps`

- `Run(cfgFilePath string)`
- `RunDir(cfgDir string)`
- `RunContent(uid, cfgContent string) (err string)`
- `RunFile(uid, cfgFilePath string) (err string)`
- `Close(uid string) bool`
- `GetUids() string`
- `IsRunning(uid string) bool`

> 注意：Android 端不启用 `pkg/vnet` 的真实虚拟网逻辑，在 `pkg/vnet/controller_android.go` 中提供的是无实际功能的 stub，实现相同方法签名以便插件和客户端代码顺利编译。

---

## 8. 常见问题速查

### 8.1 `"golang.org/x/mobile/bind" is not found`

- 配置 `GOPROXY` 后执行：

  ```bash
  go get golang.org/x/mobile/bind@latest
  gomobile init
  ```

### 8.2 `请先设置 ANDROID_SDK_ROOT 或 ANDROID_HOME 环境变量`

- 确认 SDK 路径，例如：`C:\Users\admin\AppData\Local\Android\Sdk`；
- 设置系统级或当前 shell 的 `ANDROID_SDK_ROOT`，或修改 `make_frp_dex_aar.sh` 顶部的 `SDK_ROOT` 默认路径。

### 8.3 找不到 `gomobile` / `gobind`

- 确认 `$GOPATH/bin` 已加入系统 `PATH`（Windows 示例：`C:\Users\admin\go\bin`）；

---

如需后续扩展（例如在 Android 上做真正的 VPN / TUN 模式），需要基于 Android 的 `VpnService` 重新设计一层流量处理，这超出了当前 AAR 的范围。当前 AAR 专注于：**在 Android 进程内以库方式启动/停止 frpc / frps，并通过内存配置驱动它们工作**。

