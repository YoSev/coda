package coda

type CodaStats struct {
	CodaRuntimeTotalMs float64 `json:"coda_runtime_total_ms" yaml:"coda_runtime_total_ms"`

	OperationsRuntimeTotalMs   float64 `json:"operations_runtime_total_ms" yaml:"operations_runtime_total_ms"`
	OperationsTotal            float64 `json:"operations_total" yaml:"operations_total"`
	OperationsSuccessfulTotal  float64 `json:"operations_successful_total" yaml:"operations_successful_total"`
	OperationsFailedTotal      float64 `json:"operations_failed_total" yaml:"operations_failed_total"`
	OperationsBlacklistedTotal float64 `json:"operations_blacklisted_total" yaml:"operations_blacklisted_total"`

	VariablesTotal           float64 `json:"variables_total" yaml:"variables_total"`
	VariablesFailedTotal     float64 `json:"variables_failed_total" yaml:"variables_failed_total"`
	VariablesSuccessfulTotal float64 `json:"variables_successful_total" yaml:"variables_successful_total"`
}

func (c *Coda) newStats() *CodaStats {
	return &CodaStats{
		CodaRuntimeTotalMs: 0,

		OperationsRuntimeTotalMs:   0,
		OperationsTotal:            0,
		OperationsSuccessfulTotal:  0,
		OperationsFailedTotal:      0,
		OperationsBlacklistedTotal: 0,

		VariablesTotal:           0,
		VariablesFailedTotal:     0,
		VariablesSuccessfulTotal: 0,
	}
}
