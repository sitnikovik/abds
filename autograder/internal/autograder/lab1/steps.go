// Package lab1 builds the autograder checks of lab 1.
package lab1

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/clickhouse"
	"github.com/sitnikovik/abds/autograder/internal/config"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/data"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	gauges_post "github.com/sitnikovik/abds/autograder/internal/step/gauges/post"
	gauges_post_bad_request "github.com/sitnikovik/abds/autograder/internal/step/gauges/postbadrequest"
	gauges_post_not_found "github.com/sitnikovik/abds/autograder/internal/step/gauges/postnotfound"
	gauges_replacement "github.com/sitnikovik/abds/autograder/internal/step/gauges/replacement"
	invoices_bad_request "github.com/sitnikovik/abds/autograder/internal/step/invoices/badrequest"
	invoices_list "github.com/sitnikovik/abds/autograder/internal/step/invoices/list"
	invoices_not_found "github.com/sitnikovik/abds/autograder/internal/step/invoices/notfound"
	invoices_pagination "github.com/sitnikovik/abds/autograder/internal/step/invoices/pagination"
	"github.com/sitnikovik/abds/autograder/internal/step/ping"
	readings_bad_request "github.com/sitnikovik/abds/autograder/internal/step/readings/badrequest"
	readings_correction "github.com/sitnikovik/abds/autograder/internal/step/readings/correction"
	readings_list "github.com/sitnikovik/abds/autograder/internal/step/readings/list"
	readings_list_not_found "github.com/sitnikovik/abds/autograder/internal/step/readings/listnotfound"
	readings_post "github.com/sitnikovik/abds/autograder/internal/step/readings/post"
	readings_post_bad_request "github.com/sitnikovik/abds/autograder/internal/step/readings/postbadrequest"
	readings_post_not_found "github.com/sitnikovik/abds/autograder/internal/step/readings/postnotfound"
	realty_get "github.com/sitnikovik/abds/autograder/internal/step/realty/get"
	realty_not_found "github.com/sitnikovik/abds/autograder/internal/step/realty/notfound"
	"github.com/sitnikovik/abds/autograder/internal/step/schema"
	tariffs_bad_request "github.com/sitnikovik/abds/autograder/internal/step/tariffs/badrequest"
	tariffs_boundary "github.com/sitnikovik/abds/autograder/internal/step/tariffs/boundary"
	tariffs_get "github.com/sitnikovik/abds/autograder/internal/step/tariffs/get"
	tariffs_invalid_date "github.com/sitnikovik/abds/autograder/internal/step/tariffs/invaliddate"
	tariffs_not_found "github.com/sitnikovik/abds/autograder/internal/step/tariffs/notfound"
	users_fio "github.com/sitnikovik/abds/autograder/internal/step/users/fio"
	users_get "github.com/sitnikovik/abds/autograder/internal/step/users/get"
	users_not_found "github.com/sitnikovik/abds/autograder/internal/step/users/notfound"
	users_patch "github.com/sitnikovik/abds/autograder/internal/step/users/patch"
	users_patch_bad_request "github.com/sitnikovik/abds/autograder/internal/step/users/patchbadrequest"
	users_patch_not_found "github.com/sitnikovik/abds/autograder/internal/step/users/patchnotfound"
	users_phone_history "github.com/sitnikovik/abds/autograder/internal/step/users/phonehistory"
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
		users_get.NewStep(httpCli, base),
		users_not_found.NewStep(httpCli, base),
		users_fio.NewStep(db, database, httpCli, base),
		users_phone_history.NewStep(db, database, httpCli, base),
		users_patch.NewStep(httpCli, base),
		users_patch_bad_request.NewStep(httpCli, base),
		users_patch_not_found.NewStep(httpCli, base),
		realty_get.NewStep(httpCli, base),
		realty_not_found.NewStep(httpCli, base),
		readings_list.NewStep(httpCli, base),
		readings_bad_request.NewStep(httpCli, base),
		readings_list_not_found.NewStep(httpCli, base),
		readings_post.NewStep(httpCli, base),
		readings_post_bad_request.NewStep(httpCli, base),
		readings_post_not_found.NewStep(httpCli, base),
		readings_correction.NewStep(httpCli, base),
		invoices_list.NewStep(httpCli, base),
		invoices_pagination.NewStep(httpCli, base),
		invoices_bad_request.NewStep(httpCli, base),
		invoices_not_found.NewStep(httpCli, base),
		tariffs_get.NewStep(httpCli, base),
		tariffs_bad_request.NewStep(httpCli, base),
		tariffs_invalid_date.NewStep(httpCli, base),
		tariffs_not_found.NewStep(httpCli, base),
		tariffs_boundary.NewStep(db, database, httpCli, base),
		gauges_post.NewStep(httpCli, base),
		gauges_post_bad_request.NewStep(httpCli, base),
		gauges_post_not_found.NewStep(httpCli, base),
		gauges_replacement.NewStep(db, database),
	}
}
