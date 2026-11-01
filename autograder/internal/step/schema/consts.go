package schema

// column describes a column the student must create.
type column struct {
	// name is the name of the column.
	name string
	// typ is the ClickHouse type of the column.
	typ string
}

// table describes a table the student must create with its engine and columns.
type table struct {
	// name is the name of the table.
	name string
	// engine is the name of the table engine as reported by system.tables.
	engine string
	// columns are the columns the table must have.
	columns []column
}

// tables lists the tables of lab 1 with their engines and columns.
var tables = []table{
	{
		name:   "dim_realty",
		engine: "MergeTree",
		columns: []column{
			{name: "id", typ: "UInt64"},
			{name: "address", typ: "String"},
			{name: "square", typ: "UInt32"},
			{name: "type", typ: "String"},
			{name: "number", typ: "String"},
			{name: "rooms", typ: "UInt8"},
			{name: "floor", typ: "UInt8"},
			{name: "porch", typ: "UInt8"},
			{name: "created_at", typ: "DateTime"},
		},
	},
	{
		name:   "dim_gauges",
		engine: "MergeTree",
		columns: []column{
			{name: "id", typ: "UInt64"},
			{name: "name", typ: "String"},
			{name: "country", typ: "String"},
			{name: "created_at", typ: "DateTime"},
		},
	},
	{
		name:   "dim_tariffs",
		engine: "ReplacingMergeTree",
		columns: []column{
			{name: "id", typ: "UInt64"},
			{name: "name", typ: "String"},
			{name: "price_day_kwph", typ: "UInt32"},
			{name: "price_night_kwph", typ: "UInt32"},
			{name: "valid_from", typ: "DateTime"},
			{name: "valid_to", typ: "DateTime"},
			{name: "created_at", typ: "DateTime"},
			{name: "updated_at", typ: "DateTime"},
		},
	},
	{
		name:   "dim_users",
		engine: "ReplacingMergeTree",
		columns: []column{
			{name: "id", typ: "UInt64"},
			{name: "fio", typ: "String"},
			{name: "phone", typ: "String"},
			{name: "phone_prev", typ: "String"},
			{name: "created_at", typ: "DateTime"},
			{name: "updated_at", typ: "DateTime"},
		},
	},
	{
		name:   "fact_gauge_bindings",
		engine: "VersionedCollapsingMergeTree",
		columns: []column{
			{name: "gauge_id", typ: "UInt64"},
			{name: "realty_id", typ: "UInt64"},
			{name: "created_at", typ: "DateTime"},
			{name: "sign", typ: "Int8"},
			{name: "version", typ: "UInt64"},
		},
	},
	{
		name:   "fact_meter_readings",
		engine: "ReplacingMergeTree",
		columns: []column{
			{name: "realty_id", typ: "UInt64"},
			{name: "user_id", typ: "UInt64"},
			{name: "tariff_id", typ: "UInt64"},
			{name: "gauge_id", typ: "UInt64"},
			{name: "reading_date", typ: "Date"},
			{name: "day_kw", typ: "UInt32"},
			{name: "night_kw", typ: "UInt32"},
			{name: "created_at", typ: "DateTime"},
		},
	},
	{
		name:   "fact_invoice",
		engine: "MergeTree",
		columns: []column{
			{name: "realty_id", typ: "UInt64"},
			{name: "user_id", typ: "UInt64"},
			{name: "tariff_id", typ: "UInt64"},
			{name: "gauge_id", typ: "UInt64"},
			{name: "day_kw", typ: "UInt32"},
			{name: "night_kw", typ: "UInt32"},
			{name: "price", typ: "UInt64"},
			{name: "deadline", typ: "Date"},
			{name: "created_at", typ: "DateTime"},
		},
	},
}

const (
	// enginesQuery selects the name and the engine of every table in the database.
	enginesQuery = "SELECT name, engine FROM system.tables WHERE database = ?"
	// columnsQuery selects the table, the name and the type of every column in the database.
	columnsQuery = "SELECT table, name, type FROM system.columns WHERE database = ?"
)
