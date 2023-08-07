package sender

import (
	"context"
	"errors"
	"time"

	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/fabric/calchash"
	"goAdvancedTpl/internal/fabric/logs"
	"goAdvancedTpl/internal/fabric/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GRPCSender(addr string, met *collector.MetricsList, key string) error {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	for err != nil {
		logs.Logger().Println(err.Error())
		time.Sleep(5 * time.Second)
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	c := proto.NewMetricsClient(conn)

	err = sendMetrics(c, met, key)
	if errC := conn.Close(); err != nil {
		logs.Logger().Println(errC.Error())
	}
	return err
}

func sendMetrics(c proto.MetricsClient, metrics *collector.MetricsList, key string) error {

	var met []*proto.Metric

	for name, value := range metrics.Gauge {
		met = append(met, &proto.Metric{
			ID:    name,
			MType: "gauge",
			Value: value,
			Hash:  calchash.Calculate[float64](key, "gauge", name, value),
		})
	}

	for name, value := range metrics.Counter {
		met = append(met, &proto.Metric{
			ID:    name,
			MType: "counter",
			Delta: value,
			Hash:  calchash.Calculate[int64](key, "counter", name, value),
		})
	}

	resp, err := c.AddMetrics(context.Background(), &proto.AddMetricsRequest{
		Metric: met,
	})
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	return nil
}
