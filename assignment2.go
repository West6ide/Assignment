package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = 20030625
	dbname   = "first_db"
)

type Task struct {
	ID        int
	Name      string
	Completed bool
}

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%d dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database")

	var id int
	fmt.Scan(&id)

	var name string
	fmt.Scan(&name)

	err = createTask(db, id, name)
	if err != nil {
		panic(err)
	}
	fmt.Println("Задача успешно создана")

	//tasks, err := getAllTasks(db)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Список задач:")
	//for _, task := range tasks {
	//	fmt.Printf("ID: %d, Название: %s, Завершено: %t\n", task.ID, task.Name, task.Completed)
	//}

	//err = completeTask(db, 2) // Передайте ID нужной задачи
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Задача успешно отмечена как завершенная")

	//err = deleteTask(db, 987654) // Передайте ID нужной задачи
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Задача успешно удалена")

}

func createTask(db *sql.DB, id int, name string) error {
	_, err := db.Exec("INSERT INTO tasks (id, name) VALUES ($1, $2)", id, name)
	return err
}

func getAllTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, name, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func completeTask(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE tasks SET completed = true WHERE id = $1", id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func deleteTask(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
