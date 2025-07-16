package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) JWTMiddlewareWithRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return app.JWTMiddleware(next, requiredRole)
	}
}

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, makeResponseJSON)
	authMiddleware := standardMiddleware.Append(app.JWTMiddlewareWithRole(""))

	mux := pat.New()

	// Users
	mux.Post("/user/sign_up", standardMiddleware.ThenFunc(app.userHandler.SignUp))
	mux.Post("/user/sign_in", standardMiddleware.ThenFunc(app.userHandler.SignIn))
	mux.Post("/user/logout", authMiddleware.ThenFunc(app.userHandler.Logout))

	// Reviews
	mux.Post("/reviews", authMiddleware.ThenFunc(app.reviewHandler.Create))
	mux.Get("/reviews", standardMiddleware.ThenFunc(app.reviewHandler.GetAll))
	mux.Get("/reviews/:id", standardMiddleware.ThenFunc(app.reviewHandler.GetByID))
	mux.Put("/reviews/:id", authMiddleware.ThenFunc(app.reviewHandler.Update))
	mux.Del("/reviews/:id", authMiddleware.ThenFunc(app.reviewHandler.Delete))

	// Review images
	mux.Get("/images/reviews/:filename", http.HandlerFunc(app.reviewHandler.ServeReviewImage))

	// Review pdfs
	mux.Get("/pdfs/reviews/:filename", http.HandlerFunc(app.reviewHandler.ServeReviewPDF))

	// mux.Get("/swagger/", httpSwagger.WrapHandler)

	return standardMiddleware.Then(mux)
}
