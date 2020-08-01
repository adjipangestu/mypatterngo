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
	//PublicRouter.HandleFunc("/login", controllers.Login).Methods("POST")
	//PublicRouter.HandleFunc("/forgot-password", controllers.ForgotPassword).Methods("POST")
	//PublicRouter.HandleFunc("/change-password/{token}", controllers.ChangePassword).Methods("PATCH")

	// High Admin Routes
	//ProtectedRouterHighAdminRouter.HandleFunc("/roles", controllers.GetAllRoles).Methods("GET")
	//ProtectedRouterHighAdminRouter.HandleFunc("/roles", controllers.CreateRole).Methods("POST")
	//ProtectedRouterHighAdminRouter.HandleFunc("/roles/{id}", controllers.GetRole).Methods("GET")
	//ProtectedRouterHighAdminRouter.HandleFunc("/roles/{id}", controllers.UpdateRole).Methods("PATCH")
	//ProtectedRouterHighAdminRouter.HandleFunc("/roles/{id}", controllers.DeleteRole).Methods("DELETE")
	//ProtectedRouterHighAdminRouter.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")

	// Protected Routes
	//ProtectedRouter.HandleFunc("/users/me", controllers.GetAuthenticatedUser).Methods("GET")
	//ProtectedRouter.HandleFunc("/users/me/upload-image", controllers.UploadUserImage).Methods("PATCH")
	//ProtectedRouter.HandleFunc("/users/me/delete-image", controllers.DeleteImage).Methods("DELETE")
}
