package handlers

import (
	"log"
	"net/http"

	"github.com/elitenomad/garagesale/internal/platform/database"
	"github.com/elitenomad/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

type Check struct {
	Db  *sqlx.DB
	Log *log.Logger
}

func (c *Check) Health(w http.ResponseWriter, r *http.Request) error {
	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(r.Context(), c.Db); err != nil {
		health.Status = "Not connected..."
		return web.Respond(w, health, http.StatusInternalServerError)
	}

	health.Status = "Successful"
	return web.Respond(w, health, http.StatusOK)
}
