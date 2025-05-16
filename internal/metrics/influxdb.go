package metrics

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/yosev/coda/pkg/coda"
)

func SendStatsToInfluxDB(dsn string, stats *coda.CodaStats) error {
	client := influxdb2.NewClientWithOptions(dsn, "", influxdb2.DefaultOptions())
	defer client.Close()

	// Write API
	writeAPI := client.WriteAPIBlocking("", "")

	// Create a data point
	p := influxdb2.NewPoint(
		"coda_execution",
		map[string]string{},
		map[string]interface{}{
			"coda_runtime_total_ms":        stats.CodaRuntimeTotalMs,
			"operations_runtime_total_ms":  stats.OperationsRuntimeTotalMs,
			"operations_total":             stats.OperationsTotal,
			"operations_successful_total":  stats.OperationsSuccessfulTotal,
			"operations_failed_total":      stats.OperationsFailedTotal,
			"operations_blacklisted_total": stats.OperationsBlacklistedTotal,
			"variables_total":              stats.VariablesTotal,
			"variables_successful_total":   stats.VariablesSuccessfulTotal,
			"variables_failed_total":       stats.VariablesFailedTotal,
		},
		time.Now(),
	)

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		return fmt.Errorf("Error writing to InfluxDB: %s", err)
	} else {
		return fmt.Errorf("Data written successfully")
	}
}
