# mfa 🔐

一个简单的命令行 TOTP（动态口令）生成和管理工具。

> For the English documentation, please refer to [README\_en.md](README_en.md) 📄

## 安装 🚀

```bash
go install github.com/gipuv/mfa@latest
```

安装成功后，`mfa` 可执行文件会放在你的 `$GOPATH/bin` 或 `$HOME/go/bin` 目录下。

请确保该目录已加入你的环境变量 `PATH`，以便全局使用：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## 使用示例 💡

### 添加或更新密钥 🔑

```bash
mfa -op add -name github -secret JBSWY3DPEHPK3PXP
```

如果名称已存在，程序会提示是否替换密钥。

### 获取当前验证码 🎫

```bash
mfa -op get -name github
```

### 交互模式 🤝

```bash
mfa github
```

或者

```bash
mfa
```

此时程序会提示输入名称和密钥。

## 数据库文件预览工具 🧰

如果需要查看或管理 `.db` 文件（SQLite 数据库），推荐使用以下工具：

🔍 **DB Browser for SQLite**
官网地址：[https://sqlitebrowser.org/dl/](https://sqlitebrowser.org/dl/)
下载后打开 `.db` 文件，即可可视化管理数据库内容。

## 备注 📝

* 密钥必须为合法的 Base32 编码字符串。
* 默认 TOTP 码有效期为 30 秒。

---