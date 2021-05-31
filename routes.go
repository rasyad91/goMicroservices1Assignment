package main

import (
	"assignment-5/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func routes() http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", handler.Repo.Home).Methods("GET")
	router.HandleFunc("/api/v1/courses", handler.Repo.AllCourses).Methods("GET")

	sub := router.NewRoute().Subrouter()
	sub.Use(handler.ValidationAPIMiddleware)

	sub.HandleFunc("/api/v1/courses/{courseId}", handler.Repo.GetCourse).Methods("GET")
	sub.HandleFunc("/api/v1/courses/{courseId}", handler.Repo.PostCourse).Methods("POST")
	sub.HandleFunc("/api/v1/courses/{courseId}", handler.Repo.PutCourse).Methods("PUT")
	sub.HandleFunc("/api/v1/courses/{courseId}", handler.Repo.DeleteCourse).Methods("DELETE")
	return router
}
