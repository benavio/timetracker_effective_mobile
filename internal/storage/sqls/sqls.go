package sqls

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type User struct {
	ID             int    `json:"id"`
	PassportNumber string `json:"passport_number"`
}

type Task struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Duration    int       `json:"duration"` // in minutes
}

func New(connStr string) (*Storage, error) {
	const opUser = "internal.storage.sqls.New.users"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s :%w", opUser, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		passport_number VARCHAR(11) NOT NULL UNIQUE
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s :%w", opUser, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s :%w", opUser, err)
	}

	const opTask = "internal.storage.sqls.New.task"

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		UserID INT NOT NULL REFERENCES users (id),
		Description TEXT NOT NULL,
		StartTime TIMESTAMP NOT NULL,
		EndTime TIMESTAMP NOT NULL
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s :%w", opTask, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s :%w", opTask, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddUser(PassportNumber string) (id int, err error) {
	op := "internal.storage.sqls.AddUser"

	stmt, err := s.db.Prepare("INSERT INTO users (passport_number) VALUES ($1) RETURNING id")
	if err != nil {
		return -1, fmt.Errorf("%s :%w", op, err)
	}

	var UserID int
	err = stmt.QueryRow(PassportNumber).Scan(&UserID)
	if err != nil {
		return -1, fmt.Errorf("%s :%w", op, err)
	}

	return UserID, nil
}

func (s *Storage) GetUsers() ([]User, error) {
	const op = "internal.storage.sqls.GetUsers"

	stmt, err := s.db.Prepare("SELECT * FROM users")
	if err != nil {
		return []User{}, fmt.Errorf("%s :%w", op, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return []User{}, fmt.Errorf("%s :%w", op, err)
	}

	users := []User{}

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.PassportNumber)
		if err != nil {
			return []User{}, fmt.Errorf("%s :%w", op, err)
		}
		users = append(users, user)
	}

	return users, nil
}

func UpdateUser() {

}

func (s *Storage) DeleteUser(passport_number string) error {
	op := "internal.storage.sqls.DeleteUser"

	stmt, err := s.db.Prepare("DELETE FROM users WHERE passport_number = $1")
	if err != nil {
		return fmt.Errorf("%s :%w", op, err)
	}

	result, err := stmt.Exec(passport_number)
	if err != nil {
		return fmt.Errorf("%s :%w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s :%w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s :%w", op, err)
	}

	return nil
}
