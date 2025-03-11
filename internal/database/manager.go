package database

import (
	"database/sql"
	"strings"
	"fmt"

	"golang-loki-adapter.local/pkg/models"
	_ "github.com/go-sql-driver/mysql"
)

type DBManager struct {
	db  *sql.DB
	cfg *models.Config
}

func NewDBManager(cfg *models.Config) (*DBManager, error) {
	dsn := fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &DBManager{db: db, cfg: cfg}, nil
}

func (m *DBManager) ProcessQueue() ([]models.QueueRecord, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT ID, DATA FROM %s WHERE DATA IS NOT NULL AND ATTEMPTS < ? LIMIT %d",
		m.cfg.QueueTable, m.cfg.Loki.BatchSize)

	rows, err := tx.Query(query, m.cfg.Loki.Retries)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	var records []models.QueueRecord
	for rows.Next() {
		var r models.QueueRecord
		if err := rows.Scan(&r.ID, &r.Data); err != nil {
			tx.Rollback()
			return nil, err
		}
		records = append(records, r)
	}

	// Помечаем записи как обрабатываемые
	if err := m.markAsProcessing(tx, records); err != nil {
		tx.Rollback()
		return nil, err
	}

	return records, tx.Commit()
}

func (m *DBManager) markAsProcessing(tx *sql.Tx, records []models.QueueRecord) error {
	if len(records) == 0 {
		return nil
	}

	ids := make([]interface{}, len(records))
	for i, r := range records {
		ids[i] = r.ID
	}

	if len(ids) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE %s 
		SET ATTEMPTS = ATTEMPTS + 1, 
		    LAST_ATTEMPT_AT = NOW() 
		WHERE ID IN (?%s)`,
		m.cfg.QueueTable,
		strings.Repeat(",?", len(ids)-1))

	_, err := tx.Exec(query, ids...)
	return err
}

func (m *DBManager) DeleteProcessed(records []models.QueueRecord) error {
	if len(records) == 0 {
		return nil
	}

	ids := make([]interface{}, len(records))
	for i, r := range records {
		ids[i] = r.ID
	}

	if len(ids) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		DELETE FROM %s 
		WHERE ID IN (?%s)`,
		m.cfg.QueueTable,
		strings.Repeat(",?", len(ids)-1))

	_, err := m.db.Exec(query, ids...)
	return err
}

func (m *DBManager) Close() {
	m.db.Close()
}
