package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStorage SQLite存储实现
type SQLiteStorage struct {
	db     *sql.DB
	dbPath string
}

// ExecutionTrace 执行追踪记录
type ExecutionTrace struct {
	ID          string                 `json:"id"`
	ServiceName string                 `json:"service_name"`
	FlowName    string                 `json:"flow_name"`
	Method      string                 `json:"method"`
	Status      string                 `json:"status"` // running, completed, failed
	StartTime   time.Time              `json:"start_time"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Steps       []ExecutionStep        `json:"steps"`
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output"`
	Error       string                 `json:"error,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ExecutionStep 执行步骤
type ExecutionStep struct {
	StepNumber      int                    `json:"step_number"`
	StepName        string                 `json:"step_name"`
	Method          string                 `json:"method"`
	Status          string                 `json:"status"`
	StartTime       time.Time              `json:"start_time"`
	EndTime         *time.Time             `json:"end_time,omitempty"`
	Duration        time.Duration          `json:"duration"`
	Input           map[string]interface{} `json:"input"`
	Output          map[string]interface{} `json:"output"`
	Error           string                 `json:"error,omitempty"`
	LogicFlow       []string               `json:"logic_flow"`
	BusinessContext string                 `json:"business_context"`
}

// PaginatedResult 分页结果
type PaginatedResult struct {
	Traces     []ExecutionTrace `json:"traces"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// QueryOptions 查询选项
type QueryOptions struct {
	Page        int        `json:"page"`
	PageSize    int        `json:"page_size"`
	Status      string     `json:"status"`
	Method      string     `json:"method"`
	ServiceName string     `json:"service_name"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	UserID      string     `json:"user_id"`
	OrderBy     string     `json:"order_by"`
	OrderDesc   bool       `json:"order_desc"`
}

// NewSQLiteStorage 创建SQLite存储
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开SQLite数据库失败: %v", err)
	}

	storage := &SQLiteStorage{
		db:     db,
		dbPath: dbPath,
	}

	if err := storage.initTables(); err != nil {
		return nil, fmt.Errorf("初始化数据库表失败: %v", err)
	}

	log.Printf("✅ SQLite存储初始化成功: %s", dbPath)
	return storage, nil
}

// initTables 初始化数据库表
func (s *SQLiteStorage) initTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS execution_traces (
		id TEXT PRIMARY KEY,
		service_name TEXT NOT NULL,
		flow_name TEXT NOT NULL,
		method TEXT NOT NULL,
		status TEXT NOT NULL,
		start_time DATETIME NOT NULL,
		end_time DATETIME,
		duration INTEGER DEFAULT 0,
		steps TEXT,
		input TEXT,
		output TEXT,
		error TEXT,
		user_id TEXT,
		request_id TEXT,
		tags TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_traces_service_name ON execution_traces(service_name);
	CREATE INDEX IF NOT EXISTS idx_traces_method ON execution_traces(method);
	CREATE INDEX IF NOT EXISTS idx_traces_status ON execution_traces(status);
	CREATE INDEX IF NOT EXISTS idx_traces_created_at ON execution_traces(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_traces_start_time ON execution_traces(start_time DESC);

	CREATE TABLE IF NOT EXISTS execution_statistics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_name TEXT NOT NULL,
		method TEXT NOT NULL,
		date DATE NOT NULL,
		total_executions INTEGER DEFAULT 0,
		successful_executions INTEGER DEFAULT 0,
		failed_executions INTEGER DEFAULT 0,
		avg_duration INTEGER DEFAULT 0,
		min_duration INTEGER DEFAULT 0,
		max_duration INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(service_name, method, date)
	);
	`

	if _, err := s.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}

	return nil
}

