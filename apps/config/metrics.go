package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func InitMetrics() {
	meter := otel.Meter("config.scenario")
	pendingScenarios, err := meter.Int64ObservableGauge(
		"active_scenarios",
		metric.WithDescription("The number of active scenarios"))
	if err != nil {
		panic(err)
	}
	var count int64

	meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		count = 22
		observer.ObserveInt64(pendingScenarios, count)
		return nil
	}, pendingScenarios)
}
