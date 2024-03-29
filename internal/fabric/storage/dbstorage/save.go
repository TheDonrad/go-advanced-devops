package dbstorage

import (
	"context"
	"database/sql"

	"goAdvancedTpl/internal/fabric/logs"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (m *DBStorage) Save() error {

	db, err := sql.Open("pgx",
		m.Settings.DBConnString)
	if err != nil {
		return err
	}
	m.Mutex.Lock()
	writeMetric(db, m)
	m.Mutex.Unlock()
	if err = db.Close(); err != nil {
		logs.Logger().Println(err.Error())
	}

	return nil
}

func writeMetric(db *sql.DB, m *DBStorage) {
	query := `INSERT INTO metrics(type, name, value)
		VALUES($1, $2, $3)
		ON CONFLICT (type, name) DO UPDATE
		SET value = $3;`
	var err error

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		logs.Logger().Println(err.Error())
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		logs.Logger().Println(err.Error())
		if err := tx.Rollback(); err != nil {
			logs.Logger().Println(err.Error())
		}
		return
	}

	for k, v := range m.Metrics.Gauge {
		if _, err = stmt.ExecContext(ctx, "gauge", k, v); err != nil {
			logs.Logger().Println(err.Error())

			if err = stmt.Close(); err != nil {
				logs.Logger().Println(err.Error())
			}

			if err := tx.Rollback(); err != nil {
				logs.Logger().Println(err.Error())
			}

			return
		}
	}
	for k, v := range m.Metrics.Counter {
		if _, err := stmt.ExecContext(ctx, "counter", k, v); err != nil {
			logs.Logger().Println(err.Error())

			if err = stmt.Close(); err != nil {
				logs.Logger().Println(err.Error())
			}

			if err := tx.Rollback(); err != nil {
				logs.Logger().Println(err.Error())
			}

			return
		}
	}

	if err = stmt.Close(); err != nil {
		logs.Logger().Println(err.Error())
	}

	if err = tx.Commit(); err != nil {
		logs.Logger().Println(err.Error())
	}

}
