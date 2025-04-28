package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	DueDate     time.Time `json:"due_date"`
	IsCompleted bool      `json:"is_completed"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable()

	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/tasks", getTasks)
	app.Post("/tasks", createTask)
	app.Put("/tasks/:id/complete", completeTask)
	app.Put("/tasks/:id/uncomplete", uncompleteTask)
	app.Delete("/tasks/:id", deleteTask)

	log.Fatal(app.Listen(":3001"))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		due_date DATETIME,
		is_completed BOOLEAN DEFAULT 0
	);
	`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func getTasks(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, title, due_date, is_completed FROM tasks ORDER BY due_date ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var due string
		if err := rows.Scan(&t.ID, &t.Title, &due, &t.IsCompleted); err != nil {
			return err
		}
		t.DueDate, _ = time.Parse(time.RFC3339, due)
		tasks = append(tasks, t)
	}

	return c.JSON(tasks)
}

func createTask(c *fiber.Ctx) error {
	var t Task
	if err := c.BodyParser(&t); err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO tasks (title, due_date) VALUES (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(t.Title, t.DueDate.Format(time.RFC3339))
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func completeTask(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.Exec("UPDATE tasks SET is_completed = 1 WHERE id = ?", id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func uncompleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.Exec("UPDATE tasks SET is_completed = 0 WHERE id = ?", id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
