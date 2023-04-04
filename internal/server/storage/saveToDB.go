package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (m *MetricStorage) saveToDB(dbConnString string) {
	conn, err := pgx.Connect(context.Background(), dbConnString)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		if err = conn.Close(context.Background()); err != nil {
			fmt.Println(err.Error())
		}
	}()

	writeMetric(conn, m)
}

func writeMetric(conn *pgx.Conn, m *MetricStorage) {
	query := `INSERT INTO metrics(type, name, value)
		VALUES($1::text, $2::text, $3)
		ON CONFLICT (type, name) DO UPDATE
		SET value = $3;`
	for k, v := range m.Gauge {
		if _, err := conn.Exec(context.Background(), query, "gauge", k, v); err != nil {
			fmt.Println(err.Error())
		}
	}
	for k, v := range m.Counter {
		if _, err := conn.Exec(context.Background(), query, "counter", k, v); err != nil {
			fmt.Println(err.Error())
		}
	}
}
