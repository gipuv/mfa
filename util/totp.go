package util

import "strings"

// PadBase32Secret 用于确保 base32 编码的 secret 长度为8的倍数，不足补 '='
// base32 解码要求长度必须是8的倍数，补齐后方便解码
func PadBase32Secret(secret string) string {
	// 去除空格，转大写
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))
	missing := len(secret) % 8
	if missing != 0 {
		secret += strings.Repeat("=", 8-missing) // 补齐 '='
	}
	return secret
}
