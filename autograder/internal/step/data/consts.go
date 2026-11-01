package data

const (
	// minRealties is the minimum number of rows in dim_realty after make setup.
	minRealties = 1000
	// minGauges is the minimum number of rows in dim_gauges after make setup.
	minGauges = 1000
	// minUsers is the minimum number of rows in dim_users after make setup.
	minUsers = 1000
	// minBindings is the minimum number of rows in fact_gauge_bindings after make setup.
	minBindings = 1000
	// minReadings is the minimum number of meter readings after make setup, counted with FINAL.
	minReadings = 1000000
	// minInvoices is the minimum number of invoices after make setup.
	minInvoices = 36000
	// minReplacedBindings is the minimum number of gauge bindings with version 2 or higher.
	// Such a binding exists only after a gauge was replaced, and it survives merges.
	minReplacedBindings = 1
	// minTariffsWithVersions is the minimum number of tariffs with at least two versions.
	minTariffsWithVersions = 1
)

const (
	// countQuery counts the rows of a table, the query is formatted with the database and the table.
	countQuery = "SELECT count() FROM %s.%s"
	// countFinalQuery counts the rows of a table with FINAL, the query is formatted with the database and the table.
	countFinalQuery = "SELECT count() FROM %s.%s FINAL"
	// replacedBindingsQuery counts the gauge bindings written after a gauge was replaced.
	replacedBindingsQuery = "SELECT count() FROM %s.fact_gauge_bindings WHERE version >= 2"
	// tariffsWithVersionsQuery counts the tariffs that have at least two versions.
	tariffsWithVersionsQuery = "SELECT count() FROM (SELECT id FROM %s.dim_tariffs FINAL GROUP BY id HAVING count() >= 2)"
)

const (
	// realtyTable is the name of the table with the objects of the company.
	realtyTable = "dim_realty"
	// gaugesTable is the name of the table with the gauges.
	gaugesTable = "dim_gauges"
	// usersTable is the name of the table with the owners.
	usersTable = "dim_users"
	// bindingsTable is the name of the table with the gauge installations.
	bindingsTable = "fact_gauge_bindings"
	// readingsTable is the name of the table with the meter readings.
	readingsTable = "fact_meter_readings"
	// invoicesTable is the name of the table with the invoices.
	invoicesTable = "fact_invoice"
)
