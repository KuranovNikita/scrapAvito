package handler_readiness

import (
	"net/http"
	"scrapAvito/json_app"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	json_app.RespondWithJSON(w, 200, struct{}{})
}
