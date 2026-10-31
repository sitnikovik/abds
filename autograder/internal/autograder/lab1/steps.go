// Package lab1 builds the autograder checks of lab 1.
package lab1

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/clickhouse"
	"github.com/sitnikovik/abds/autograder/internal/config"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/data"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/gauges"
	"github.com/sitnikovik/abds/autograder/internal/step/invoices"
	"github.com/sitnikovik/abds/autograder/internal/step/ping"
	"github.com/sitnikovik/abds/autograder/internal/step/readings"
	"github.com/sitnikovik/abds/autograder/internal/step/realty"
	"github.com/sitnikovik/abds/autograder/internal/step/schema"
	"github.com/sitnikovik/abds/autograder/internal/step/tariffs"
	"github.com/sitnikovik/abds/autograder/internal/step/users"
)

// Steps returns the checks of lab 1 in the order they run.
//
// The application is reached through httpCli, the database through db.
// The read checks run first, the checks that change data run after them:
// the owner's name and phone, the readings and the gauges.
// The checks that read the stored data run after the write checks they depend on.
func Steps(
	cfg config.Config,
	db clickhouse.Client,
	httpCli *http.Client,
) []step.Step {
	base := "http://localhost:" + cfg.AppPort()
	database := cfg.ClickHouseDatabase()
	return []step.Step{
		ping.NewStep(
			httpCli,
			base+"/ping",
		),
		schema.NewStep(
			db,
			database,
		),
		data.NewStep(
			db,
			database,
		),
		fixture.NewStep(
			db,
			database,
		),
		users.NewGet(httpCli, base),
		users.NewNotFound(httpCli, base),
		users.NewFio(db, database, httpCli, base),
		users.NewPhoneHistory(db, database, httpCli, base),
		users.NewPatch(httpCli, base),
		users.NewPatchBadRequest(httpCli, base),
		users.NewPatchNotFound(httpCli, base),
		realty.NewGet(httpCli, base),
		realty.NewNotFound(httpCli, base),
		readings.NewList(httpCli, base),
		readings.NewBadRequest(httpCli, base),
		readings.NewListNotFound(httpCli, base),
		readings.NewPost(httpCli, base),
		readings.NewPostBadRequest(httpCli, base),
		readings.NewPostNotFound(httpCli, base),
		readings.NewCorrection(httpCli, base),
		invoices.NewList(httpCli, base),
		invoices.NewPagination(httpCli, base),
		invoices.NewBadRequest(httpCli, base),
		invoices.NewNotFound(httpCli, base),
		tariffs.NewGet(httpCli, base),
		tariffs.NewBadRequest(httpCli, base),
		tariffs.NewInvalidDate(httpCli, base),
		tariffs.NewNotFound(httpCli, base),
		tariffs.NewBoundary(db, database, httpCli, base),
		gauges.NewPost(httpCli, base),
		gauges.NewPostBadRequest(httpCli, base),
		gauges.NewPostNotFound(httpCli, base),
		gauges.NewReplacement(db, database),
	}
}
