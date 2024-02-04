<p align="center">
  ![logo_kw](https://github.com/kunwu2023/kunwu/assets/31091395/956c21f3-f829-49b3-8358-d2e3ebced65b)
</p>
<h1 align="center"> KunWu </h1>
<p align="center">
<img alt="Static Badge" src="https://img.shields.io/badge/Build-v0.1.0-blue">
<img alt="Static Badge" src="https://img.shields.io/badge/License-MIT-green">

<p align="center"> 默安科技打造的新一代 WebShell 检测工具「<b>KunWu</b>」</p>
<p align="center"> 集<b>模糊规则</b>、<b>污点分析模拟执行</b>、<b>机器学习</b>三种高效检测策略，精准无误地发现 WebShell 风险</p>


## 🚀 快速开始

KunWu 支持 GUI 客户端、CLI 工具、在线 WebShell 检测。

cli编译：go build -ldflags="-w -s"

gui编译：

1. 后端编译cd gui_go/go build -ldflags="-w -s"cd ..

2. 前端编译修改src/background.js文件中goPath变量路径为编译后的二进制文件

npm inpm run electron:build

### GUI 客户端使用

#### mac 端

下载 [KunWu 客户端](https://github.com/kunwu2023/kunwu/releases/tag/0.1.0) ，将app拖拽到Applications目录下。

![image](https://github.com/kunwu2023/kunwu/assets/131849947/3dc40a5e-8ef8-4452-a04b-fafdeab20c9e)

> ![NOTE]
> 
> 如果提示无法验证开发者，请在设置种选择仍要打开。

![image](https://github.com/kunwu2023/kunwu/assets/131849947/494fe0ad-b474-4480-9d31-4aba7503482f)

![image](https://github.com/kunwu2023/kunwu/assets/131849947/e834aa87-e3ff-429d-858a-1bb24a68b212)


#### Windows 端

下载 [KunWu 客户端](https://github.com/kunwu2023/kunwu/releases/tag/0.1.0) ，解压后运行「昆吾WebShell检测.exe」文件。


#### 开始扫描

KunWu 支持**快速扫描**、**本地扫描**、**远程扫描**；上传或拖拽文件开始扫描！

- 🛫 快速扫描：执行临时的扫描任务，在客户端重启后数据将会清空。
- 🚄 本地扫描：执行本地扫描任务，在任务列表中生成扫描任务记录。
- 🛸 远程扫描：通过SFTP读取远程文件进行扫描。

![image](https://github.com/kunwu2023/kunwu/assets/131849947/083351e2-6139-49eb-9c56-883ee2797612)

### CLI 工具使用

下载 [KunWu CLI For Win](https://github.com/kunwu2023/kunwu/releases/download/0.1.0/kw_0.1.2_windows_amd64_cli.exe) \ [KunWu CLI For Mac](https://github.com/kunwu2023/kunwu/releases/download/0.1.0/kw_0.1.2_macOS_arm64_cli) \ [KunWu CLI For Linux](https://github.com/kunwu2023/kunwu/releases/download/0.1.0/kw_0.1.2_linux_amd64_cli)。

#### 使用 CLI 工具快速扫描本地文件

```bash
$ kw -file /home/sample.php
File path: /home/sample.php
Cloud scan: false
filter normal files: true
--------------------------start----------------------------
Local engine scanning...
+------------------------+----------+------+
|        文件路径         | 检出引擎 | 结果  |
+------------------------+----------+------+
| /home/sample.php | 本地引擎 | 恶意 |
+------------------------+----------+------+
--------------------------end----------------------------
```

Linux 环境使用截图👇

![image](https://github.com/kunwu2023/kunwu/assets/131849947/002daae6-f86e-4643-b4e5-bdb09ecb3178)

#### 高级选项

**云端引擎检测（-cloud）：该选项开启后，会将本地无法检出的文件上传到云端进行扫描**

```bash
$ kw -cloud -file /home/sample.php
```

## 交流反馈:

![1541701248885_ pic](https://github.com/kunwu2023/kunwu/assets/20842613/5f8486d6-6f83-4df8-8d3b-b4550b5aef02)

![image](https://github.com/kunwu2023/kunwu/assets/131849947/67bf7658-5c2a-4fe7-91ad-92cc37ccdb3a)


## Stargazers over time

[![Stargazers over time](https://starchart.cc/kunwu2023/kunwu.svg)](https://starchart.cc/kunwu2023/kunwu)
