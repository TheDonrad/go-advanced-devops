package dbstorage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func (m *DBStorage) Save() error {
	db, err := sql.Open("pgx",
		m.Settings.DBConnString)
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Print(err.Error())
		}
	}()

	writeMetric(db, m)
	return nil
}

func writeMetric(db *sql.DB, m *DBStorage) {
	query := `INSERT INTO metrics(type, name, value)
		VALUES($1, $2, $3)
		ON CONFLICT (type, name) DO UPDATE
		SET value = $3;`
	var err error
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer func() {
		if err != nil {
			log.Println(err.Error())
			if err = tx.Rollback(); err != nil {
				log.Println(err.Error())
			}
		} else {
			if err = tx.Commit(); err != nil {
				log.Println(err.Error())
			}
		}
	}()
	ctx := context.Background()
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

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
}
