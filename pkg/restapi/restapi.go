package restapi

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pantheon-systems/pod-heartbeat/pkg/heartbeat"
)

// API is the state management for webapi
type API struct {
	check  *heartbeat.Check
	Router *httprouter.Router
}

// NewAPI is the constructor for the API object
func NewAPI(c *heartbeat.Check) *API {
	router := httprouter.New()
	a := API{
		check:  c,
		Router: router,
	}
	router.GET("/", a.Status)
	return &a
}

// Status reports the status of the checker. 200 if ok 503 if not
func (a *API) Status(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if a.check.OK != true {
		http.Error(w, "Status check failed", http.StatusInternalServerError)
	} else {
		fmt.Fprint(w, "OK")
	}
}
