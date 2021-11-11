package main

import (
	"database/sql"
	"fmt"
)

type task struct {
	ID          int    `json: "id"`
	Title       string `json: "title"`
	Description string `json: 'description"`
	AssignedTo  string `json: "assignedTo"`
}

func (t *task) getTask(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT title, description, assignedTo FROM tasks WHERE id=%d", t.ID)
	return db.QueryRow(statement).Scan(&t.Title, &t.Description, &t.AssignedTo)
}

func (t *task) updateTask(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE tasks SET title='%s', description='%s', assignedTo='%s' WHERE id=%d", t.Title, t.Description, t.AssignedTo, t.ID)
	_, err := db.Exec(statement)
	return err
}

func (t *task) deleteTask(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM tasks WHERE id=%d", t.ID)
	_, err := db.Exec(statement)
	return err
}

func (t *task) createTask(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO tasks(title, description, assignedTo) VALUES('%s', '%s','%s')", t.Title, t.Description, t.AssignedTo)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}

func getTasks(db *sql.DB, start, count int) ([]task, error) {
	statement := fmt.Sprintf("SELECT id, title, description, assignedTo FROM tasks LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []task{}
	for rows.Next() {
		var t task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.AssignedTo); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
