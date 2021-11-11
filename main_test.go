package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("root", "Kamath@123", "golang_restAPI")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM tasks")
	a.DB.Exec("ALTER TABLE tasks AUTO_INCREMENT = 1")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS tasks
(
	id INT NOT NULL AUTO_INCREMENT,
	title VARCHAR(50) NOT NULL,
	description VARCHAR(50) NOT NULL,
	assignedTo VARCHAR(50) NOT NULL,
	PRIMARY KEY(id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentTask(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/task/45", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Task not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Task not found'. Got '%s'", m["error"])
	}
}

func TestCreateTask(t *testing.T) {
	clearTable()
	payload := []byte(`{"Title":"test Title","Description":"test Description", "AssignedTo": "test AssignedTo"}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["Title"] != "test Title" {
		t.Errorf("Expected task title to be 'test Title'. Got '%v'", m["Title"])
	}
	if m["Description"] != "test Description" {
		t.Errorf("Expected task Description to be 'test Description'. Got '%v'", m["Description"])
	}
	if m["AssignedTo"] != "test AssignedTo" {
		t.Errorf("Expected task assignedTo age to be 'test AssignedTo'. Got '%v'", m["AssignedTo"])
	}
	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["ID"] != 1.0 {
		t.Errorf("Expected task ID to be '1'. Got '%v'", m["ID"])
	}
}

func TestGetTask(t *testing.T) {
	clearTable()
	addTasks(1)
	req, _ := http.NewRequest("GET", "/task/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func addTasks(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO tasks(title, description, assignedTo) VALUES('%s', '%s', '%s')", ("Task " + strconv.Itoa(i+1)), ("Task " + strconv.Itoa(i+1)), ("Task " + strconv.Itoa(i+1)))
		a.DB.Exec(statement)
	}
}

func TestUpdateTask(t *testing.T) {
	clearTable()
	addTasks(1)
	req, _ := http.NewRequest("GET", "/task/1", nil)
	response := executeRequest(req)
	var originalTask map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalTask)
	payload := []byte(`{"title":"test title - updated title","description":"loremmm", "assignedTo": "pavan"}`)
	req, _ = http.NewRequest("PUT", "/task/1", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["ID"] != originalTask["ID"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalTask["ID"], m["ID"])
	}
	if m["Title"] == originalTask["Title"] {
		t.Errorf("Expected the title to change from '%v' to '%v'. Got '%v'", originalTask["Title"], m["Title"], m["Title"])
	}
	if m["Description"] == originalTask["Description"] {
		t.Errorf("Expected the description to change from '%v' to '%v'. Got '%v'", originalTask["Description"], m["Description"], m["Description"])
	}
	if m["AssignedTo"] == originalTask["AssignedTo"] {
		t.Errorf("Expected the assignedTo to change from '%v' to '%v'. Got '%v'", originalTask["AssignedTo"], m["AssignedTo"], m["AssignedTo"])
	}
}

func TestDeleteTask(t *testing.T) {
	clearTable()
	addTasks(1)
	req, _ := http.NewRequest("GET", "/task/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("DELETE", "/task/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("GET", "/task/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
