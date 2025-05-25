package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gipuv/mfa/cmd"
	"github.com/gipuv/mfa/config"
	"github.com/gipuv/mfa/database"
)

func main() {
	// 加载配置文件
	cfg, err := config.LoadOrCreateConfig("example.config.json")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库连接（只做一次）
	if err := database.Init(cfg); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		}
	}()

	// 命令行参数
	var (
		op     string // 操作：add / get
		name   string // TOTP 名称
		secret string // TOTP 密钥（add时需要）
		code   string // 验证码（get时用于验证）
	)
	flag.StringVar(&op, "op", "", "操作：add 或 get")
	flag.StringVar(&name, "name", "", "TOTP 名称")
	flag.StringVar(&secret, "secret", "", "TOTP 密钥（add时必填）")
	flag.StringVar(&code, "code", "", "验证码（get时可选，用于验证）")
	flag.Parse()

	args := flag.Args()

	// 场景1：仅 mfa.exe，交互输入name
	if op == "" && len(args) == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入TOTP名称: ")
		n, _ := reader.ReadString('\n')
		name = strings.TrimSpace(n)
		if name == "" {
			log.Println("错误: 名称不能为空")
			os.Exit(1)
		}
		if err := cmd.HandleInteractive(name); err != nil {
			log.Printf("处理交互式输入失败: %v", err)
			os.Exit(1)
		}
		return
	}

	// 场景2：直接 mfa.exe github 形式
	if op == "" && len(args) == 1 {
		name = args[0]
		if err := cmd.HandleInteractive(name); err != nil {
			log.Printf("处理交互式输入失败: %v", err)
			os.Exit(1)
		}
		return
	}

	// 显式指定操作时，走原流程
	if err := cmd.Run(op, name, secret, code); err != nil {
		log.Printf("执行命令失败: %v", err)
		os.Exit(1)
	}
}
