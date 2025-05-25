package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// Path 配置 & .db 目录
	Path = "data"

	// 内置的默认配置
	defaultConfig = `{
    "type": "sqlite",
    "host": "localhost",
    "port": "3306",
    "user": "root",
    "password": "",
    "database": "",
    "db_file": "mfa.db",
    "table_prefix": "",
    "max_idle_conn": 10,
    "max_open_conn": 100,
    "conn_max_lifetime": 3600
  }`
)

// Config 定义数据库配置结构体，用于支持 MySQL 或 SQLite
type Config struct {
	Type            string `json:"type"`              // 数据库类型: "mysql" 或 "sqlite"，默认 "sqlite"
	Host            string `json:"host"`              // MySQL: 数据库主机地址
	Port            string `json:"port"`              // MySQL: 数据库端口号
	User            string `json:"user"`              // MySQL: 用户名
	Password        string `json:"password"`          // MySQL: 密码
	Database        string `json:"database"`          // MySQL: 数据库名（必填）
	DBFile          string `json:"db_file"`           // SQLite: mfa.db
	TablePrefix     string `json:"table_prefix"`      // 表名前缀（可选）
	MaxIdleConn     int    `json:"max_idle_conn"`     // 数据库连接池中的最大空闲连接数
	MaxOpenConn     int    `json:"max_open_conn"`     // 数据库连接池中的最大打开连接数
	ConnMaxLifetime int    `json:"conn_max_lifetime"` // 连接的最大寿命, 超过这个时间会被关闭
}

// getConfigPath 查找配置文件路径的函数
// 优先级顺序:
// 1. 读取环境变量 CONFIG_PATH 指定的路径
// 2. 尝试可执行文件所在目录下的 config.json
// 3. 尝试当前工作目录下的 config.json
func getConfigPath() (string, error) {
	// 优先检查环境变量 CONFIG_PATH
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
	}

	// 可执行文件目录的 config 子目录
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		configPath := filepath.Join(exeDir, Path, "config.json")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}

	// 当前工作目录的 config 子目录
	if wd, err := os.Getwd(); err == nil {
		configPath := filepath.Join(wd, Path, "config.json")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}

	// 都未找到，返回错误
	return "", fmt.Errorf("未找到配置文件")
}

// LoadOrCreateConfig 加载配置文件
// 如果配置文件不存在，将根据示例文件（examplePath）创建一个新的配置文件，并提示用户修改后重新运行程序
//
// 参数:
// - examplePath: 示例配置文件路径，通常为 example.config.json
//
// 返回:
// - *Config: 解析后的配置结构体指针
// - error: 错误信息
func LoadOrCreateConfig(examplePath string) (*Config, error) {
	// 获取配置文件路径（优先读取 CONFIG_PATH 环境变量，其次查找当前可执行文件所在目录，最后查找当前工作目录）
	configPath, err := getConfigPath()
	if err != nil {
		// configPath 找不到时，可直接默认生成 exe 目录下的 path/config.json
		if exePath, err2 := os.Executable(); err2 == nil {
			exeDir := filepath.Dir(exePath)
			configPath = filepath.Join(exeDir, Path, "config.json")
			// 确保 config 目录存在
			os.MkdirAll(filepath.Dir(configPath), 0755)
		} else {
			return nil, fmt.Errorf("未找到配置文件，且无法确定默认路径: %w", err)
		}
	}

	// 检查配置文件是否存在
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		if err = os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
			return nil, fmt.Errorf("创建默认配置文件失败: %w", err)
		}

		// 提示用户修改配置文件后重新运行程序
		fmt.Printf("未找到配置文件 %s，已自动创建默认配置文件，内容如下：\n%s\n", configPath, defaultConfig)

		// 这里不退出，直接加载刚创建的配置文件
		// os.Exit(0)
	}

	// 打开已存在的配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件 %s: %w", configPath, err)
	}
	defer file.Close()

	// 使用 JSON 解码配置文件内容到 Config 结构体
	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件 %s 失败: %w", configPath, err)
	}

	// 补充默认配置（如未指定数据库类型，则默认为 sqlite）
	if cfg.Type == "" {
		cfg.Type = "sqlite"
	}

	if cfg.MaxIdleConn == 0 {
		cfg.MaxIdleConn = 10
	}
	if cfg.MaxOpenConn == 0 {
		cfg.MaxOpenConn = 100
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = 3600 // 1小时，单位秒
	}

	return &cfg, nil
}
