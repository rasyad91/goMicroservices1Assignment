package mysql

import (
	"assignment-5/internal/model"
	"assignment-5/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBrepo struct {
	*sql.DB
}

// NewRepo creates the repository
func NewRepo(Conn *sql.DB) repository.DatabaseRepo {
	return &DBrepo{
		DB: Conn,
	}
}

func (m *DBrepo) GetAllCourses() ([]model.Course, error) {
	courses := make([]model.Course, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, created_at, updated_at from courses`

	rows, err := m.QueryContext(ctx, query)
	if err != nil {
		return courses, err
	}
	defer rows.Close()

	for rows.Next() {
		course := model.Course{}

		err := rows.Scan(&course.ID, &course.Title, &course.CreatedAt, &course.UpdatedAt)
		if err != nil {
			return courses, err
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return courses, err
	}

	return courses, nil
}

func (m *DBrepo) GetCourseByID(id int) (model.Course, error) {

	course := model.Course{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title from courses where id = ?`

	row := m.QueryRowContext(ctx, query, id)
	err := row.Scan(&course.ID, &course.Title)
	if err != nil {
		return course, err
	}

	return course, nil
}

func (r *DBrepo) AddNewCourse(course model.Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO 
				courses(id, title, created_at, updated_at)
    		VALUES
				(?, ?, ? ,?);
			`

	if _, err := r.ExecContext(ctx, stmt, course.ID, course.Title, time.Now(), time.Now()); err != nil {
		return err
	}

	return nil
}

func (r *DBrepo) DeleteCourse(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM
				courses
			WHERE id = ?;
			`
	result, err := r.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	if x, _ := result.RowsAffected(); x == 0 {
		return fmt.Errorf("id not found")
	}
	return nil

}

func (r *DBrepo) UpdateCourse(course model.Course) error {
	fmt.Println(course)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE courses
			SET
				title = ?, updated_at = ?
    		WHERE
				id = ?
			`

	result, err := r.ExecContext(ctx, stmt, course.Title, time.Now(), course.ID)
	if err != nil {
		return err
	}
	if x, _ := result.RowsAffected(); x == 0 {
		return fmt.Errorf("id not found")
	}

	return nil
}
