package storage

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// UserStorage 用户数据库操作
type UserStorage struct {
	db *sql.DB
}

// NewUserStorage 创建用户存储实例
func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

// InitUserTables 初始化用户相关表
func (us *UserStorage) InitUserTables() error {
	// 创建用户表
	userTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		api_key TEXT UNIQUE,
		is_active BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := us.db.Exec(userTableQuery); err != nil {
		log.Printf("创建用户表失败: %v", err)
		return err
	}

	// 为现有的qa_records表添加user_id字段
	alterTableQuery := `
	ALTER TABLE qa_records ADD COLUMN user_id INTEGER REFERENCES users(id);`

	// 忽略错误，因为字段可能已经存在
	us.db.Exec(alterTableQuery)

	log.Println("用户表初始化成功")
	return nil
}

// CreateUser 创建新用户
func (us *UserStorage) CreateUser(username, email, password string) (interface{}, error) {
	// 检查用户名是否已存在
	if exists, err := us.UserExists(username, email); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("用户名或邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 生成API密钥
	apiKey, err := us.generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("生成API密钥失败: %v", err)
	}

	// 插入用户
	query := `
	INSERT INTO users (username, email, password_hash, api_key) 
	VALUES (?, ?, ?, ?)
	`
	result, err := us.db.Exec(query, username, email, string(hashedPassword), apiKey)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("获取用户ID失败: %v", err)
	}

	// 返回创建的用户
	userInterface, err := us.GetUserByID(int(userID))
	if err != nil {
		return nil, fmt.Errorf("获取创建的用户失败: %v", err)
	}

	user, ok := userInterface.(*User)
	if !ok {
		return nil, fmt.Errorf("用户类型断言失败")
	}

	log.Printf("用户创建成功: %s (ID: %d)", user.Username, user.ID)
	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func (us *UserStorage) GetUserByUsername(username string) (*User, error) {
	query := `
	SELECT id, username, email, password_hash, api_key, is_active, created_at, updated_at 
	FROM users WHERE username = ? AND is_active = 1
	`

	var user User
	err := us.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.APIKey, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return &user, nil
}

// GetUserByID 根据ID获取用户
func (us *UserStorage) GetUserByID(id int) (interface{}, error) {
	query := `
	SELECT id, username, email, password_hash, api_key, is_active, created_at, updated_at 
	FROM users WHERE id = ? AND is_active = 1
	`

	var user User
	err := us.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.APIKey, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return &user, nil
}

// ValidateUser 验证用户登录
func (us *UserStorage) ValidateUser(username, password string) (interface{}, error) {
	user, err := us.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("密码错误")
	}

	return user, nil
}

// UserExists 检查用户是否存在
func (us *UserStorage) UserExists(username, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ? OR email = ?`

	var count int
	err := us.db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("检查用户存在性失败: %v", err)
	}

	return count > 0, nil
}

// UpdateUserLastLogin 更新用户最后登录时间
func (us *UserStorage) UpdateUserLastLogin(userID int) error {
	query := `UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := us.db.Exec(query, userID)
	if err != nil {
		log.Printf("更新用户最后登录时间失败: %v", err)
		return err
	}

	return nil
}

// generateAPIKey 生成API密钥
func (us *UserStorage) generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetAllUsers 获取所有用户（管理员功能）
func (us *UserStorage) GetAllUsers() (interface{}, error) {
	query := `
	SELECT id, username, email, api_key, is_active, created_at, updated_at 
	FROM users ORDER BY created_at DESC
	`

	rows, err := us.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email,
			&user.APIKey, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			log.Printf("扫描用户记录失败: %v", err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
