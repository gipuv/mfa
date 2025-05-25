package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gipuv/mfa/config"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

var (
	db     *sql.DB // 全局数据库连接
	dbType string  // 全局数据库类型(mysql/sqlite)
)

// Init 初始化数据库连接并创建表，程序启动时调用一次
func Init(cfg *config.Config) error {
	var err error

	// 连接数据库，赋值给包变量 db
	db, err = openDB(cfg)
	if err != nil {
		return err
	}

	// 记录数据库类型
	dbType = cfg.Type

	// 创建数据表
	if err = createTableIfNotExists(); err != nil {
		db.Close()
		return err
	}

	return nil
}

// Close 关闭数据库连接，程序退出时调用
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// openDB 连接数据库，支持 mysql 或 sqlite，内部使用
func openDB(cfg *config.Config) (*sql.DB, error) {
	var (
		dsn    string
		driver string
		err    error
	)

	switch cfg.Type {
	case "mysql":
		driver = "mysql"
		if cfg.Database == "" {
			return nil, fmt.Errorf("MySQL 配置缺少数据库名")
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	case "sqlite":
		driver = "sqlite"

		dbFile := "mfa.db"
		if cfg.DBFile != "" {
			dbFile = cfg.DBFile
		} else if cfg.Database != "" {
			dbFile = cfg.Database
		}

		fullPath := filepath.Join(config.Path, dbFile)
		if err = os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return nil, fmt.Errorf("创建 SQLite 数据库目录失败: %w", err)
		}

		dsn = fullPath

	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Type)
	}

	db, err = sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	return db, nil
}

// createTableIfNotExists 创建 mfa_secret 表，使用包变量 db 和 dbType
func createTableIfNotExists() error {
	var (
		query string
		err   error
	)

	switch dbType {
	case "mysql":
		query = `
        CREATE TABLE IF NOT EXISTS mfa_secret (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE,
            secret VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`
	case "sqlite":
		query = `
        CREATE TABLE IF NOT EXISTS mfa_secret (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            secret TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`
	default:
		return fmt.Errorf("不支持的数据库类型: %s", dbType)
	}

	_, err = db.Exec(query)
	return err
}

// saveSecret 保存密钥，事务处理，存在则更新，不存在则插入
func saveSecret(name, secret string) error {
	var err error
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var existing string
	err = tx.QueryRow("SELECT secret FROM mfa_secret WHERE name = ?", name).Scan(&existing)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		_, err = tx.Exec("INSERT INTO mfa_secret (name, secret) VALUES (?, ?)", name, secret)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec("UPDATE mfa_secret SET secret = ? WHERE name = ?", secret, name)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// getSecret 根据名称获取密钥
func getSecret(name string) (string, error) {
	var (
		secret string
		err    error
	)
	err = db.QueryRow("SELECT secret FROM mfa_secret WHERE name = ?", name).Scan(&secret)
	return secret, err
}

// GetSecretByName 通过名称获取密钥（对外接口）
func GetSecretByName(name string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("数据库未初始化")
	}
	return getSecret(name)
}

// SaveSecretByName 通过名称保存密钥（对外接口）
func SaveSecretByName(name, secret string) error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return saveSecret(name, secret)
}
