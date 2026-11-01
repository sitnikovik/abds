// Package schema provides the rows a correct lab 1 database reports from the system tables.
package schema

// Engines returns the rows of system.tables for a correct lab 1 database.
func Engines() [][]string {
	return [][]string{
		{"dim_realty", "MergeTree"},
		{"dim_gauges", "MergeTree"},
		{"dim_tariffs", "ReplacingMergeTree"},
		{"dim_users", "ReplacingMergeTree"},
		{"fact_gauge_bindings", "VersionedCollapsingMergeTree"},
		{"fact_meter_readings", "ReplacingMergeTree"},
		{"fact_invoice", "MergeTree"},
	}
}

// Columns returns the rows of system.columns for a correct lab 1 database.
func Columns() [][]string {
	return [][]string{
		{"dim_realty", "id", "UInt64"},
		{"dim_realty", "address", "String"},
		{"dim_realty", "square", "UInt32"},
		{"dim_realty", "type", "String"},
		{"dim_realty", "number", "String"},
		{"dim_realty", "rooms", "UInt8"},
		{"dim_realty", "floor", "UInt8"},
		{"dim_realty", "porch", "UInt8"},
		{"dim_realty", "created_at", "DateTime"},
		{"dim_gauges", "id", "UInt64"},
		{"dim_gauges", "name", "String"},
		{"dim_gauges", "country", "String"},
		{"dim_gauges", "created_at", "DateTime"},
		{"dim_tariffs", "id", "UInt64"},
		{"dim_tariffs", "name", "String"},
		{"dim_tariffs", "price_day_kwph", "UInt32"},
		{"dim_tariffs", "price_night_kwph", "UInt32"},
		{"dim_tariffs", "valid_from", "DateTime"},
		{"dim_tariffs", "valid_to", "DateTime"},
		{"dim_tariffs", "created_at", "DateTime"},
		{"dim_tariffs", "updated_at", "DateTime"},
		{"dim_users", "id", "UInt64"},
		{"dim_users", "fio", "String"},
		{"dim_users", "phone", "String"},
		{"dim_users", "phone_prev", "String"},
		{"dim_users", "created_at", "DateTime"},
		{"dim_users", "updated_at", "DateTime"},
		{"fact_gauge_bindings", "gauge_id", "UInt64"},
		{"fact_gauge_bindings", "realty_id", "UInt64"},
		{"fact_gauge_bindings", "created_at", "DateTime"},
		{"fact_gauge_bindings", "sign", "Int8"},
		{"fact_gauge_bindings", "version", "UInt64"},
		{"fact_meter_readings", "realty_id", "UInt64"},
		{"fact_meter_readings", "user_id", "UInt64"},
		{"fact_meter_readings", "tariff_id", "UInt64"},
		{"fact_meter_readings", "gauge_id", "UInt64"},
		{"fact_meter_readings", "reading_date", "Date"},
		{"fact_meter_readings", "day_kw", "UInt32"},
		{"fact_meter_readings", "night_kw", "UInt32"},
		{"fact_meter_readings", "created_at", "DateTime"},
		{"fact_invoice", "realty_id", "UInt64"},
		{"fact_invoice", "user_id", "UInt64"},
		{"fact_invoice", "tariff_id", "UInt64"},
		{"fact_invoice", "gauge_id", "UInt64"},
		{"fact_invoice", "day_kw", "UInt32"},
		{"fact_invoice", "night_kw", "UInt32"},
		{"fact_invoice", "price", "UInt64"},
		{"fact_invoice", "deadline", "Date"},
		{"fact_invoice", "created_at", "DateTime"},
	}
}
