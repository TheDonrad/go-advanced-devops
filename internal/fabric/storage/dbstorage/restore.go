package dbstorage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func (m *DBStorage) Restore() {
	conn, err := pgx.Connect(context.Background(), m.Settings.DBConnString)
	if err != nil {
		log.Println(err.Error())
	}

	defer func() {
		if err = conn.Close(context.Background()); err != nil {
			log.Print(err.Error())
		}
	}()

	rows, err := conn.Query(context.Background(), "SELECT name, type, value FROM metrics")
	if err != nil {
		createTable(conn)
		log.Println(err.Error())
	}
	defer rows.Close()

	type met struct {
		Name    string
		MetType string
		Value   float64
	}

	for rows.Next() {
		var r met
		err = rows.Scan(&r.Name, &r.MetType, &r.Value)
		if err != nil {
			log.Println(err)
			return
		}
		switch r.MetType {
		case "gauge":
			m.Metrics.Gauge[r.Name] = r.Value
		default:
			m.Metrics.Counter[r.Name] = int64(r.Value)
		}
	}
}

func createTable(conn *pgx.Conn) {

	query := `CREATE TABLE IF NOT EXISTS public.metrics
			(	
			type text COLLATE pg_catalog.default NOT NULL,	
			name text COLLATE pg_catalog.default NOT NULL,
			value double precision NOT NULL,
			CONSTRAINT main UNIQUE (name, type)
			)`

	if _, err := conn.Exec(context.Background(), query); err != nil {
		log.Println(err.Error())
	}

}
