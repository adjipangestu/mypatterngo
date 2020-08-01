package controllers

import (
	"mypatterngo/app/helpers"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"status": "200",
		"messages": "hello bray",
	}
	helpers.Success(w, http.StatusOK, "Success", data)
	return
}
