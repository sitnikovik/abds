// Package lab1 builds the autograder checks of lab 1.
package lab1

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/clickhouse"
	"github.com/sitnikovik/abds/autograder/internal/config"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/data"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	gaugesPost "github.com/sitnikovik/abds/autograder/internal/step/gauges/post"
	gaugesPostBadRequest "github.com/sitnikovik/abds/autograder/internal/step/gauges/postbadrequest"
	gaugesPostNotFound "github.com/sitnikovik/abds/autograder/internal/step/gauges/postnotfound"
	gaugesReplacement "github.com/sitnikovik/abds/autograder/internal/step/gauges/replacement"
	invoicesBadRequest "github.com/sitnikovik/abds/autograder/internal/step/invoices/badrequest"
	invoicesList "github.com/sitnikovik/abds/autograder/internal/step/invoices/list"
	invoicesNotFound "github.com/sitnikovik/abds/autograder/internal/step/invoices/notfound"
	invoicesPagination "github.com/sitnikovik/abds/autograder/internal/step/invoices/pagination"
	"github.com/sitnikovik/abds/autograder/internal/step/ping"
	readingsBadRequest "github.com/sitnikovik/abds/autograder/internal/step/readings/badrequest"
	readingsCorrection "github.com/sitnikovik/abds/autograder/internal/step/readings/correction"
	readingsList "github.com/sitnikovik/abds/autograder/internal/step/readings/list"
	readingsListNotFound "github.com/sitnikovik/abds/autograder/internal/step/readings/listnotfound"
	readingsPost "github.com/sitnikovik/abds/autograder/internal/step/readings/post"
	readingsPostBadRequest "github.com/sitnikovik/abds/autograder/internal/step/readings/postbadrequest"
	readingsPostNotFound "github.com/sitnikovik/abds/autograder/internal/step/readings/postnotfound"
	realtyGet "github.com/sitnikovik/abds/autograder/internal/step/realty/get"
	realtyNotFound "github.com/sitnikovik/abds/autograder/internal/step/realty/notfound"
	"github.com/sitnikovik/abds/autograder/internal/step/schema"
	tariffsBadRequest "github.com/sitnikovik/abds/autograder/internal/step/tariffs/badrequest"
	tariffsBoundary "github.com/sitnikovik/abds/autograder/internal/step/tariffs/boundary"
	tariffsGet "github.com/sitnikovik/abds/autograder/internal/step/tariffs/get"
	tariffsInvalidDate "github.com/sitnikovik/abds/autograder/internal/step/tariffs/invaliddate"
	tariffsNotFound "github.com/sitnikovik/abds/autograder/internal/step/tariffs/notfound"
	usersFio "github.com/sitnikovik/abds/autograder/internal/step/users/fio"
	usersGet "github.com/sitnikovik/abds/autograder/internal/step/users/get"
	usersNotFound "github.com/sitnikovik/abds/autograder/internal/step/users/notfound"
	usersPatch "github.com/sitnikovik/abds/autograder/internal/step/users/patch"
	usersPatchBadRequest "github.com/sitnikovik/abds/autograder/internal/step/users/patchbadrequest"
	usersPatchNotFound "github.com/sitnikovik/abds/autograder/internal/step/users/patchnotfound"
	usersPhoneHistory "github.com/sitnikovik/abds/autograder/internal/step/users/phonehistory"
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
		usersGet.NewStep(httpCli, base),
		usersNotFound.NewStep(httpCli, base),
		usersFio.NewStep(db, database, httpCli, base),
		usersPhoneHistory.NewStep(db, database, httpCli, base),
		usersPatch.NewStep(httpCli, base),
		usersPatchBadRequest.NewStep(httpCli, base),
		usersPatchNotFound.NewStep(httpCli, base),
		realtyGet.NewStep(httpCli, base),
		realtyNotFound.NewStep(httpCli, base),
		readingsList.NewStep(httpCli, base),
		readingsBadRequest.NewStep(httpCli, base),
		readingsListNotFound.NewStep(httpCli, base),
		readingsPost.NewStep(httpCli, base),
		readingsPostBadRequest.NewStep(httpCli, base),
		readingsPostNotFound.NewStep(httpCli, base),
		readingsCorrection.NewStep(httpCli, base),
		invoicesList.NewStep(httpCli, base),
		invoicesPagination.NewStep(httpCli, base),
		invoicesBadRequest.NewStep(httpCli, base),
		invoicesNotFound.NewStep(httpCli, base),
		tariffsGet.NewStep(httpCli, base),
		tariffsBadRequest.NewStep(httpCli, base),
		tariffsInvalidDate.NewStep(httpCli, base),
		tariffsNotFound.NewStep(httpCli, base),
		tariffsBoundary.NewStep(db, database, httpCli, base),
		gaugesPost.NewStep(httpCli, base),
		gaugesPostBadRequest.NewStep(httpCli, base),
		gaugesPostNotFound.NewStep(httpCli, base),
		gaugesReplacement.NewStep(db, database),
	}
}
