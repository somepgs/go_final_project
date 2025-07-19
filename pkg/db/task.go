package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask inserts a new task into the database and returns the ID of the newly created task.
func AddTask(task *Task) (int64, error) {
	var id int64
	// Prepare the SQL statement to insert a new task
	stmt := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(stmt, task.Date, task.Title, task.Comment, task.Repeat)
	// Check for errors during the execution of the query
	if err == nil {
		id, err = result.LastInsertId()
	}
	return id, err
}

// Tasks retrieves a limited number of tasks from the database, ordered by date.
func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getTasks(rows)
}

// SearchTasks searches for tasks by title, comment, or date.
func SearchTasks(search string, limit int) ([]*Task, error) {
	if date, err := time.Parse("02.01.2006", search); err == nil {
		rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit",
			sql.Named("date", date.Format("20060102")), sql.Named("limit", limit))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return getTasks(rows)
	}

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE LOWER(title) LIKE LOWER(:search) "+
		"OR LOWER(comment) LIKE LOWER(:search) ORDER BY date LIMIT :limit",
		sql.Named("search", "%"+search+"%"), sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getTasks(rows)
}

// GetTask retrieves a task by its ID from the database.
func GetTask(id string) (*Task, error) {
	var task Task
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No task found with the given ID
		}
		return nil, err // Return any other error
	}
	return &task, nil
}

// UpdateTask updates an existing task in the database.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// DeleteTask removes a task from the database by its ID.
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id`
	res, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}

// UpdateDate updates the date of a task in the database by its ID.
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := db.Exec(query, sql.Named("date", next), sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task date`)
	}
	return nil
}

// getTasks scans the rows returned by a query and returns a slice of Task pointers.
func getTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
