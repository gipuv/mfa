package cmd

import (
	"bufio"
	"encoding/base32"
	"fmt"
	"os"
	"strings"

	"github.com/gipuv/mfa/database"
	"github.com/gipuv/mfa/totp"
)

// isValidSecret 简单验证secret是否为合法的Base32字符串
func isValidSecret(secret string) bool {
	secret = strings.ToUpper(secret)
	_, err := base32.StdEncoding.DecodeString(secret)
	return err == nil
}

// askConfirm 交互式确认，提示用户输入 Y 或 N
func askConfirm(prompt string) (bool, error) {
	fmt.Print(prompt + " (Y/N): ")
	reader := bufio.NewReader(os.Stdin)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	resp = strings.TrimSpace(strings.ToUpper(resp))
	return resp == "Y" || resp == "YES", nil
}

// HandleInteractive 根据给定的名称查找TOTP密钥，若不存在则提示用户输入并保存密钥，
// 最后生成并打印当前验证码。增加secret格式校验避免无效输入。
func HandleInteractive(name string) error {
	secret, err := database.GetSecretByName(name)
	if err != nil {
		if strings.Contains(err.Error(), "未找到") || strings.Contains(err.Error(), "no rows") {
			fmt.Printf("名称 %q 不存在，请输入TOTP密钥: ", name)
			reader := bufio.NewReader(os.Stdin)
			s, _ := reader.ReadString('\n')
			secret = strings.TrimSpace(s)
			if secret == "" {
				return fmt.Errorf("密钥不能为空")
			}

			if !isValidSecret(secret) {
				return fmt.Errorf("无效的密钥格式，请输入正确的Base32编码密钥")
			}

			if err = database.SaveSecretByName(name, secret); err != nil {
				return fmt.Errorf("保存密钥失败: %w", err)
			}
		} else {
			return fmt.Errorf("查询密钥失败: %w", err)
		}
	} else {
		// 已存在的secret也验证下格式，避免数据库中存了无效的
		if !isValidSecret(secret) {
			return fmt.Errorf("数据库中密钥格式无效，请检查")
		}
	}

	code, err := totp.GenerateTOTP(secret, 30)
	if err != nil {
		return fmt.Errorf("生成验证码失败: %w", err)
	}

	fmt.Printf("名称 %q 的当前验证码: %s\n", name, code)
	return nil
}
