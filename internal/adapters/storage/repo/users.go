package repo

import (
	"database/sql"
	"log"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/storage"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db}
}

func (ur *UsersRepo) Create(login string, password string) (int, error) {
	id := 0
	err := ur.db.QueryRow(storage.UsersQueryCreate, login, password).Scan(&id)
	if err != nil {
		log.Println("Create. QueryRow. err: " + err.Error())
		return 0, err
	}

	log.Println("Create. err: nil")
	return id, nil
}

func (ur *UsersRepo) Delete(login string) error {
	_, err := ur.db.Exec(storage.UsersQueryDelete, login)
	if err != nil {
		log.Println("Delete. Exec. err: " + err.Error())
		return err
	}

	log.Println("Delete. err: nil")
	return nil
}

func (ur *UsersRepo) Check(login string) (bool, error) {
	var exists bool
	err := ur.db.QueryRow(storage.UsersQueryCheck, login).Scan(&exists)
	if err != nil {
		log.Println("Check. QueryRow. err: " + err.Error())
		return false, err
	}

	log.Println("Check. err: nil")
	return exists, nil
}

func (ur *UsersRepo) Get(login string) (string, error) {
	var password string
	err := ur.db.QueryRow(storage.UsersQueryGet, login).Scan(&password)
	if err != nil {
		log.Println("Get. QueryRow. err: " + err.Error())
		return "", err
	}

	log.Println("Get. err: nil")
	return password, nil
}
