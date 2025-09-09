package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// QAStorage QA记录数据库操作
type QAStorage struct {
	db *sql.DB
}

// NewQAStorage 创建新的QA存储实例
func NewQAStorage(dbPath string) (*QAStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	storage := &QAStorage{db: db}

	// 初始化数据库表
	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

// initTables 初始化数据库表
func (s *QAStorage) initTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS qa_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		question TEXT NOT NULL,
		answer TEXT DEFAULT '',
		user_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Printf("创建表失败: %v", err)
		return err
	}

	log.Println("QA数据库表初始化成功")
	return nil
}

// SaveQuestion 保存新问题，返回记录ID（支持用户关联）
func (s *QAStorage) SaveQuestion(question string, userID *int) (int, error) {
	query := `INSERT INTO qa_records (question, user_id) VALUES (?, ?)`

	result, err := s.db.Exec(query, question, userID)
	if err != nil {
		log.Printf("保存问题失败: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入ID失败: %v", err)
		return 0, err
	}

	log.Printf("问题保存成功，ID: %d, UserID: %v", id, userID)
	return int(id), nil
}

// UpdateAnswer 更新答案
func (s *QAStorage) UpdateAnswer(id int, answer string) error {
	query := `UPDATE qa_records SET answer = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	result, err := s.db.Exec(query, answer, id)
	if err != nil {
		log.Printf("更新答案失败: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取影响行数失败: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("未找到ID为%d的记录", id)
		return sql.ErrNoRows
	}

	log.Printf("答案更新成功，ID: %d", id)
	return nil
}

// GetRecord 根据ID获取记录
func (s *QAStorage) GetRecord(id int) (interface{}, error) {
	query := `SELECT id, question, answer, user_id, created_at, updated_at FROM qa_records WHERE id = ?`

	row := s.db.QueryRow(query, id)

	var record QARecord
	err := row.Scan(&record.ID, &record.Question, &record.Answer, &record.UserID, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("未找到ID为%d的记录", id)
		} else {
			log.Printf("查询记录失败: %v", err)
		}
		return nil, err
	}

	return &record, nil
}

// GetAllRecords 获取所有记录
func (s *QAStorage) GetAllRecords() (interface{}, error) {
	query := `SELECT id, question, answer, user_id, created_at, updated_at FROM qa_records ORDER BY created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("查询所有记录失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var records []QARecord
	for rows.Next() {
		var record QARecord
		err := rows.Scan(&record.ID, &record.Question, &record.Answer, &record.UserID, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			log.Printf("扫描记录失败: %v", err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

// GetRecordsByUserID 获取指定用户的问答记录
func (s *QAStorage) GetRecordsByUserID(userID int) (interface{}, error) {
	query := `SELECT id, question, answer, user_id, created_at, updated_at FROM qa_records WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		log.Printf("查询用户记录失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var records []QARecord
	for rows.Next() {
		var record QARecord
		err := rows.Scan(&record.ID, &record.Question, &record.Answer, &record.UserID, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			log.Printf("扫描记录失败: %v", err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

// GetDB 获取数据库连接（用于用户存储）
func (s *QAStorage) GetDB() *sql.DB {
	return s.db
}

// Close 关闭数据库连接
func (s *QAStorage) Close() error {
	return s.db.Close()
}
