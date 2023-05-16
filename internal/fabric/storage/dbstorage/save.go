package dbstorage

import (
	"context"
	"database/sql"
	"log"

	"goAdvancedTpl/internal/fabric/logs"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (m *DBStorage) Save() error {
	m.Mutex.RLock()
	db, err := sql.Open("pgx",
		m.Settings.DBConnString)
	if err != nil {
		return err
	}
	writeMetric(db, m)
	if err = db.Close(); err != nil {
		logs.New().Println(err.Error())
	}
	m.Mutex.RUnlock()
	return nil
}

func writeMetric(db *sql.DB, m *DBStorage) {
	query := `INSERT INTO metrics(type, name, value)
		VALUES($1, $2, $3)
		ON CONFLICT (type, name) DO UPDATE
		SET value = $3;`
	var err error

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for k, v := range m.Metrics.Gauge {
		if _, err = stmt.ExecContext(ctx, "gauge", k, v); err != nil {
			log.Println(err.Error())
		}
	}
	for k, v := range m.Metrics.Counter {
		if _, err := stmt.ExecContext(ctx, "counter", k, v); err != nil {
			log.Println(err.Error())
		}
	}

	if err = stmt.Close(); err != nil {
		logs.New().Println(err.Error())
	}

	if err = tx.Commit(); err != nil {
		logs.New().Println(err.Error())
	}

}
