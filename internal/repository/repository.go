package repository

import "assignment-5/internal/model"

type DatabaseRepo interface {
	GetAllCourses() ([]model.Course, error)
	GetCourseByID(id int) (model.Course, error)
	AddNewCourse(course model.Course) error
	DeleteCourse(id int) error
	UpdateCourse(course model.Course) error
}
