package store

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
)

// SQLiteStore SQLite 存储
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore 创建 SQLite 存储
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.init(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

// init 初始化数据库表
func (s *SQLiteStore) init() error {
	schema := `
	CREATE TABLE IF NOT EXISTS local_tasks (
		id              TEXT PRIMARY KEY,
		name            TEXT,
		type            TEXT DEFAULT 'shell',
		command         TEXT NOT NULL,
		args            TEXT,
		workdir         TEXT,
		env             TEXT,
		timeout         INTEGER DEFAULT 3600,
		priority        INTEGER DEFAULT 5,
		retry_count     INTEGER DEFAULT 0,
		retry_delay     INTEGER DEFAULT 60,
		max_retries     INTEGER DEFAULT 3,
		status          TEXT DEFAULT 'pending',
		exit_code       INTEGER,
		stdout          TEXT,
		stderr          TEXT,
		error           TEXT,
		machine_id      TEXT,
		group_id        TEXT,
		parent_id       TEXT,
		assigned_agent_id TEXT,
		lease_expires_at  TEXT,
		attempt_id        TEXT,
		created_at      TEXT,
		assigned_at     TEXT,
		started_at      TEXT,
		ended_at        TEXT,
		synced          INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_local_tasks_status ON local_tasks(status);
	CREATE INDEX IF NOT EXISTS idx_local_tasks_priority ON local_tasks(priority);
	`
	_, err := s.db.Exec(schema)
	return err
}

// Close 关闭数据库连接
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// Save 保存任务
func (s *SQLiteStore) Save(task *models.Task) error {
	argsJSON, _ := json.Marshal(task.Args)
	envJSON, _ := json.Marshal(task.Env)

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO local_tasks (
			id, name, type, command, args, workdir, env, timeout,
			priority, retry_count, retry_delay, max_retries,
			status, exit_code, stdout, stderr, error,
			machine_id, group_id, parent_id,
			assigned_agent_id, lease_expires_at, attempt_id,
			created_at, assigned_at, started_at, ended_at, synced
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.ID, task.Name, task.Type, task.Command, string(argsJSON), task.WorkDir, string(envJSON), task.Timeout,
		task.Priority, task.RetryCount, task.RetryDelay, task.MaxRetries,
		task.Status, task.ExitCode, task.Stdout, task.Stderr, task.Error,
		task.MachineID, task.GroupID, task.ParentID,
		task.AssignedAgentID, formatTime(task.LeaseExpiresAt), task.AttemptID,
		formatTime(task.CreatedAt), formatTime(task.AssignedAt), formatTime(task.StartedAt), formatTime(task.EndedAt),
		boolToInt(task.Synced),
	)
	return err
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Get 获取单个任务
func (s *SQLiteStore) Get(id string) (*models.Task, error) {
	row := s.db.QueryRow(`SELECT * FROM local_tasks WHERE id = ?`, id)
	return s.scanTask(row)
}

// scanner 接口用于统一处理 Row 和 Rows
type scanner interface {
	Scan(dest ...any) error
}

func (s *SQLiteStore) scanTask(row scanner) (*models.Task, error) {
	var task models.Task
	var argsJSON, envJSON string
	var leaseExpires, createdAt, assignedAt, startedAt, endedAt string
	var synced int

	err := row.Scan(
		&task.ID, &task.Name, &task.Type, &task.Command, &argsJSON, &task.WorkDir, &envJSON, &task.Timeout,
		&task.Priority, &task.RetryCount, &task.RetryDelay, &task.MaxRetries,
		&task.Status, &task.ExitCode, &task.Stdout, &task.Stderr, &task.Error,
		&task.MachineID, &task.GroupID, &task.ParentID,
		&task.AssignedAgentID, &leaseExpires, &task.AttemptID,
		&createdAt, &assignedAt, &startedAt, &endedAt, &synced,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(argsJSON), &task.Args)
	json.Unmarshal([]byte(envJSON), &task.Env)
	task.LeaseExpiresAt = parseTime(leaseExpires)
	task.CreatedAt = parseTime(createdAt)
	task.AssignedAt = parseTime(assignedAt)
	task.StartedAt = parseTime(startedAt)
	task.EndedAt = parseTime(endedAt)
	task.Synced = synced == 1

	return &task, nil
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

// ListByStatus 按状态查询任务
func (s *SQLiteStore) ListByStatus(status models.TaskStatus) ([]*models.Task, error) {
	rows, err := s.db.Query(`SELECT * FROM local_tasks WHERE status = ? ORDER BY priority, created_at`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task, err := s.scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// Delete 删除任务
func (s *SQLiteStore) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM local_tasks WHERE id = ?`, id)
	return err
}

// ListUnsynced 获取未同步的任务
func (s *SQLiteStore) ListUnsynced() ([]*models.Task, error) {
	rows, err := s.db.Query(`SELECT * FROM local_tasks WHERE synced = 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task, err := s.scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// MarkSynced 标记任务已同步
func (s *SQLiteStore) MarkSynced(id string) error {
	_, err := s.db.Exec(`UPDATE local_tasks SET synced = 1 WHERE id = ?`, id)
	return err
}
