package gojinn

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func (r *Gojinn) setupDB() error {
	if r.DBDriver == "" || r.DBDSN == "" {
		return nil
	}

	db, err := sql.Open(r.DBDriver, r.DBDSN)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	maxConns := r.PoolSize
	if maxConns > 20 {
		maxConns = 20
	}

	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	r.db = db
	r.logger.Info("host database connection pool established",
		zap.String("driver", r.DBDriver),
		zap.Int("max_conns", maxConns))

	return nil
}

func (r *Gojinn) executeQueryToJSON(query string) ([]byte, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not configured on host")
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)

	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	return json.Marshal(tableData)
}
