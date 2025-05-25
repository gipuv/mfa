package cmd

import (
	"database/sql"
	"fmt"

	"github.com/gipuv/mfa/database"
	"github.com/gipuv/mfa/totp"
)

// Run 根据命令行参数执行相应操作：add 添加/更新密钥，get 获取当前验证码
func Run(op, name, secret, code string) error {
	// 校验操作参数
	if op != "add" && op != "get" {
		return fmt.Errorf("请使用 -op add 或 -op get")
	}
	if name == "" {
		return fmt.Errorf("请指定 -name")
	}

	switch op {
	case "add":
		if secret == "" {
			return fmt.Errorf("add 操作必须指定 -secret")
		}

		// 格式校验
		if !isValidSecret(secret) {
			return fmt.Errorf("无效的密钥格式，请输入正确的Base32编码密钥")
		}

		// 查询名称对应的密钥是否存在
		existingSecret, err := database.GetSecretByName(name)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("查询密钥失败: %w", err)
		}

		var confirm bool
		if existingSecret != "" {
			// 存在时交互式询问是否替换
			confirm, err = askConfirm(fmt.Sprintf("名称 %q 已存在，是否替换？", name))
			if err != nil {
				return fmt.Errorf("读取用户输入失败: %w", err)
			}
			if !confirm {
				fmt.Println("取消操作")
				return nil
			}
			// 更新密钥
			err = database.SaveSecretByName(name, secret)
			if err != nil {
				return fmt.Errorf("更新密钥失败: %w", err)
			}
			fmt.Printf("成功替换名称 %q 的密钥\n", name)
		} else {
			// 不存在则新增保存
			err = database.SaveSecretByName(name, secret)
			if err != nil {
				return fmt.Errorf("保存密钥失败: %w", err)
			}
			fmt.Printf("成功保存名称 %q 的密钥\n", name)
		}

		// 保存或替换成功后，生成并打印当前验证码
		curCode, err := totp.GenerateTOTP(secret, 30)
		if err != nil {
			return fmt.Errorf("生成验证码失败: %w", err)
		}
		fmt.Printf("名称 %q 的当前验证码: %s\n", name, curCode)

	case "get":
		// 获取密钥
		secret, err := database.GetSecretByName(name)
		if err != nil {
			return fmt.Errorf("查询密钥失败: %w", err)
		}

		// 生成当前 TOTP 码
		curCode, err := totp.GenerateTOTP(secret, 30)
		if err != nil {
			return fmt.Errorf("生成验证码失败: %w", err)
		}

		fmt.Printf("名称 %q 的当前验证码: %s\n", name, curCode)

		// 如果传入用户验证码，进行验证
		if code != "" {
			valid := totp.ValidateTOTP(secret, code, 30)
			fmt.Printf("用户验证码 %q 是否有效: %v\n", code, valid)
		}
	}

	return nil
}
