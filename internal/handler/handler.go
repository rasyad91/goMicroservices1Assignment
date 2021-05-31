package handler

import (
	"assignment-5/internal/driver/mysqlDriver"
	"assignment-5/internal/model"
	"assignment-5/internal/repository"
	"assignment-5/internal/repository/mysql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Repository struct {
	DB    repository.DatabaseRepo
	Error *log.Logger
	Info  *log.Logger
}

var Repo *Repository

// NewMySQLHandlers creates db repo for postgres
func NewMySQLHandlers(db *mysqlDriver.DB, errorLog, infoLog *log.Logger) *Repository {
	return &Repository{
		DB:    mysql.NewRepo(db.SQL),
		Error: errorLog,
		Info:  infoLog,
	}
}

// NewHandlers creates the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func ValidationAPIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		v := r.URL.Query()
		key, ok := v["key"]
		if !ok || key[0] != "2c78afaf-97da-4816-bbee-9ad239abb296" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("401 - Invalid key"))
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to REST API\n")
}

func (m *Repository) AllCourses(w http.ResponseWriter, r *http.Request) {

	courses, err := m.DB.GetAllCourses()
	if err != nil {
		m.Error.Println("GetAllCourses: ", err)
	}

	for _, course := range courses {
		if err := json.NewEncoder(w).Encode(course); err != nil {
			m.Error.Println(err)
		}
	}
}

func (m *Repository) PostCourse(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-type") == "application/json" {
		// read the string sent to the service
		var newCourse model.Course
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply course information in JSON format\n"))
			return
		}
		// convert JSON to object
		json.Unmarshal(reqBody, &newCourse)
		if newCourse.Title == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply course information in JSON format\n"))
			return
		}
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["courseId"])
		if err != nil {
			m.Error.Println(err)
			return
		}
		newCourse.ID = id
		if err := m.DB.AddNewCourse(newCourse); err != nil {
			if strings.Contains(err.Error(), "Duplicate") {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - Existing ID\n"))
			} else {
				m.Error.Println(err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("201 - Course added: %d\n", id)))
	}
}

func (m *Repository) PutCourse(w http.ResponseWriter, r *http.Request) {

	// read the string sent to the service
	var newCourse model.Course
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply course information in JSON format\n"))
		return
	}
	// convert JSON to object
	json.Unmarshal(reqBody, &newCourse)
	if newCourse.Title == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply course information in JSON format\n"))
		return
	}
	fmt.Println(newCourse)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["courseId"])
	if err != nil {
		m.Error.Println(err)
		return
	}
	newCourse.ID = id
	err = m.DB.UpdateCourse(newCourse)
	if err != nil && err.Error() == "id not found" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No course found\n"))
		return
	}
	if err != nil {
		m.Error.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("200 - Course updated: %d\n", id)))
}

func (m *Repository) GetCourse(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["courseId"])
	if err != nil {
		m.Error.Println(err)
		return
	}

	course, err := m.DB.GetCourseByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No course found\n"))
		m.Error.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(course)
}

func (m *Repository) DeleteCourse(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["courseId"])
	if err != nil {
		m.Error.Println(err)
		return
	}
	if err := m.DB.DeleteCourse(id); err != nil && err.Error() == "id not found" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No course found\n"))

	} else if err != nil {
		m.Error.Println(err)

	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("200 - Course id: %d deleted\n", id)))
		m.Info.Println("Successfully deleted course with id: ", id)
	}
}
