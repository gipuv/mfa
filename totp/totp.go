package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/gipuv/mfa/util"
)

// GenerateTOTP 生成基于时间的一次性密码（TOTP）
// secret: base32 编码的密钥字符串
// timestep: 时间步长，单位秒，常用30秒
// 返回6位数字验证码，或错误
func GenerateTOTP(secret string, timestep int64) (string, error) {
	// 先补齐 Base32 密钥，确保长度是8的倍数（Base32编码要求）
	secret = util.PadBase32Secret(secret)

	// 使用无填充的 Base32 解码成字节数组（密钥的二进制形式）
	secretKey, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("Base32解码失败: %w", err)
	}

	// 获取当前Unix时间戳（秒），除以步长，得到计数器
	counter := time.Now().Unix() / timestep

	// 把计数器转成8字节的 大端字节数组（big-endian）
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(counter))

	// 用HMAC-SHA1算法，以secretKey为key，对计数器的字节数组进行加密
	h := hmac.New(sha1.New, secretKey)
	h.Write(buf[:])
	hash := h.Sum(nil)

	// 动态截断算法：
	// 取 hash 的最后一个字节的低4位作为偏移量 offset
	offset := hash[len(hash)-1] & 0x0F

	// 从 offset 开始取4个字节，组合成一个31位的整数（最高位去掉，防止负数）
	binCode := (uint32(hash[offset])&0x7F)<<24 |
		(uint32(hash[offset+1])&0xFF)<<16 |
		(uint32(hash[offset+2])&0xFF)<<8 |
		(uint32(hash[offset+3]) & 0xFF)

	// 对 10^6 取余，得到6位数字验证码
	code := binCode % 1000000

	// 格式化成6位字符串，不足补0返回
	return fmt.Sprintf("%06d", code), nil
}

// ValidateTOTP 验证用户输入的验证码是否正确，允许当前时间步长前后各一个步长的误差
// secret: base32 密钥
// code: 用户输入的6位验证码
// timestep: 时间步长，单位秒
// 返回是否验证通过
func ValidateTOTP(secret, code string, timestep int64) bool {
	// 验证当前步长和前后各一个步长内的验证码，容忍时间误差
	for i := -1; i <= 1; i++ {
		validCode, err := GenerateTOTPWithTime(secret, timestep, time.Now().Add(time.Duration(i)*time.Duration(timestep)*time.Second))
		if err == nil && validCode == code {
			return true
		}
	}
	return false
}

// GenerateTOTPWithTime 生成指定时间点的 TOTP，功能同 GenerateTOTP，但支持自定义时间
// secret: base32 编码密钥
// timestep: 时间步长秒数
// t: 指定时间点
// 返回6位验证码或错误
func GenerateTOTPWithTime(secret string, timestep int64, t time.Time) (string, error) {
	// 同 GenerateTOTP，只不过计数器改为 t.Unix()/timestep
	secret = util.PadBase32Secret(secret)
	secretKey, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("Base32解码失败: %w", err)
	}

	counter := t.Unix() / timestep

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(counter))

	h := hmac.New(sha1.New, secretKey)
	h.Write(buf[:])
	hash := h.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F
	binCode := (uint32(hash[offset])&0x7F)<<24 |
		(uint32(hash[offset+1])&0xFF)<<16 |
		(uint32(hash[offset+2])&0xFF)<<8 |
		(uint32(hash[offset+3]) & 0xFF)

	codeNum := binCode % 1000000
	return fmt.Sprintf("%06d", codeNum), nil
}
