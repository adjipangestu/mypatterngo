package routes

import (
	"mypatterngo/app/controllers"
	"mypatterngo/app/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	Router *mux.Router
}

// ServeRoutes handle the public routes
func (api *Api) InitializeRoutes() {
	api.Router = mux.NewRouter()

	// Server static file
	var imgServer = http.FileServer(http.Dir("./static/"))
	api.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", imgServer))

	// Route List
	GuestRouter := api.Router.PathPrefix("/api").Subrouter()
	AuthRouter := api.Router.PathPrefix("/api/v1").Subrouter()
	AuthAdminRouter := api.Router.PathPrefix("/api/v1").Subrouter()

	//Middleware
	GuestRouter.Use(middlewares.GuestMiddleware)
	AuthAdminRouter.Use(middlewares.GuestMiddleware)
	AuthRouter.Use(middlewares.GuestMiddleware)
	AuthRouter.Use(middlewares.AuthenticationMiddleware)
	AuthAdminRouter.Use(middlewares.AuthenticationMiddleware)
	AuthAdminRouter.Use(middlewares.AdminMiddleware)

	// Open Routes
	GuestRouter.HandleFunc("/welcome", controllers.Welcome).Methods("GET")
}