// SaveTrace 保存执行追踪
func (s *SQLiteStorage) SaveTrace(trace *ExecutionTrace) error {
	trace.UpdatedAt = time.Now()
	if trace.CreatedAt.IsZero() {
		trace.CreatedAt = time.Now()
	}

	// 序列化复杂字段
	stepsJSON, _ := json.Marshal(trace.Steps)
	inputJSON, _ := json.Marshal(trace.Input)
	outputJSON, _ := json.Marshal(trace.Output)
	tagsJSON, _ := json.Marshal(trace.Tags)

	insertSQL := `
	INSERT OR REPLACE INTO execution_traces 
	(id, service_name, flow_name, method, status, start_time, end_time, duration, 
	 steps, input, output, error, user_id, request_id, tags, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var endTime *time.Time
	if trace.EndTime != nil {
		endTime = trace.EndTime
	}

	_, err := s.db.Exec(insertSQL,
		trace.ID, trace.ServiceName, trace.FlowName, trace.Method, trace.Status,
		trace.StartTime, endTime, int64(trace.Duration),
		string(stepsJSON), string(inputJSON), string(outputJSON), trace.Error,
		trace.UserID, trace.RequestID, string(tagsJSON),
		trace.CreatedAt, trace.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("保存执行追踪失败: %v", err)
	}

	// 更新统计信息
	s.updateStatistics(trace)

	return nil
}

// GetTraces 分页查询执行追踪
func (s *SQLiteStorage) GetTraces(options *QueryOptions) (*PaginatedResult, error) {
	if options == nil {
		options = &QueryOptions{Page: 1, PageSize: 10}
	}

	// 构建查询条件
	whereClause, args := s.buildWhereClause(options)

	// 查询总数
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM execution_traces %s", whereClause)
	var total int
	err := s.db.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("查询总数失败: %v", err)
	}

	// 计算分页
	totalPages := (total + options.PageSize - 1) / options.PageSize
	offset := (options.Page - 1) * options.PageSize

	// 构建排序
	orderBy := "created_at DESC"
	if options.OrderBy != "" {
		direction := "ASC"
		if options.OrderDesc {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", options.OrderBy, direction)
	}

	// 查询数据
	querySQL := fmt.Sprintf(`
		SELECT id, service_name, flow_name, method, status, start_time, end_time, 
		       duration, steps, input, output, error, user_id, request_id, 
		       tags, created_at, updated_at
		FROM execution_traces %s 
		ORDER BY %s 
		LIMIT ? OFFSET ?
	`, whereClause, orderBy)

	queryArgs := append(args, options.PageSize, offset)
	rows, err := s.db.Query(querySQL, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("查询数据失败: %v", err)
	}
	defer rows.Close()

	var traces []ExecutionTrace
	for rows.Next() {
		var trace ExecutionTrace
		var endTime sql.NullTime
		var duration sql.NullInt64
		var stepsJSON, inputJSON, outputJSON, tagsJSON sql.NullString

		err := rows.Scan(
			&trace.ID, &trace.ServiceName, &trace.FlowName, &trace.Method, &trace.Status,
			&trace.StartTime, &endTime, &duration, &stepsJSON, &inputJSON, &outputJSON,
			&trace.Error, &trace.UserID, &trace.RequestID, &tagsJSON,
			&trace.CreatedAt, &trace.UpdatedAt,
		)
		if err != nil {
			log.Printf("❌ 扫描行数据失败: %v", err)
			continue
		}

		// 处理可空字段
		if endTime.Valid {
			trace.EndTime = &endTime.Time
		}
		if duration.Valid {
			trace.Duration = time.Duration(duration.Int64)
		}

		// 反序列化JSON字段
		if stepsJSON.Valid {
			json.Unmarshal([]byte(stepsJSON.String), &trace.Steps)
		}
		if inputJSON.Valid {
			json.Unmarshal([]byte(inputJSON.String), &trace.Input)
		}
		if outputJSON.Valid {
			json.Unmarshal([]byte(outputJSON.String), &trace.Output)
		}
		if tagsJSON.Valid {
			json.Unmarshal([]byte(tagsJSON.String), &trace.Tags)
		}

		traces = append(traces, trace)
	}

	return &PaginatedResult{
		Traces:     traces,
		Total:      total,
		Page:       options.Page,
		PageSize:   options.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetTraceByID 根据ID获取执行追踪
func (s *SQLiteStorage) GetTraceByID(id string) (*ExecutionTrace, error) {
	querySQL := `
		SELECT id, service_name, flow_name, method, status, start_time, end_time, 
		       duration, steps, input, output, error, user_id, request_id, 
		       tags, created_at, updated_at
		FROM execution_traces 
		WHERE id = ?
	`

	var trace ExecutionTrace
	var endTime sql.NullTime
	var duration sql.NullInt64
	var stepsJSON, inputJSON, outputJSON, tagsJSON sql.NullString

	err := s.db.QueryRow(querySQL, id).Scan(
		&trace.ID, &trace.ServiceName, &trace.FlowName, &trace.Method, &trace.Status,
		&trace.StartTime, &endTime, &duration, &stepsJSON, &inputJSON, &outputJSON,
		&trace.Error, &trace.UserID, &trace.RequestID, &tagsJSON,
		&trace.CreatedAt, &trace.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("查询执行追踪失败: %v", err)
	}

	// 处理可空字段和JSON反序列化
	if endTime.Valid {
		trace.EndTime = &endTime.Time
	}
	if duration.Valid {
		trace.Duration = time.Duration(duration.Int64)
	}
	if stepsJSON.Valid {
		json.Unmarshal([]byte(stepsJSON.String), &trace.Steps)
	}
	if inputJSON.Valid {
		json.Unmarshal([]byte(inputJSON.String), &trace.Input)
	}
	if outputJSON.Valid {
		json.Unmarshal([]byte(outputJSON.String), &trace.Output)
	}
	if tagsJSON.Valid {
		json.Unmarshal([]byte(tagsJSON.String), &trace.Tags)
	}

	return &trace, nil
}

// buildWhereClause 构建WHERE子句
func (s *SQLiteStorage) buildWhereClause(options *QueryOptions) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if options.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, options.Status)
	}

	if options.Method != "" {
		conditions = append(conditions, "method = ?")
		args = append(args, options.Method)
	}

	if options.ServiceName != "" {
		conditions = append(conditions, "service_name = ?")
		args = append(args, options.ServiceName)
	}

	if options.UserID != "" {
		conditions = append(conditions, "user_id = ?")
		args = append(args, options.UserID)
	}

	if options.StartTime != nil {
		conditions = append(conditions, "start_time >= ?")
		args = append(args, options.StartTime)
	}

	if options.EndTime != nil {
		conditions = append(conditions, "start_time <= ?")
		args = append(args, options.EndTime)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	return whereClause, args
}

// updateStatistics 更新统计信息
func (s *SQLiteStorage) updateStatistics(trace *ExecutionTrace) {
	date := trace.StartTime.Format("2006-01-02")
	duration := int64(trace.Duration)

	// 检查是否已存在统计记录
	var exists bool
	checkSQL := "SELECT EXISTS(SELECT 1 FROM execution_statistics WHERE service_name = ? AND method = ? AND date = ?)"
	s.db.QueryRow(checkSQL, trace.ServiceName, trace.Method, date).Scan(&exists)

	if exists {
		// 更新现有记录
		updateSQL := `
		UPDATE execution_statistics SET
			total_executions = total_executions + 1,
			successful_executions = successful_executions + CASE WHEN ? = 'completed' THEN 1 ELSE 0 END,
			failed_executions = failed_executions + CASE WHEN ? = 'failed' THEN 1 ELSE 0 END,
			avg_duration = (avg_duration * (total_executions - 1) + ?) / total_executions,
			min_duration = CASE WHEN ? < min_duration OR min_duration = 0 THEN ? ELSE min_duration END,
			max_duration = CASE WHEN ? > max_duration THEN ? ELSE max_duration END,
			updated_at = CURRENT_TIMESTAMP
		WHERE service_name = ? AND method = ? AND date = ?
		`
		s.db.Exec(updateSQL, trace.Status, trace.Status, duration, duration, duration, duration, duration, trace.ServiceName, trace.Method, date)
	} else {
		// 插入新记录
		insertSQL := `
		INSERT INTO execution_statistics 
		(service_name, method, date, total_executions, successful_executions, failed_executions, 
		 avg_duration, min_duration, max_duration) 
		VALUES (?, ?, ?, 1, ?, ?, ?, ?, ?)
		`
		successCount := 0
		failedCount := 0
		if trace.Status == "completed" {
			successCount = 1
		} else if trace.Status == "failed" {
			failedCount = 1
		}

		s.db.Exec(insertSQL, trace.ServiceName, trace.Method, date, successCount, failedCount, duration, duration, duration)
	}
}

// GetStatistics 获取统计信息
func (s *SQLiteStorage) GetStatistics(serviceName, method string, days int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 基础统计
	basicSQL := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN status = 'running' THEN 1 ELSE 0 END) as running,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed,
			AVG(CASE WHEN duration > 0 THEN duration ELSE NULL END) as avg_duration,
			MIN(CASE WHEN duration > 0 THEN duration ELSE NULL END) as min_duration,
			MAX(duration) as max_duration
		FROM execution_traces
		WHERE created_at >= datetime('now', '-' || ? || ' days')
	`

	args := []interface{}{days}
	whereClause := ""

	if serviceName != "" {
		whereClause += " AND service_name = ?"
		args = append(args, serviceName)
	}
	if method != "" {
		whereClause += " AND method = ?"
		args = append(args, method)
	}

	var total, completed, running, failed int
	var avgDuration, minDuration, maxDuration sql.NullFloat64

	err := s.db.QueryRow(basicSQL+whereClause, args...).Scan(
		&total, &completed, &running, &failed, &avgDuration, &minDuration, &maxDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("查询基础统计失败: %v", err)
	}

	stats["total"] = total
	stats["completed"] = completed
	stats["running"] = running
	stats["failed"] = failed

	if total > 0 {
		stats["success_rate"] = float64(completed) / float64(total) * 100
	} else {
		stats["success_rate"] = 0
	}

	if avgDuration.Valid {
		stats["avg_duration_ms"] = int(avgDuration.Float64 / 1000000)
	}
	if minDuration.Valid {
		stats["min_duration_ms"] = int(minDuration.Float64 / 1000000)
	}
	if maxDuration.Valid {
		stats["max_duration_ms"] = int(maxDuration.Float64 / 1000000)
	}

	return stats, nil
}

// CleanupOldRecords 清理过期记录
func (s *SQLiteStorage) CleanupOldRecords(retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}

	// 删除过期的执行追踪记录
	deleteSQL := "DELETE FROM execution_traces WHERE created_at < datetime('now', '-' || ? || ' days')"
	result, err := s.db.Exec(deleteSQL, retentionDays)
	if err != nil {
		return fmt.Errorf("清理过期记录失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("🧹 清理了 %d 条过期记录", rowsAffected)
	}

	// 清理过期的统计记录
	deleteStatsSQL := "DELETE FROM execution_statistics WHERE date < date('now', '-' || ? || ' days')"
	s.db.Exec(deleteStatsSQL, retentionDays)

	return nil
}

// Close 关闭数据库连接
func (s *SQLiteStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
