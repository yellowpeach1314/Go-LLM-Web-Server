# Go 环境安装指南

## macOS 安装 Go

### 方法一：使用 Homebrew（推荐）

1. 安装 Homebrew（如果还没有安装）：


2. 使用 Homebrew 安装 Go：


### 方法二：官方安装包

1. 访问 Go 官网：https://golang.org/dl/
2. 下载适用于 macOS 的安装包
3. 双击安装包并按照提示安装

## 验证安装

安装完成后，打开终端运行：



如果显示版本信息，说明安装成功。

## 设置环境变量

在 `~/.zshrc` 文件中添加：



然后运行：


## 运行项目

1. 进入项目目录：


2. 下载依赖：


3. 运行服务器：


4. 打开浏览器访问：http://localhost:8080

## 常见问题

### 问题1：go: command not found
- 解决方案：确保Go已正确安装并添加到PATH环境变量中

### 问题2：端口被占用
- 解决方案：修改main.go中的端口号，或者停止占用8080端口的其他程序

### 问题3：依赖下载失败 - 网络超时
当你看到类似 `dial tcp 142.250.198.81:443: i/o timeout` 的错误时，这是网络连接问题。

**解决方案1：设置中国代理（推荐）**
```bash
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn
```

**解决方案2：使用阿里云代理**
```bash
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

**解决方案3：使用七牛云代理**
```bash
go env -w GOPROXY=https://goproxy.qiniu.com,direct
```

**解决方案4：禁用模块验证（临时方案）**
```bash
go env -w GOSUMDB=off
```

设置完成后，重新运行：
```bash
go mod tidy
```

**验证代理设置：**
```bash
go env GOPROXY
go env GOSUMDB
```