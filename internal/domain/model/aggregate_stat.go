package model

// AggregatedStat is a generic structure for holding aggregated analytics data.
// For example: Key="Chrome", Value=150 or Key="2025-09-08", Value=25.
type AggregatedStat struct {
	Key   string
	Value int64
}

// AggregatedStatDetailed is a structure for holding multi-key aggregated data.
type AggregatedStatDetailed struct {
	TimeKey string
	UAKey   string
	Value   int64
}
